package tool

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
)

/*
	跟文件相关的方法
*/


//判断文件或文件夹是否存在
// filepath 为包含路径的文件/文件夹的名称
func FileExist(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


//判断 路径是 文件还是文件夹
func IsFile(filepath string) bool {
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return !finfo.IsDir()
}


//实现 单个文件 下载功能
func downloadFile(c *gin.Context, filepath string) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	_, filename := path.Split(filepath)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Writer.Write(content)
}


//实现 整个文件夹 下载功能
func downloadDir(c *gin.Context, dirpath string) {
	//todo
}



//抽象的下载方法
func Download(c *gin.Context, path string) {
	if path == "" {
		return
	}
	if exist, _ := FileExist(path); exist {
		if IsFile(path) {
			downloadFile(c, path)
		} else {
			downloadDir(c, path)
		}
	}
}


//压缩文件
func ZipFile(filepath string) error {
	//todo
	return nil
}


//压缩文件夹
func ZipDir(dirpath string) error {
	//todo
	//ioutil.ReadDir()
	return nil
}