package todolist

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
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
	case 1000: // all
		data, err := queryAllTasks()
		if err != nil {
			response(c, -1000, "服务器读取数据时出错！", "")
		} else if data == nil {
			response(c, 1000, "查询结果为空！", "")
		} else {
			response(c, 1000, "请求数据成功！", data)
		}
	case 1001: // update task status
		data := new(struct {
			Id   int
			Done int
		})
		if reflect.TypeOf(req.Data).Kind().String() != "map" {
			response(c, -1, "请求参数格式错误！", "")
			return
		}
		if mp, ok := req.Data.(map[string]interface{}); ok {
			id, ok1 := mp["id"]
			done, ok2 := mp["done"]
			if len(mp) == 2 && ok1 && ok2 {
				var err1, err2 error
				data.Id, err1 = strconv.Atoi(fmt.Sprint(id))
				data.Done, err2 = strconv.Atoi(fmt.Sprint(done))
				if err1 == nil && err2 == nil {
					if updateActiveTaskStatus(data.Id, data.Done) {
						response(c, 1001, "更新数据成功！", "")
						return
					} else {
						response(c, -1001, "服务器更新数据时出错！", "")
						return
					}
				}
			}
		}
		response(c, -1, "请求参数格式错误！", "")

	case 1002: // update task content
		data := new(struct {
			Id      int
			Content string
		})
		if reflect.TypeOf(req.Data).Kind().String() != "map" {
			response(c, -1, "请求参数格式错误！", "")
			return
		}
		if mp, ok := req.Data.(map[string]interface{}); ok {
			id, ok1 := mp["id"]
			content, ok2 := mp["content"]
			if len(mp) == 2 && ok1 && ok2 {
				var err error
				data.Content = fmt.Sprint(content)
				data.Id, err = strconv.Atoi(fmt.Sprint(id))
				if err == nil {
					task, tag := dataHandler(data.Content)
					if dDate, dTime, ok := updateActiveTaskContent(data.Id, task, tag); ok {
						dt := make(map[string]string)
						dt["date"] = dDate
						dt["time"] = dTime
						dt["task"] = task
						dt["tag"] = tag
						response(c, 1002, "更新数据成功！", dt)
						return
					} else {
						response(c, -1002, "服务器更新数据时出错！", "")
						return
					}
				}
			}
		}
		response(c, -1, "请求参数格式错误！", "")
	case 1003: //delete active task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if deleteActiveTask(id) {
				data, _ := tool.TimeFormat(time.Now())
				response(c, 1003, "更新数据成功！", data)
			} else {
				response(c, -1003, "服务器更新数据时出错！", "")
			}
		}
	case 1004: //recover history task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if recoverHistoryTask(id) {
				response(c, 1004, "更新数据成功！", "")
			} else {
				response(c, -1004, "服务器更新数据时出错！", "")
			}
		}
	case 1005: // delete history task
		res := fmt.Sprint(req.Data)
		if id, err := strconv.Atoi(res); err != nil {
			response(c, -1, "请求参数格式错误！", "")
		} else {
			if deleteHistoryTask(id) {
				response(c, 1005, "更新数据成功！", "")
			} else {
				response(c, -1005, "服务器更新数据时出错！", "")
			}
		}
	case 1006: // add new active task
		res := fmt.Sprint(req.Data)
		if data, ok := addActiveTask(res); ok {
			response(c, 1006, "添加数据成功！", data)
		} else {
			response(c, -1006, "服务器更新数据时出错！", "")
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
