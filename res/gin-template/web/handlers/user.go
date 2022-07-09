package handlers

import (
	"github.com/IfanTsai/go-lib/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	middlewares.SetResp(c, "hello world")
}
