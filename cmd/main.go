package main

import (
	"livestreamall/api"
	"livestreamall/dao"
)

func main() {
	dao.InitDB()
	api.InitRouter()
}
