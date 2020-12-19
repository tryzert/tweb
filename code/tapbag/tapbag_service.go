package tapbag

import (
	"github.com/gin-gonic/gin"
	"tweb/code/tool"
)

//注册我的云盘 TapBag 服务
func RegisterService(r *gin.Engine, srcPath string) {
	r.Static("tapbag/static", "template/tapbag/static")
	r.GET("/tapbag", tool.AuthMiddleWare, indexViewHandler())
	r.POST("/tapbag/api", apiHandler(srcPath))
	r.GET("/tapbag/api/online", onlineFileHandler(srcPath))
	r.POST("/tapbag/api/upload", uploadFileHandler(srcPath))
	r.GET("/tapbag/api/download", downloadFileHandler())
}
