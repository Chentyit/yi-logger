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

// Compress
// @author Tianyi
// @param filePath 需要压缩文件或者目录的路径
// @param dest 压缩目标文件
// @description 将文件压缩，触发压缩有以下几种情况：
// 				1. 手动触发，会将当前 fo.file 打包
// 				2. 没跨天，但是超过了 maxSize 会打包
//				3. 跨天打包
func Compress(pkgPath string, paths ...string) error {
	// 获取上级目录路径
	preDir := filepath.Dir(pkgPath)
	if err := os.MkdirAll(preDir, os.ModePerm); err != nil {
		return err
	}

	// 创建压缩文件
	archive, err := os.Create(pkgPath)
	if err != nil {
		return err
	}
	defer func(archive *os.File) {
		_ = archive.Close()
	}(archive)

	// 创建 zip writer
	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	// 遍历需要打包的路径
	for _, srcPath := range paths {
		// 删除最后一个 '/'
		srcPath = strings.TrimSuffix(srcPath, string(os.PathSeparator))

		// 开始检查文件树
		err = filepath.Walk(
			srcPath,
			func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				header, err := zip.FileInfoHeader(info)
				if err != nil {
					return err
				}

				// 设置压缩方式
				header.Method = zip.Deflate

				// 将文件的相对路径设置为头名称
				header.Name, err = filepath.Rel(filepath.Dir(srcPath), path)
				if err != nil {
					return err
				}
				if info.IsDir() {
					header.Name += string(os.PathSeparator)
				}

				// 创建文件头写入器并保存文件内容
				headerWriter, err := zipWriter.CreateHeader(header)
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer func(f *os.File) {
					_ = f.Close()
				}(f)
				_, err = io.Copy(headerWriter, f)
				return err
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Decompress
// @author Tianyi
// @param srcPath 压缩包路径
// @param dstPath 解压路径
// @description
func Decompress(srcPath, dstPath string) error {
	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	for _, file := range reader.File {
		if err := decompress(file, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func decompress(file *zip.File, dstPath string) error {
	// create the directory of file
	filePath := path.Join(dstPath, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	r, err := file.Open()
	if err != nil {
		return err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(w)

	// save the decompressed file content
	_, err = io.Copy(w, r)
	return err
}
