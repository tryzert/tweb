package tapbag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"tweb/code/tool"
)



func response(c *gin.Context, code int, tip string, data []*File) {
	c.JSON(http.StatusOK, ResponseContent{
		Code: code,
		Tip:  tip,
		Data: data,
	})
}


func apiHandler(srcPath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := RequestContent{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}

		switch req.Code {
		//request files
		case 2000:
			requestFiles(c, srcPath, req.Data)
		// 2001 new folder
		case 2001:
			requestMakeNewFolder(c, srcPath, req.Data)
		// 2002 upload
		// 2003 download
		// 2004 move
		// 2005 rename
		// 2006 delete
		default:
			response(c, -1, "请求参数错误！", nil)
			return
		}

	}
}


// 2000
func requestFiles(c *gin.Context, srcPath string, rdata string) {
	requestDirFullPath := filepath.Join(srcPath, rdata)
	if requestDirFullPath == "" {
		requestDirFullPath = "./"
	}
	fs, err := ioutil.ReadDir(requestDirFullPath)

	if err != nil {
		response(c, -1, "请求参数错误！", nil)
		return
	}
	sort.Slice(fs, func(i, j int) bool {
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
	fmt.Println(rdata, requestDirFullPath)
	for id, file := range fs {
		relpath := filepath.Join(rdata, file.Name())
		if relpath == "" {
			relpath = "/"
		}
		fileinfo := &File{Id: id, Name: file.Name(), Relpath: relpath}
		if file.IsDir() {
			fileinfo.Type = "folder"
			fileinfo.Openable = true
		} else {
			fileextname := strings.ToLower(filepath.Ext(file.Name()))
			if fileextname == ".zip" || fileextname == ".gz"  || fileextname == ".tar" || fileextname == ".xz"{
				fileinfo.Type = "archive"
				fileinfo.Openable = false
			} else if fileextname == ".mp3" {
				fileinfo.Type = "audio"
				fileinfo.Openable = true
			} else if fileextname == ".doc" {
				fileinfo.Type = "doc"
				fileinfo.Openable = false
			} else if fileextname == ".png" || fileextname == ".jpg" || fileextname == ".svg"{
				fileinfo.Type = "image"
				fileinfo.Openable = true
			} else if fileextname == ".pdf" {
				fileinfo.Type = "pdf"
				fileinfo.Openable = false
			} else if fileextname == ".ppt" {
				fileinfo.Type = "ppt"
				fileinfo.Openable = false
			} else if fileextname == ".psd" {
				fileinfo.Type = "psd"
				fileinfo.Openable = false
			} else if fileextname == ".txt" {
				fileinfo.Type = "text"
				fileinfo.Openable = false
			} else if fileextname == ".mp4" || fileextname == ".avi" {
				fileinfo.Type = "video"
				fileinfo.Openable = true
			} else if fileextname == ".xls" {
				fileinfo.Type = "xls"
				fileinfo.Openable = false
			} else {
				fileinfo.Type = "file"
				fileinfo.Openable = false
			}
		}
		files = append(files, fileinfo)
	}
	response(c, 2000, "请求成功！", files)
}
func requestMakeNewFolder(c *gin.Context, srcPath, relpath string) {
	fullpath := filepath.Join(srcPath, filepath.Clean(relpath))
	exist, err := tool.FileExist(fullpath)
	if err != nil {
		response(c, -1, "服务器检查文件夹是否存在时出错！", nil)
		return
	}
	if exist {
		response(c, -1, "文件夹已存在！", nil)
		return
	}
	err = os.Mkdir(fullpath, os.ModePerm)
	if err != nil {
		response(c, -1, "服务器创建文件夹时出错！", nil)
		return
	}
	response(c, 2001, "创建文件夹成功！", nil)
}