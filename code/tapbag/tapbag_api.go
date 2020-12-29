package tapbag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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
	filemanager := NewFileManager(srcPath)
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
				requestFiles(c, filemanager, srcPath, data)
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
		case 2004:
			if reflect.TypeOf(req.Data).Kind().String() == "map" {
				if data, ok := req.Data.(map[string]interface{}); ok {
					fromPath, ok1 := data["from"]
					toPath, ok2 := data["to"]
					moveList, ok3 := data["moveList"]
					if len(data) == 3 && ok1 && ok2 && ok3 {
						if reflect.TypeOf(moveList).Kind().String() == "slice" {
							if mList, ok4 := moveList.([]interface{}); ok4 {
								requestMove(c, srcPath, fromPath, toPath, mList)
								return
							}

						}
					}
				}
			}
			response(c, -1, "请求参数错误！", nil)
		// 2005 rename
		case 2005:
			if reflect.TypeOf(req.Data).Kind().String() == "map" {
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
					requestDelete(c, srcPath, data)
					return
				}
			}
			response(c, -1, "请求参数错误！", nil)
		// 2007 request folders
		case 2007:
			if reflect.TypeOf(req.Data).Name() == "string" {
				data := reflect.ValueOf(req.Data).String()
				requestFolders(c, srcPath, data)
			} else {
				response(c, -1, "请求参数错误！", nil)
			}
		default:
			response(c, -1, "请求参数错误！", nil)
			return
		}

	}
}

// 2000
func requestFiles(c *gin.Context, fm *FileManager, srcPath, rdata string) {
	requestDirFullPath := filepath.Join(srcPath, rdata)
	if requestDirFullPath == "" {
		requestDirFullPath = "./"
	}
	fs, err := ioutil.ReadDir(requestDirFullPath)

	if err != nil {
		response(c, -1, "请求参数错误！", nil)
		return
	}
	//sort.Slice(fs, func(i, j int) bool {
	//	a, b := fs[i].IsDir(), fs[j].IsDir()
	//	if a && b {
	//		return fs[i].Name() < fs[j].Name()
	//	}
	//	if !a && !b {
	//		return fs[i].Name() < fs[j].Name()
	//	}
	//	if a && !b {
	//		return true
	//	}
	//	return false
	//})
	files := []*File{}
	for id, fileInfo := range fs {
		relpath := filepath.Join(rdata, fileInfo.Name())
		if relpath == "" {
			relpath = "/"
		}
		fileItem := &File{Id: id, Name: fileInfo.Name(), Relpath: relpath}
		fileItem.Type = fm.getType(fileInfo)
		fileItem.Openable = fm.fileOpenable(fileInfo)
		files = append(files, fileItem)
	}
	response(c, 2000, "请求数据成功！", files)
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
		response(c, -1, "无法对不存在的文件重命名！", nil)
		return
	}
	if exist2 {
		response(c, -1, "文件名重复了！", nil)
		return
	}
	if os.Rename(opath, npath) != nil {
		response(c, -1, "服务器对文件重命名时出错！", nil)
	} else {
		response(c, 2005, "文件重命名成功！！", nil)
	}
}

func requestDelete(c *gin.Context, srcPath string, data []interface{}) {
	// srcPath/.tweb/recycleBin
	if len(data) == 0 {
		response(c, -1, "请求参数错误！", nil)
		return
	}
	destPath := filepath.Join(srcPath, ".tweb/recycleBin")
	if exist, err := tool.FileExist(destPath); !exist || err != nil {
		if os.MkdirAll(destPath, os.ModePerm) != nil {
			panic("srcPath下的 .tweb/recycleBin 文件夹创建失败！")
		}
	}
	var err error
	errorList := []string{}
	successList := []string{}
	for _, relpath := range data {
		absPath := filepath.Join(srcPath, fmt.Sprint(relpath))
		baseName := filepath.Base(absPath)
		// todo 当回收站存在同名文件，删除文件会失败
		err = os.Rename(absPath, filepath.Join(destPath, baseName))
		if err != nil {
			errorList = append(errorList, fmt.Sprint(relpath))
		} else {
			successList = append(successList, fmt.Sprint(relpath))
		}
	}
	go addDeleteRecord(successList)
	if len(errorList) == len(data) {
		response(c, -1, "服务器所有删除操作全部失败！", errorList)
	} else if len(errorList) > 0 {
		response(c, -1, "服务器在删除这些文件时遇到错误，其余文件删除成功！", errorList)
	} else {
		response(c, 2006, "所有数据删除成功！", nil)
	}
}

func requestDownload(c *gin.Context, srcPath string, data []interface{}) {
	//files := []string{}
	//for _, relpath := range data {
	//	absPath := filepath.Join(srcPath, fmt.Sprint(relpath))
	//	if exist, _ := tool.FileExist(absPath); exist {
	//		files = append(files, absPath)
	//	} else {
	//		response(c, -2003, "服务器检测到要下载的部分文件已经不存在！", nil)
	//		return
	//	}
	//}
	//fileKey := tool.Encryption(filepath.Join(files...))
	//if Fman.Exist(fileKey) {
	//	response(c, 2003, "服务器准备中，下载即将开始！", fileKey)
	//	return
	//}
	//fileKey = Fman.Put(files)
	//response(c, 2003, "服务器准备中，下载即将开始！", fileKey)
}

func requestMove(c *gin.Context, srcPath string, fromPath, toPath interface{}, moveList []interface{}) {
	from := fmt.Sprint(fromPath)
	to := fmt.Sprint(toPath)
	if from == "" || to == "" || len(moveList) == 0 {
		response(c, -1, "请求参数错误！", nil)
		return
	}
	if from == to {
		response(c, -2004, "文件已在此路径下！", nil)
		return
	}
	for _, fp := range moveList {
		fpath := fmt.Sprint(fp)
		if strings.HasPrefix(to, fpath) {
			response(c, -2004, "路径冲突：文件夹不能移到子文件夹或本身！", nil)
			return
		}
	}
	destPath := filepath.Join(srcPath, to)
	errCount := 0
	for _, fp := range moveList {
		fpath := filepath.Join(srcPath, fmt.Sprint(fp))
		fBaseName := filepath.Base(fpath)
		err := os.Rename(fpath, filepath.Join(destPath, fBaseName))
		if err != nil {
			errCount++
		}
	}
	if errCount == len(moveList) {
		response(c, -2004, "所有文件均移动失败！", nil)
	} else if errCount > 0 {
		response(c, -2004, "部分文件移动失败！", nil)
	} else {
		response(c, 2004, "文件移动成功！", nil)
	}
}

func requestFolders(c *gin.Context, srcPath, relpath string) {
	fullpath := filepath.Join(srcPath, relpath)
	if exist, err := tool.FileExist(fullpath); exist && err == nil {
		folders := make([]*Folder, 0)
		infos, err := ioutil.ReadDir(fullpath)
		if err != nil {
			response(c, -1, "请求参数错误！", nil)
			return
		}
		for _, info := range infos {
			if info.IsDir() {
				folder := &Folder{Name: info.Name(), Src: filepath.Join(relpath, info.Name())}
				if folderHasFolder(srcPath, folder.Src) {
					folder.HasChildren = true
				} else {
					folder.HasChildren = false
				}
				folders = append(folders, folder)
			}
		}
		response(c, 2007, "请求数据成功！", folders)
	} else {
		response(c, -1, "请求参数错误！", nil)
	}
}

func folderHasFolder(srcPath, path string) bool {
	infos, err := ioutil.ReadDir(filepath.Join(srcPath, path))
	if err != nil {
		return false
	}
	for _, info := range infos {
		if info.IsDir() {
			return true
		}
	}
	return false
}
