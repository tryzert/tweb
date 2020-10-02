package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterService(r *gin.Engine) {
	r.GET("admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_index.html", nil)
	})
}
