package tapbag

import (
	"fmt"
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
		if exist, _ := tool.FileExist(fullPath); exist {
			if tool.IsFile(fullPath) {
				c.File(fullPath)
				return
			}

		}
		c.HTML(http.StatusNotFound, "notFound404_index.html", nil)
	}
}

func downloadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileKey := c.DefaultQuery("filekey", "")
		fmt.Println(fileKey)
		//if fileKey == "" || !Fman.Exist(fileKey){
		//	c.HTML(http.StatusNotFound, "notFound404_index.html", nil)
		//	return
		//}
		//Download(c, fileKey)
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
