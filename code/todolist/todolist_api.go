package todolist

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	主要处理前端过来的各种api请求
*/
func apiHandler(c *gin.Context) {
	//rcode := c.DefaultPostForm("requestCode", "-1")
	//fmt.Println(rcode)
	//switch rcode {
	//case "-1":
	//	response(c, -1, "参数错误！", "")
	//case "100":
	//	data, err := queryAll()
	//	if err != nil {
	//		response(c, 100, "获取数据失败！", "")
	//		return
	//	}
	//	response(c, 100, "获取数据成功！", data)
	//default:
	//	response(c, 300, "其他", "300 code")
	//}
	qy, _ := queryAll()
	response(c, 100, "success", qy)
}

func response(c *gin.Context, responseCode int, tip string, data interface{}) {
	c.JSON(http.StatusOK, ResponseContent{
		ResponseCode: responseCode,
		Tip:          tip,
		Data:         data,
	})
}
