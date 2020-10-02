package login

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"tweb/code/tool"
)

func RegisterSevice(r *gin.Engine) {
	r.Static("/login/static", "template/login/static")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login_index.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		// username := c.PostForm("username")
		// password := c.PostForm("password")
		username := c.DefaultPostForm("username", "匿名")
		password := c.DefaultPostForm("password", "******")

		if tool.UserLoginValidate(username, password) { //用户名和密码都正确，跳转到首页
			sess := sessions.Default(c)
			sess.Set("username", username)
			userPassKey := tool.CreateRandPassKey(32)
			sess.Set("userPassKey", userPassKey)
			tool.Upks.Set(username, userPassKey, time.Hour * 2)
			sess.Save()
			c.Redirect(http.StatusFound, "/")
		} else { //用户名或密码错误
			c.HTML(http.StatusOK, "login_index.html", gin.H{
				"loginTip": "用户名或密码错误",
			})
		}
	})


	r.GET("/logout", func(c *gin.Context) {
		sess := sessions.Default(c)
		username, ok := sess.Get("username").(string)
		fmt.Println("this is logout, username:", username)
		if ok {
			fmt.Println("delete before", tool.Upks)
			tool.Upks.Delete(username)
			fmt.Println("delete after", tool.Upks)
		}
		sess.Clear()
		//sess.Save()
		c.Redirect(http.StatusFound, "/login")
	})
}




