package tool

import "os"

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
