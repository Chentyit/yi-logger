package logger

import (
	"testing"
)

func TestExists(t *testing.T) {
	flag := isExists("../go.mod")
	if !flag {
		t.Fatal("文件不存在")
	}
	t.Log("文件存在")
}
