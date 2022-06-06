package file_op

import (
	"errors"
	"io/fs"
	"os"
)

type FileOp struct {
	file *os.File
	path string
}

func CreateFileOp(path string) *FileOp {
	return &FileOp{
		path: path,
	}
}

// Ready
// @author Tianyi
// @description 用于进行文件操作前的准备工作
func (fo *FileOp) Ready() (err error) {
	if fo.file == nil {
		if IsExists(fo.path) {
			fo.file, err = MustOpenFile(fo.path)
		} else {
			fo.file, err = CreateFile(fo.path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (fo *FileOp) Write(buf []byte) error {
	buf = append(buf, '\n')
	_, err := fo.file.Write(buf)
	return err
}

func (fo *FileOp) Close() error {
	err := fo.file.Close()
	return err
}

// IsExists
// @author Tianyi
// @description 判断路径是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

// IsPermission
// @author Tianyi
// description 判断文件是否有权限操作
func IsPermission(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrPermission)
}

// Mkdir
// @author Tianyi
// @description 创建一个目录
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile
// @author Tianyi
// @description 创建文件，先检查文件是否存在，存在就报错，不存在就创建
func CreateFile(path string) (*os.File, error) {
	exist := IsExists(path)
	if exist {
		return nil, errors.New("file already exists")
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// MustOpenFile
// @author Tianyi
// @description 直接打开文件，使用该方法的前提是确定文件一定存在
func MustOpenFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0666)
	return file, err
}
