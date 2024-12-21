package main

import (
	"livestreamall/api"
	"livestreamall/config"
	"livestreamall/dao"
)

func main() {
	// 加载配置
	config.LoadConfig()

	dao.InitDB()
	api.InitRouter()
}
