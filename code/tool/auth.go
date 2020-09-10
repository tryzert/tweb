package tool

import (
	"github.com/gin-gonic/gin"
)


func AuthMiddleWare(c *gin.Context) {
	//sess := sessions.Default(c)
	//username := sess.Get("username")
	//userPassKey := sess.Get("userPassKey")
	////fmt.Println(username, userPassKey)
	//if UserPassKeyValidate(username, userPassKey) {
	//	//todo
	//	c.Next()
	//	return
	//} else {
	//	c.Redirect(http.StatusMovedPermanently, "/login")
	//	c.Abort()
	//	return
	//}
}