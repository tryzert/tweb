package tapbag

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"tweb/code/tool"
)

//注册我的云盘 TapBag 服务
func RegisterService(r *gin.Engine, srcPath string) {
	r.Static("tapbag/static", "template/tapbag/static")
	r.GET("/tapbag", tool.AuthMiddleWare, func(c *gin.Context) {
		c.HTML(http.StatusOK, "tapbag_index.html", nil)
	})

	r.POST("/tapbag/api", apiHandler(srcPath))


	r.POST("/tapbag/api/upload", func(c *gin.Context) {
		//f, err := c.FormFile("file")
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
	})

	r.StaticFS("/tapbag/api/online", http.Dir(srcPath))
}


