package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"livestreamall/config"
	"livestreamall/dao"
	"livestreamall/model"
	"net/http"
	"time"
)

// GetPullStreamURL 获取拉流相关地址
func GetPullStreamURL(c *gin.Context) {
	// 从请求中获取参数
	streamName := c.Param("stream_name")

	// 验证参数
	if streamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing stream_name"})
		return
	}

	// 查询直播间是否存在
	var liveRoom model.LiveRoom
	err := dao.DB.Where("stream_name = ?", streamName).First(&liveRoom).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "live room not found or not live"})
		return
	}

	// 动态生成拉流地址
	// 动态生成拉流和推流地址
	hlsURL := fmt.Sprintf("%s/hls/%s.m3u8", config.Stream.PullBaseURL, streamName)
	dashURL := fmt.Sprintf("%s/dash/%s.mpd", config.Stream.PullBaseURL, streamName)
	playURL := fmt.Sprintf("http://%s:%d/live/play?stream_id=%s", config.App.Host, config.App.Port, streamName)

	// 返回拉流地址
	c.JSON(http.StatusOK, gin.H{
		"hls_url":  hlsURL,
		"dash_url": dashURL,
		"play_url": playURL,
	})
}

// GenerateSecureURL 生成带签名的拉流地址
func GenerateSecureURL(baseURL, streamName, secret string) string {
	// 设置过期时间
	expiration := time.Now().Add(1 * time.Hour).Unix() // 1小时有效期

	// 生成签名字符串
	data := fmt.Sprintf("%s|%d", streamName, expiration)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	// 拼接带签名的地址
	secureURL := fmt.Sprintf("%s/%s.m3u8?expires=%d&signature=%s", baseURL, streamName, expiration, signature)
	return secureURL
}

func GetLiveRooms(c *gin.Context) {

	// 查询数据库中的直播间信息
	var liveRooms []model.LiveRoom
	if err := dao.DB.Where("is_live = ?", true).Find(&liveRooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch live rooms"})
		return
	}

	// 返回直播间列表
	c.JSON(http.StatusOK, gin.H{
		"live_rooms": liveRooms,
	})
}

func GetLiveRoomByStreamName(c *gin.Context) {
	// 获取 URL 参数中的 StreamName
	streamName := c.Param("stream_name")

	// 检查参数是否存在
	if streamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stream_name parameter is required"})
		return
	}

	// 查询数据库，查找对应的直播间信息
	var liveRoom model.LiveRoom
	if err := dao.DB.Where("stream_name = ?", streamName).First(&liveRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "live room not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch live room"})
		}
		return
	}

	// 返回直播间数据
	c.JSON(http.StatusOK, gin.H{
		"id":          liveRoom.ID,
		"title":       liveRoom.Title,
		"description": liveRoom.Description,
		"is_live":     liveRoom.IsLive,
	})
}

// StartLive 开始直播，生成推流密钥和流名称
func StartLive(c *gin.Context) {
	var request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 生成 JWT 推流密钥
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("uID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 检查直播间是否已存在
	var existingRoom model.LiveRoom
	err := dao.DB.Where("user_id = ?", userID).First(&existingRoom).Error

	if err == nil {
		// 如果直播间已存在且为当前用户的直播间，更新直播间信息
		existingRoom.Title = request.Title
		existingRoom.StreamName = token.(string)
		existingRoom.Description = request.Description
		dao.DB.Save(&existingRoom) // 更新数据库中的记录

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"stream_name": userID,
			"push_url":    config.Stream.PushBaseURL,
			"token":       token,
			"user_id":     userID,
		})
	} else {
		// 创建新的直播间记录
		newRoom := model.LiveRoom{
			Title:       request.Title,
			StreamName:  token.(string),
			Description: request.Description,
			IsLive:      false,
			UserID:      userID.(string), // 关联用户ID
		}
		dao.DB.Create(&newRoom) // 创建新记录

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"stream_name": userID,
			"push_url":    config.Stream.PushBaseURL,
			"token":       token,
			"user_id":     userID,
		})
	}
}

//// StreamForwarding 处理客户端拉取 HLS 流
//func StreamForwarding(c *gin.Context) {
//	// 获取客户端请求中的 streamID
//	streamID := c.Param("stream_id")
//
//	// 查询直播间信息，通过 streamID 获取对应的 liveRoom
//	var liveRoom model.LiveRoom
//	err := dao.DB.Where("user_id = ?", streamID).First(&liveRoom).Error
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
//		return
//	}
//
//	// 构建实际的 HLS 流地址
//	realStreamName := liveRoom.StreamName // 或根据需要加上 liveRoom.UserID 来进行流名称的映射
//	m3u8URL := fmt.Sprintf("http://127.0.0.1:8080/hls/%s.m3u8", realStreamName)
//
//	// 将 HLS 流的 m3u8 文件返回给客户端
//	resp, err := http.Get(m3u8URL)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get m3u8 stream"})
//		return
//	}
//	defer resp.Body.Close()
//
//	// 设置 HTTP 响应头为 m3u8 文件类型
//	c.Header("Content-Type", "application/vnd.apple.mpegurl")
//
//	// 将 m3u8 文件内容传递给客户端
//	_, err = io.Copy(c.Writer, resp.Body)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to forward m3u8 stream"})
//		return
//	}
//}
