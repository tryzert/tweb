package home

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

func RegisterService(r *gin.Engine) {
	r.Static("/home/static", "template/home/static")

	r.GET("/", tool.AuthMiddleWare, func(c *gin.Context) {
		sess := sessions.Default(c)
		username := sess.Get("username")
		c.HTML(http.StatusOK, "home_index.html", gin.H{
			"headTip": fmt.Sprintf("你好呀：%s", username),
		})
	})
	r.GET("/home", tool.AuthMiddleWare, func(c *gin.Context) {
		sess := sessions.Default(c)
		username := sess.Get("username")
		c.HTML(http.StatusOK, "home_index.html", gin.H{
			"headTip": fmt.Sprintf("你好呀：%s", username),
		})
	})
}
