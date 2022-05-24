package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExists(t *testing.T) {
	a := assert.New(t)

	flag := isExists("../go.mod")
	a.Equal(flag, true, "文件不存在")
}
