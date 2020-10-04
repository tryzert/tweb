package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//管理页服务
func RegisterService(r *gin.Engine) {
	r.GET("admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_index.html", nil)
	})
}
