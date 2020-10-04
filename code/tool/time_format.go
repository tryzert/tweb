package tool

import (
	"strings"
	"time"
)

//将 time.Time 格式化为字符串
func TimeFormat(t time.Time) (string, string) {
	fres := strings.Split(t.Format("2006-01-02 15:04:05"), " ")
	date, tm := fres[0], fres[1]
	return date, tm
}
