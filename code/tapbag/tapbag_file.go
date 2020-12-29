package tapbag

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"tweb/code/tool"
)

//实现 单个文件 下载功能
func downloadSingleFile(c *gin.Context, abspath string) {
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	filename := filepath.Base(abspath)
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	// c.Writer.Header().Add("Content-Length", 1024)
	c.File(abspath)
}

//实现 多个文件包括文件夹 下载功能
func downloadMultiFiles(c *gin.Context, zipName string, filepaths []string) {
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+zipName+"\"")
	tool.ZipFilesToStream(filepaths, c.Writer, false)
}

////抽象的下载方法
//func Download(c *gin.Context, fileKey string) {
//	downloadSingleFile(c, fileKey)
//}
//
//// file system center
//type FileManager struct {
//	FileSystem map[string]*FileItem
//	Lock       sync.RWMutex
//}
//
//func NewFileManager() *FileManager {
//	return &FileManager{
//		FileSystem: make(map[string]*FileItem),
//	}
//}
//
//func (this *FileManager) Exist(fileKey string) bool {
//	this.Lock.RLock()
//	defer this.Lock.RUnlock()
//	_, ok := this.FileSystem[fileKey]
//	return ok
//}
//
//func (this *FileManager) Get(fileKey string) *FileItem {
//	this.Lock.Lock()
//	defer this.Lock.Unlock()
//	if fit, ok := this.FileSystem[fileKey]; ok {
//		fit.VisitHot += 1
//		fit.LastVisitTime = time.Now()
//		return fit
//	}
//	return nil
//}
//
//func (this *FileManager) Put(files []string) string {
//	this.Lock.Lock()
//	defer this.Lock.Unlock()
//	fileKey := tool.Encryption(filepath.Join(files...))
//	fit := &FileItem{
//		Paths:    files,
//		Size:     0,
//		VisitHot: 0,
//	}
//	for _, file := range files {
//		if size, err := tool.CalculateFileSize(file); err != nil {
//			return ""
//		} else {
//			fit.Size += size
//		}
//	}
//	this.FileSystem[fileKey] = fit
//	return fileKey
//}
//
////定时清理map
//func (this *FileManager) Check() {
//	for {
//		if len(this.FileSystem) > 100 {
//			this.Lock.Lock()
//
//			for k, v := range this.FileSystem {
//				if v.LastVisitTime.Add(time.Hour*24).Before(time.Now()) || v.VisitHot < 3 {
//					delete(this.FileSystem, k)
//				}
//			}
//
//			this.Lock.Unlock()
//		}
//
//		time.Sleep(time.Hour * 3)
//	}
//}
