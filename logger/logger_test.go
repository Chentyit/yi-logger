package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLoggerConfig(t *testing.T) {
	ass := assert.New(t)

	BuildConfig().SetDateFormat(LogDateFormat.ShortLine).SetTimeFormat(LogTimeFormat.Default).Build()

	df := config.dateFormat
	fmt.Println(df)
	fmt.Println(config)
	ass.EqualValues(df, "yyyy-MM-dd", "日期格式错误")
}
