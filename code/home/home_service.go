package home

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)


//注册主页服务
func RegisterService(r *gin.Engine) {
	r.Static("/home/static", "template/home/static")

	r.GET("/", tool.AuthMiddleWare, func(c *gin.Context) {
		sess := sessions.Default(c)
		username := sess.Get("username")
		c.HTML(http.StatusOK, "home_index.html", gin.H{
			"headTip": fmt.Sprintf("你好：%s", username),
		})
	})
	r.GET("/home", tool.AuthMiddleWare, func(c *gin.Context) {
		sess := sessions.Default(c)
		username := sess.Get("username")
		c.HTML(http.StatusOK, "home_index.html", gin.H{
			"headTip": fmt.Sprintf("你好：%s", username),
		})
	})

	//搜索框输入内容后，跳转服务
	r.POST("/home", tool.AuthMiddleWare, func(c *gin.Context) {
		keyword := c.DefaultPostForm("keyword", "")
		if keyword != "" {
			c.Redirect(http.StatusFound, "http://www.baidu.com/s?wd=" + keyword)
			return
		} else {
			sess := sessions.Default(c)
			username := sess.Get("username")
			c.HTML(http.StatusOK, "home_index.html", gin.H{
				"headTip": fmt.Sprintf("你好：%s", username),
			})
			return
		}
	})
}
