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
)

// IsExists
// @description 判断路径是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

// IsPermission
// description 判断文件是否有权限操作
func IsPermission(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrPermission)
}

// Mkdir
// @description 创建一个目录
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile
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

// Remove
// @description 删除文件
func Remove(path string) error {
	return os.RemoveAll(path)
}

// ChangeFileName
// description 修改文件名
func ChangeFileName(path string, dstName string) (string, error) {
	preDir := filepath.Dir(path)
	newFilePath := filepath.Join(preDir, dstName)
	err := os.Rename(path, newFilePath)
	return newFilePath, err
}

// MustOpenFile
// @description 直接打开文件，使用该方法的前提是确定文件一定存在
func MustOpenFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0666)
	return file, err
}

// Compress
// @param filePath 需要压缩文件或者目录的路径
// @param dest 压缩目标文件
// @description 压缩文件
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
// @param srcPath 压缩包路径
// @param dstPath 解压路径
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
