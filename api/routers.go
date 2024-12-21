package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"livestreamall/config"
	"net/http"
)

func InitRouter() {

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.LoadHTMLGlob("html/*")

	r.Static("/static", "./static")

	r.GET("/home", HomePage)

	user := r.Group("/user")
	{
		user.GET("/login", LoginPage)                    // 用户注册登录页
		user.POST("/register", Register)                 // 注册API
		user.POST("/login", Login)                       // 登录API
		user.GET("/info", AuthMiddleware(), GetUserInfo) // 获取用户信息
	}

	r.POST("/auth/publish", PublishAuth)            // 推流认证
	r.POST("/auth/stop_publish", StopPublish)       // 停止推流
	r.GET("/stream/:stream_name", GetPullStreamURL) // 获取拉流相关地址
	//r.GET("/stream/:stream_id", StreamForwarding) // 配置流转发路由

	live := r.Group("/live")
	{

		live.GET("/play", Liveroom)                      // 观看直播页
		live.GET("/start", StartLivePage)                // 开始直播页
		live.POST("/start", AuthMiddleware(), StartLive) // 开始直播,获取推流地址,获取推流密钥
		live.GET("/live_rooms", GetLiveRooms)            // 获取直播间列表
	}

	r.Run(fmt.Sprintf(":%d", config.App.Port))
}

// CORSMiddleware 自定义CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
