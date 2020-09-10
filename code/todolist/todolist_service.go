package todolist
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tweb/code/tool"
)

func RegisterSevice(r *gin.Engine) {
	//r.Static("/home/static", "template/home/static")
	//r.StaticFile("/todolist", "template/home/todolist_index.html")
	r.GET("/todolist", tool.AuthMiddleWare, func(c *gin.Context) {
		c.HTML(http.StatusOK, "todolist_index.html", nil)
	})
}
