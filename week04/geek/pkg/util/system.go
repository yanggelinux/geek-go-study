package util

import (
	"io/ioutil"
	"os"
	"path"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//删除目录下所有内容
func RemoveAndCreateDir(path string) error {
	//var err error
	if IsDir(path) {
		if Exists(path) {
			err := os.RemoveAll(path)
			if err != nil {

				return err
			}
			err = os.Mkdir(path, 0744)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//清除目录下的第一层的文件
func RemoveDirFiles(dirPath string) error {
	if IsDir(dirPath) {
		if Exists(dirPath) {
			rd, _ := ioutil.ReadDir(dirPath)
			for _, item := range rd {
				name := item.Name()
				filePath := path.Join(dirPath, name)
				if IsFile(filePath) {
					err := os.Remove(filePath)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
