package tapbag

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
		requestDirFullPath := filepath.Join(srcPath, req.Data)
		if requestDirFullPath == "" {
			requestDirFullPath = "."
		}
		fs, err := ioutil.ReadDir(requestDirFullPath)
		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		files := []*File{}
		for id, file := range fs {
			fileinfo := &File{Id: id, FileName: file.Name()}
			if file.IsDir() {
				fileinfo.FileType = "folder"
			} else {
				fileextname := strings.ToLower(filepath.Ext(file.Name()))
				if fileextname == ".zip" || fileextname == ".gz"  || fileextname == ".tar" || fileextname == ".xz"{
					fileinfo.FileType = "archive"
				} else if fileextname == ".mp3" {
					fileinfo.FileType = "audio"
				} else if fileextname == ".doc" {
					fileinfo.FileType = "doc"
				} else if fileextname == ".png" || fileextname == ".jpg" {
					fileinfo.FileType = "image"
				} else if fileextname == ".pdf" {
					fileinfo.FileType = "pdf"
				} else if fileextname == ".ppt" {
					fileinfo.FileType = "ppt"
				} else if fileextname == ".psd" {
					fileinfo.FileType = "psd"
				} else if fileextname == ".txt" {
					fileinfo.FileType = "text"
				} else if fileextname == ".mp4" || fileextname == ".avi" {
					fileinfo.FileType = "video"
				} else if fileextname == ".xls" {
					fileinfo.FileType = "xls"
				} else {
					fileinfo.FileType = "file"
				}
			}
			files = append(files, fileinfo)
		}
		response(c, 100, "请求成功！", files)
	})
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
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	//FullPath string `json:"full_path"`
}


func response(c *gin.Context, code int, tip string, data []*File) {
	c.JSON(http.StatusOK, ResponseContent{
		Code: code,
		Tip:  tip,
		Data: data,
	})
}