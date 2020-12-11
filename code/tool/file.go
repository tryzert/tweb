package tool

import (
	"archive/zip"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
func downloadSingleFile(c *gin.Context, filepath string) {
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
func downloadMultiFiles(c *gin.Context, dirpath string) {
	//todo
}

//抽象的下载方法
func Download(c *gin.Context, path string) {
	if path == "" {
		return
	}
	if exist, _ := FileExist(path); exist {
		if IsFile(path) {
			downloadSingleFile(c, path)
		} else {
			downloadMultiFiles(c, path)
		}
	}
}

//压缩单文件
func ZipSingleFile(filepath, zipName string, compressed bool) error {
	d, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// 指定文件压缩方式 默认为 Store 方式 该方式不压缩文件 只是转换为zip保存
	if compressed {
		header.Method = zip.Deflate
	}
	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, f)
	return err
}


func ZipMultiFiles(files []string, dest string, compressed bool) error {
	pfs := make([]*os.File, len(files))
	for i, file := range files {
		fd, err := os.Open(file)
		if err != nil {
			return err
		}
		pfs[i] = fd
	}
	return compressFiles(pfs, dest, compressed)
}


//压缩文件夹
func compressFiles(files []*os.File, dest string, compressed bool) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w, compressed)
		if err != nil {
			return err
		}
	}
	return nil
}


func compress(file *os.File, prefix string, zw *zip.Writer, compressed bool) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = filepath.Join(prefix, info.Name())
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(filepath.Join(file.Name(), fi.Name()))
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw, compressed)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(prefix, header.Name)
		if compressed {
			header.Method = zip.Deflate
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
