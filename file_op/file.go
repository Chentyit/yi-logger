package file_op

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileOp struct {
	file         *os.File
	isOpen       bool // 用于判断是否可以进行操作
	needCompress bool // 是否需要压缩
	maxSize      int  // 以 MB 为单位
	curDate      time.Time
	path         string
}

func CreateFileOp(path string, maxSize int, needCompress bool) *FileOp {
	return &FileOp{
		path:         path,
		needCompress: needCompress,
		isOpen:       false,
		maxSize:      maxSize,
	}
}

// ready
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
// @param buf 需要写入的字节
// @description 这里不做并发控制，由 Logger 传递过来的日志数据是通过 channel 发送过来的，
//				，并不会出现多个协程往同一个文件里面写数据，文件操作模块主要集中于对日志文
//				件的分片管理，对历史日志打包
func (fo *FileOp) Write(buf []byte) error {
	if !fo.isOpen {
		_ = fo.ready()
	}

	var wg sync.WaitGroup

	// 判断当前文件是否超出 maxSize
	// 如果超出了最大限制，则需要进行以下操作:
	// - 断开 fo.file 指针
	// - 创建新文件，并将 fo.file 指向新的文件
	// - 将原来的文件压缩打包
	if fo.overMaxSize() {
		_ = fo.Close()

		now := time.Now()

		date, month, day := now.Date()
		timestamp := now.Unix()

		// 获取原文件路径
		filePreDir := filepath.Dir(fo.path)
		// 获取原文件名称和扩展名
		fileNameAndExt := strings.Split(filepath.Base(fo.path), ".")
		fileName := fileNameAndExt[0]
		fileExt := fileNameAndExt[1]
		// 拼接新文件名（fileName-year-month-day-timestamp.fileExt)
		changeFileName := fmt.Sprintf("%s-%v-%v-%v-%v.%s", fileName, date, int(month), day, timestamp, fileExt)
		// 先改名再压缩是为了防止数据写入时因为压缩速度太慢而造成阻塞
		changeFilePath, err := ChangeFileName(fo.path, changeFileName)
		if err != nil {
			return err
		}
		// 重新初始化 fo.file 继续写
		_ = fo.ready()

		// 判断用户是否设置压缩
		if fo.needCompress {
			// 使用同步锁保证压缩过程不会中断
			wg.Add(1)
			go func() {
				wg.Done()
				pkgName := fmt.Sprintf("%s-%v-%v-%v-%v.%s", fileName, date, int(month), day, timestamp, "zip")
				pkgPath := filepath.Join(filePreDir, pkgName)
				_ = Compress(pkgPath, changeFilePath)
				// 删除原文件
				_ = Remove(changeFilePath)
			}()
		}
	}

	buf = append(buf, '\n')
	_, err := fo.file.Write(buf)

	// 等待压缩完成
	wg.Wait()
	return err
}

func (fo *FileOp) Close() error {
	err := fo.file.Close()
	fo.isOpen = false
	fo.file = nil
	return err
}

// overMaxSize
// @description 判断该 FileOp 指向的文件是否超过最大值
func (fo *FileOp) overMaxSize() bool {
	info, _ := fo.file.Stat()

	if info.Size() > int64(fo.maxSize*1024*1024) {
		return true
	}

	return false
}
