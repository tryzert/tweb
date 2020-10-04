package todolist

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

//注册待办事项 todolist 服务
// to do
func RegisterService(r *gin.Engine) {
	//r.Static("/home/static", "template/home/static")
	//r.StaticFile("/todolist", "template/home/todolist_index.html")
	r.GET("/todolist", tool.AuthMiddleWare, func(c *gin.Context) {
		c.HTML(http.StatusOK, "todolist_index.html", nil)
	})

	r.POST("/todolist/api", apiHandler)
}
