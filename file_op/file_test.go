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
