package mycloud

import (
	"github.com/gin-gonic/gin"
	"tweb/code/tool"
)

func RegisterService(r *gin.Engine) {
	r.GET("/mycloud", tool.AuthMiddleWare, func(c *gin.Context) {
		c.File("/home/maple/workspace/tweb/template/home/static/img/shimei.png")
	})
}