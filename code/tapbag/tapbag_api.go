package tapbag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"tweb/code/tool"
)



func response(c *gin.Context, code int, tip string, data interface{}) {
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
			//requestFiles(c, srcPath, req.Data)
			if reflect.TypeOf(req.Data).Name() == "string" {
				data := reflect.ValueOf(req.Data).String()
				requestFiles(c, srcPath, data)
			} else {
				response(c, -1, "请求参数错误！", nil)
			}

		// 2001 new folder
		case 2001:
			if reflect.TypeOf(req.Data).Name() == "string" {
				data := reflect.ValueOf(req.Data).String()
				requestMakeNewFolder(c, srcPath, data)
			} else {
				response(c, -1, "请求参数错误！", nil)
			}
		// 2002 upload
		// 2003 download
		case 2003:
			if reflect.TypeOf(req.Data).Kind().String() == "slice" {
				if data, ok := req.Data.([]interface{}); ok {
					if len(data) == 0 {
						response(c, -1, "请求参数错误！", nil)
						return
					}
					requestDownload(c, srcPath, data)
					return
				}
			}
			response(c, -1, "请求参数错误！", nil)
		// 2004 move
		// 2005 rename
		case 2005:
			if reflect.TypeOf(req.Data).Kind().String() == "map"{
				if data, ok := req.Data.(map[string]interface{}); ok {
					oldpath, ok1 := data["oldpath"]
					newpath, ok2 := data["newpath"]
					if len(data) == 2 && ok1 && ok2 {
						requestRename(c, srcPath, filepath.Clean(fmt.Sprint(oldpath)), filepath.Clean(fmt.Sprint(newpath)))
						return
					}
				}
				response(c, -1, "请求参数错误！", nil)
			}
		// 2006 delete
		case 2006:
			if reflect.TypeOf(req.Data).Kind().String() == "slice" {
				//response(c, 2006, "success", nil)
				if data, ok := req.Data.([]interface{}); ok {
					//response(c, 2006, "success", data)
					if len(data) == 0 {
						response(c, -1, "请求参数错误！", nil)
						return
					}
					requestRemove(c, srcPath, data)
					return
				}
			}
			response(c, -1, "请求参数错误！", nil)
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
		response(c, -1,"文件夹已存在！", nil)
		return
	}
	err = os.Mkdir(fullpath, os.ModePerm)
	if err != nil {
		response(c, -1, "服务器创建文件夹时出错！", nil)
		return
	}
	response(c, 2001, "创建文件夹成功！", nil)
}


// oldName and newName are full path
func requestRename(c *gin.Context, srcPath, oldPath, newPath string) {
	opath := filepath.Join(srcPath, oldPath)
	npath := filepath.Join(srcPath, newPath)
	exist1, err1 := tool.FileExist(opath)
	exist2, err2 := tool.FileExist(npath)
	if err1 != nil || err2 != nil {
		response(c, -1, "服务器检索文件是否存在时发生未知错误！", nil)
		return
	}
	if !exist1 {
		response(c, -1,"无法对不存在的文件重命名！", nil)
		return
	}
	if exist2 {
		response(c, -1,"文件名重复了！", nil)
		return
	}
	if os.Rename(opath, npath) != nil {
		response(c, -1, "服务器对文件重命名时出错！", nil)
	} else {
		response(c, 2005, "文件重命名成功！！", nil)
	}
}


func requestRemove(c *gin.Context, srcPath string, data []interface{}) {
	var err error
	errorList := []string{}
	for _, relpath := range data {
		absPath := filepath.Join(srcPath, fmt.Sprint(relpath))
		err = os.RemoveAll(absPath)
		if err != nil {
			errorList = append(errorList, fmt.Sprint(relpath))
		}
	}
	if len(errorList) == len(data) {
		response(c, -1, "服务器所有删除操作全部失败！", errorList)
	} else if len(errorList) > 0 {
		response(c, -1, "服务器在删除这些文件时遇到错误，其余文件删除成功！", errorList)
	} else {
		response(c, 2006, "所有数据删除成功！", nil)
	}
}


func requestDownload(c *gin.Context, srcPath string, data []interface{}) {
	for _, relpath := range data {
		absPath := filepath.Join(srcPath, fmt.Sprint(relpath))
		fmt.Println(absPath)
	}
	response(c, 2003, "hhh", data)
}