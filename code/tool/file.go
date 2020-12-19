package tool

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
	跟文件相关的方法
*/

const (
	STREAM = iota
	FILE
)

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

//压缩多个文件
func ZipFilesToFile(files []string, dest string, compressed bool) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	return ZipFilesToStream(files, zipFile, compressed)
}

//压缩文件到流
func ZipFilesToStream(files []string, stream io.Writer, compressed bool) error {
	zw := zip.NewWriter(stream)
	defer zw.Close()

	for _, path := range files {
		err := filepath.Walk(path, func(root string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = strings.TrimPrefix(root, filepath.Dir(path)+"/")
			if info.IsDir() {
				header.Name += "/"
			} else {
				if compressed {
					header.Method = zip.Deflate
				}
				reader, err := os.Open(root)
				if err != nil {
					return err
				}
				defer reader.Close()
				writer, err := zw.CreateHeader(header)
				if err != nil {
					return err
				}
				_, err = io.Copy(writer, reader)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}


//计算文件或文件夹总大小
func CalculateFileSize(abspath string) (int64, error) {
	info, err := os.Stat(abspath)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		var size int64 = 0
		err = filepath.Walk(abspath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
		return size, nil
	}
	return info.Size(), nil
}