package api

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"livestreamall/config"
	"livestreamall/model"
	"livestreamall/service"
	"livestreamall/util"
	"net/http"
	"time"
)

var jwtKey = []byte("12190711")

func SetToken(mail string, password string, c *gin.Context) {

	// 获取密码
	uID, nickname, userPassword, err := service.SearchUserPassword("mail", mail)

	if err != nil {
		if err == sql.ErrNoRows {
			//处理该用户不存在(回复邮箱或密码错误)
			util.RespIncorrectPassword(c)
		} else if err == util.FieldsError {
			//处理数据库查询字段不符的错误
			util.RespFieldsMatchError(c)
		} else if err != sql.ErrNoRows {
			//处理意料之外的错误
			util.RespUnexceptedError(c)
		}
		return
	}

	var user struct {
		UID      string `json:"uID"`
		Password string `json:"password"`
		Username string `json:"Username"`
	}

	user.UID = uID
	user.Password = userPassword
	user.Username = nickname

	/*
		if err := c.ShouldBindJSON(&user); err != nil {
			util.RespUnexceptedError(c)
			fmt.Println("执行到3了")
			fmt.Printf("err: %v\n", err)
			return
		}
	*/

	// 验证密码
	if userPassword != password {
		util.RespIncorrectPassword(c)
		return
	}
	// 创建 JWT
	expireTime := time.Now().Add(time.Hour * 24).Unix()
	claims := &model.Claims{
		UID:      user.UID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't generate JWT"})
		util.RespUnexceptedError(c)
		return
	}
	// 返回 JWT
	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
	c.SetCookie("token", tokenString, 604800, "/", "", false, false)
	util.RespSetTokenSuccess(c, tokenString)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			//util.RespDidNotLogin(c)
			c.Redirect(http.StatusFound, fmt.Sprintf("http://%s:%d/user/login", config.App.Host, config.App.Port))
			c.Abort()
			return
		}
		// 验证 token
		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			//util.RespInvalidToken(c)
			c.Redirect(http.StatusFound, fmt.Sprintf("http://%s:%d/user/login", config.App.Host, config.App.Port))
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
			c.Set("uID", claims.UID)
			c.Set("Username", claims.Username)
			c.Set("token", tokenString)
		} else {
			//util.RespInvalidToken(c)
			c.Redirect(http.StatusFound, fmt.Sprintf("http://%s:%d/user/login", config.App.Host, config.App.Port))
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证 token
func validateToken(tokenString string) (bool, *model.Claims, error) {
	// 验证 token 是否有效
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, nil, fmt.Errorf("token 解析错误: %v", err)
	}

	// 如果 token 合法且 claims 类型正确
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return true, claims, nil
	}

	return false, nil, fmt.Errorf("无效的 token")
}
