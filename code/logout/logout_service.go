package logout

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

//注册登出页服务
//这个服务不需要公开访问
func RegisterService(r *gin.Engine) {
	r.GET("/logout", func(c *gin.Context) {
		sess := sessions.Default(c)
		username, ok := sess.Get("username").(string)
		//fmt.Println("this is logout, username:", username)
		if ok {
			//fmt.Println("delete before", tool.Upks)
			tool.Upks.Delete(username)
			//fmt.Println("delete after", tool.Upks)
		}
		sess.Clear()
		sess.Save()
		c.Redirect(http.StatusFound, "/login")
	})
}
