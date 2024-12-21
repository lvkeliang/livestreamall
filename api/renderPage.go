package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Liveroom(c *gin.Context) {
	c.HTML(http.StatusOK, "liveroom.html", gin.H{})
}

func StartLivePage(c *gin.Context) {
	c.HTML(http.StatusOK, "startlive.html", gin.H{})
}

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}
