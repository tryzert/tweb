package todolist

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
	"time"
	"tweb/code/tool"
)

/*
	主要处理前端过来的各种api请求
*/
func apiHandler(c *gin.Context) {
	req := RequestContent{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println(err)
		response(c, -1, "请求参数格式错误！", "")
		return
	}
	switch req.Code {
	case 100: // all
		data, err := queryAllTasks()
		if err != nil {
			response(c, req.Code, "服务器读取数据出错！", "")
		} else if data == nil {
			response(c, req.Code, "查询不到数据！", "")
		} else {
			response(c, req.Code, "请求数据成功！", data)
		}
	case 101: // update task status
		data := new(struct {
			Id   int
			Done int
		})
		if mapstructure.Decode(req.Data, data) != nil {
			response(c, -1, "请求参数格式错误！", "")
			return
		}
		if updateActiveTaskStatus(data.Id, data.Done) {
			response(c, 101, "更新数据成功！", "")
		} else {
			response(c, 101, "更新数据失败！", "")
		}
	case 102: // update task content
		data := new(struct {
			Id int
			//Task string
			//Tag string
			Content string
		})
		if mapstructure.Decode(req.Data, data) != nil {
			response(c, -1, "请求参数格式错误！", "")
			return
		}
		task, tag := dataHandler(data.Content)
		if updateActiveTaskContent(data.Id, task, tag) {
			dt := make(map[string]string)
			dt["task"] = task
			dt["tag"] = tag
			response(c, 102, "更新数据成功！", dt)
		} else {
			response(c, 102, "更新数据失败！", "")
		}
	case 103: //delete active task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if deleteActiveTask(id) {
				data, _ := tool.TimeFormat(time.Now())
				response(c, 103, "更新数据成功！", data)
			} else {
				response(c, 103, "更新数据失败！", "")
			}
		}
	case 104: //recover history task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if recoverHistoryTask(id) {
				response(c, 104, "更新数据成功！", "")
			} else {
				response(c, 104, "更新数据失败！", "")
			}
		}
	case 105: // delete history task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if deleteHistoryTask(id) {
				response(c, 105, "更新数据成功！", "")
			} else {
				response(c, 105, "更新数据失败！", "")
			}
		}
	case 106: // add new active task
		res := fmt.Sprint(req.Data)
		if data, ok := addActiveTask(res); ok {
			response(c, 106, "添加数据成功！", data)
		} else {
			response(c, 106, "添加数据失败！", "")
		}
	default:
		response(c, -1, "请求参数格式错误！", "")
	}
}

func response(c *gin.Context, code int, tip string, data interface{}) {
	c.JSON(http.StatusOK, ResponseContent{
		Code: code,
		Tip:  tip,
		Data: data,
	})
}
