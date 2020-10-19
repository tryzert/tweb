package tapbag

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

//注册我的云盘 TapBag 服务
func RegisterService(r *gin.Engine, srcPath string) {
	r.Static("tapbag/static", "template/tapbag/static")
	r.GET("/tapbag", tool.AuthMiddleWare, func(c *gin.Context) {
		//c.String(http.StatusOK, "hello tapbag!")
		c.HTML(http.StatusOK, "tapbag_index.html", nil)
	})



	r.POST("/tapbag/api", apiHandler(srcPath))

	r.StaticFS("/tapbag/api/online", http.Dir(srcPath))
}


