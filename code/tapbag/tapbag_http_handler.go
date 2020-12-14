package tapbag

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"tweb/code/tool"
)


func indexViewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "tapbag_index.html", nil)
	}
}

func onlineFileHandler(srcPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.DefaultQuery("path", "")
		fullPath := filepath.Join(srcPath, path)
		if tool.IsFile(fullPath) {
			c.File(fullPath)
		} else {
			c.HTML(http.StatusNotFound, "notFound404_index.html", nil)
		}
	}
}

func downloadFileHandler(srcPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullpath := filepath.Join(srcPath, c.DefaultQuery("path", ""))
		if exist, err := tool.FileExist(fullpath); exist && err == nil {
			if tool.IsFile(fullpath) {
				tool.Download(c, fullpath, tool.FILE, true)
				return
			}
		}
		c.HTML(http.StatusNotFound, "notFound404_index.html", nil)
	}
}

func uploadFileHandler(srcPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			response(c, -1, "上传文件出错！", nil)
			return
		}
		savePath := filepath.Join(srcPath, c.DefaultPostForm("path", ""))
		files := form.File["files"]
		errCount := 0
		for _, file := range files {
			fpath := filepath.Join(savePath, file.Filename)
			err = c.SaveUploadedFile(file, fpath)
			if err != nil {
				errCount++
			}
		}
		if errCount == len(files) {
			response(c, -1, "所有文件均上传失败！", nil)
		} else if errCount == 0 {
			response(c, 2002, "文件全部上传成功！", nil)
		} else {
			response(c, -1, "部分文件上传失败！", nil)
		}
	}
}
