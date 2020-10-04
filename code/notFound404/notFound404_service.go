package notFound404

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//注册404页面
func RegisterService(r *gin.Engine) {
	r.Static("/notFound404/static", "template/notFound404/static")
	//r.LoadHTMLFiles("./template/notFound404/notFound404_index.html", "./template/home/home_index.html")
	//r.LoadHTMLGlob("./template/notFound404/notFound404_index.html")

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "notFound404_index.html", nil)
	})
}
