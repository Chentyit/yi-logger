package file_op

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FileOp struct {
	file       *os.File
	isOpen     bool // 用于判断是否可以进行操作
	maxSize    int  // 以 MB 为单位
	curSize    int
	maxAge     int // 以天为单位
	maxBackups int
	curDate    time.Time
	path       string
}

func CreateFileOp(path string, maxSize int, maxAge int, maxBackups int) *FileOp {
	return &FileOp{
		path:       path,
		isOpen:     false,
		curSize:    0,
		maxSize:    maxSize,
		maxAge:     maxAge,
		maxBackups: maxBackups,
	}
}

// ready
// @author Tianyi
// @description 用于进行文件操作前的准备工作
func (fo *FileOp) ready() (err error) {
	if fo.file == nil {
		if IsExists(fo.path) {
			fo.file, err = MustOpenFile(fo.path)
			if err != nil {
				return err
			}
		} else {
			fo.file, err = CreateFile(fo.path)
			if err != nil {
				return err
			}
		}
	}
	fo.isOpen = true
	fo.curDate = time.Now()
	return nil
}

// Write
// @author Tianyi
// @description 这里不做并发控制，由 Logger 传递过来的日志数据是通过 channel 发送过来的，
//				，并不会出现多个协程往同一个文件里面写数据，文件操作模块主要集中于对日志文
//				件的分片管理，对历史日志打包
func (fo *FileOp) Write(buf []byte) error {
	if !fo.isOpen {
		_ = fo.ready()
	}

	buf = append(buf, '\n')
	_, err := fo.file.Write(buf)
	return err
}

func (fo *FileOp) Close() error {
	err := fo.file.Close()
	fo.isOpen = false
	fo.file = nil
	return err
}
