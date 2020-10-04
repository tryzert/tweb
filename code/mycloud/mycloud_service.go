package mycloud

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

//注册我的云盘 mycloud 服务
func RegisterService(r *gin.Engine) {
	r.GET("/mycloud", tool.AuthMiddleWare, func(c *gin.Context) {
		c.String(http.StatusOK, "hello mycloud!")
	})

	r.GET("/mycloud/download", tool.AuthMiddleWare, func(c *gin.Context) {
		filepath := c.DefaultQuery("filepath", "")
		tool.Download(c, filepath)
	})
}
