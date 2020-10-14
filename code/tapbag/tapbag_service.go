package tapbag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"tweb/code/tool"
)

//注册我的云盘 TapBag 服务
func RegisterService(r *gin.Engine, srcPath string) {
	r.Static("tapbag/static", "template/tapbag/static")
	r.GET("/tapbag", tool.AuthMiddleWare, func(c *gin.Context) {
		c.String(http.StatusOK, "hello tapbag!")
	})



	r.POST("/tapbag/api", func(c *gin.Context) {
		req := RequestContent{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		if req.Code != 100 {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		requestDirFullPath := filepath.Join(srcPath, req.Data)
		if requestDirFullPath == "" {
			requestDirFullPath = "./"
		}
		fs, err := ioutil.ReadDir(requestDirFullPath)

		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		sort.Slice(fs, func(i, j int) bool {
			//return fs[i].Name() > fs[j].Name()
			a, b := fs[i].IsDir(), fs[j].IsDir()
			if a && b {
				return fs[i].Name() < fs[j].Name()
			}
			if !a && !b {
				return fs[i].Name() < fs[j].Name()
			}
			if a && !b {
				return true
			}
			return false
		})
		files := []*File{}
		fmt.Println(req.Data, requestDirFullPath)
		for id, file := range fs {
			relpath := filepath.Join(req.Data, file.Name())
			if relpath == "" {
				relpath = "/"
			}
			fileinfo := &File{Id: id, Name: file.Name(), Relpath: relpath}
			if file.IsDir() {
				fileinfo.Type = "folder"
			} else {
				fileextname := strings.ToLower(filepath.Ext(file.Name()))
				if fileextname == ".zip" || fileextname == ".gz"  || fileextname == ".tar" || fileextname == ".xz"{
					fileinfo.Type = "archive"
				} else if fileextname == ".mp3" {
					fileinfo.Type = "audio"
				} else if fileextname == ".doc" {
					fileinfo.Type = "doc"
				} else if fileextname == ".png" || fileextname == ".jpg" || fileextname == ".svg"{
					fileinfo.Type = "image"
				} else if fileextname == ".pdf" {
					fileinfo.Type = "pdf"
				} else if fileextname == ".ppt" {
					fileinfo.Type = "ppt"
				} else if fileextname == ".psd" {
					fileinfo.Type = "psd"
				} else if fileextname == ".txt" {
					fileinfo.Type = "text"
				} else if fileextname == ".mp4" || fileextname == ".avi" {
					fileinfo.Type = "video"
				} else if fileextname == ".xls" {
					fileinfo.Type = "xls"
				} else {
					fileinfo.Type = "file"
				}
			}
			files = append(files, fileinfo)
		}
		response(c, 100, "请求成功！", files)
	})

	r.POST("/tapbag/api/online", func(c *gin.Context) {
		req := RequestContent{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		if req.Code != 101 {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		onlineUrl := filepath.Join(srcPath, req.Data)
		c.JSON(http.StatusOK, gin.H{
			"code": 101,
			"tip": "请求成功！",
			"data": onlineUrl,
		})
	})
	r.StaticFS("/tapbag/api/online", http.Dir(srcPath))
}




type RequestContent struct {
	Code int `json:"code"`
	Data string `json:"data"`
}

type ResponseContent struct {
	Code int `json:"code"`
	Tip string `json:"tip"`
	Data []*File `json:"data"`
}


type File struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Relpath string `json:"relpath"`
}


func response(c *gin.Context, code int, tip string, data []*File) {
	c.JSON(http.StatusOK, ResponseContent{
		Code: code,
		Tip:  tip,
		Data: data,
	})
}