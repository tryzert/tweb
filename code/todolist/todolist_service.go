package todolist

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

//注册待办事项 todolist 服务
func RegisterService(r *gin.Engine) {
	r.Static("/todolist/static", "template/todolist/static")
	r.GET("/todolist", tool.AuthMiddleWare, func(c *gin.Context) {
		c.HTML(http.StatusOK, "todolist_index.html", nil)
	})

	r.POST("/todolist/api", apiHandler)
}
