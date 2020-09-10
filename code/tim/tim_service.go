package tim

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

func RegisterSevice(r *gin.Engine) {
	r.GET("/tim", tool.AuthMiddleWare, func(c *gin.Context) {
		c.HTML(http.StatusOK, "tim_index.html", nil)
	})
}
