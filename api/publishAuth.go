package api

import (
	"github.com/dgrijalva/jwt-go"
	"livestreamall/config"
	"livestreamall/dao"
	"livestreamall/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublishAuth(c *gin.Context) {
	// 从请求参数中获取 token 和推流信息
	streamName := c.PostForm("name") // 推流名称
	app := c.PostForm("app")         // RTMP 应用名称
	claim := c.PostForm("claim")     // 此处就是用户token
	ip := c.ClientIP()               // 推流客户端 IP

	// 验证必要参数
	if app == "" || streamName == "" || claim == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		c.Abort()
		return
	}

	// 验证 token 是否有效
	token, err := jwt.ParseWithClaims(claim, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return
	}

	// 从 token 中提取用户 ID
	var userID string
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		userID = claims.UID
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	// 检查直播间是否已存在
	var existingRoom model.LiveRoom
	err = dao.DB.Where("stream_name = ?", streamName).First(&existingRoom).Error
	if err == nil {
		// 如果直播间已存在且正在直播，拒绝推流
		if existingRoom.IsLive {
			c.JSON(http.StatusConflict, gin.H{"error": "stream is already live"})
			return
		} else if existingRoom.UserID != userID {
			// 不是这个直播间的主播
			c.JSON(http.StatusForbidden, gin.H{"error": "user not allowed"})
			return
		}
		// 更新直播间状态为直播中
		existingRoom.IsLive = true
		dao.DB.Save(&existingRoom)
	} else {
		// 如果直播间不存在,拒绝
		c.JSON(http.StatusForbidden, gin.H{"error": "live room not exist"})
		return
	}

	// 返回成功响应，允许 Nginx 推流
	c.JSON(http.StatusOK, gin.H{
		"message":     "publish authorized",
		"stream_name": streamName,
		"ip":          ip,
		"user_id":     userID,
	})
}

func StopPublish(c *gin.Context) {
	// 从请求中获取流名称
	streamName := c.PostForm("name")
	app := c.PostForm("app") // RTMP 应用名称
	// ip := c.ClientIP()       // 推流客户端 IP

	if app == "" || streamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		return
	}

	// 查找直播间
	var existingRoom model.LiveRoom
	err := dao.DB.Where("stream_name = ?", streamName).First(&existingRoom).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "live room not found"})
		return
	}

	// 更新直播间状态为未直播
	existingRoom.IsLive = false
	dao.DB.Save(&existingRoom)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{

		"message":     "stream stopped, status updated",
		"stream_name": streamName,
		"push_url":    config.Stream.PushBaseURL,
	})
}
