package file_op

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExists(t *testing.T) {
	a := assert.New(t)

	flag := IsExists("../logger")
	a.Equal(flag, true, "路径不存在")
}

func TestIsPermission(t *testing.T) {
	a := assert.New(t)

	flag := IsPermission("../go.mod")
	a.Equal(flag, true, "文件没有权限")
}

func TestMkDir(t *testing.T) {
	a := assert.New(t)

	err := Mkdir("../tt")
	a.Equal(err, nil, "创建目录失败")
}

func TestCreateFile(t *testing.T) {
	a := assert.New(t)

	_, err := CreateFile("../tt.txt")
	a.Equal(err, nil, err)
}

func TestCreateFileOp(t *testing.T) {
	a := assert.New(t)

	fileOp := CreateFileOp("../test.txt", 10, true)
	a.EqualValues(fileOp.path, "../test.txt", "FileOp 对象创建失败")
}

func TestFileOpWrite(t *testing.T) {
	a := assert.New(t)

	fileOp := CreateFileOp("../test.txt", 10, true)
	err := fileOp.ready()
	a.Equal(err, nil, err)
	_ = fileOp.Write([]byte("hello world"))
	_ = fileOp.Close()
	_ = fileOp.Write([]byte("hello world2"))
}

func TestCompress(t *testing.T) {
	_ = Compress("../test.zip", "../logger", "../test_big_file.txt")
}

func TestDecompress(t *testing.T) {
	_ = Decompress("../test.zip", "../zip")
}

func TestChangeFileName(t *testing.T) {
	_, _ = ChangeFileName("../test_big_file.txt", "test-2022-06-12")
}

func TestFileWriteOverSize(t *testing.T) {
	a := assert.New(t)

	fileOp := CreateFileOp("../test_big_file.txt", 10, true)

	err := fileOp.Write([]byte("hello world"))
	a.Equal(err, nil, err)
}

func TestBuildBigFile(t *testing.T) {
	fileOp := CreateFileOp("../test_big_file.txt", 20, false)
	for {
		_ = fileOp.Write([]byte("hello world"))
	}
}
