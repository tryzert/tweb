package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//管理页服务
func RegisterService(r *gin.Engine, signal chan int) {
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_index.html", nil)
	})

	r.POST("/admin/api", func(c *gin.Context) {
		data := struct {
			Code int `json:"code"`
		}{}
		if c.ShouldBindJSON(&data) == nil && data.Code != 0 {
			// code == 100 : start
			// code == 101 : shutdown
			// code == 102 : reboot
			switch data.Code {
			case 100:
				signal <- 100
			case 101:
				signal <- 101
			case 102:
				signal <- 102
			}
		}
	})
}
