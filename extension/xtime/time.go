package xtime

import (
	"time"
)

const (
	TimeLayout        = "2006-01-02 15:04:05"
	TimeLayoutCompact = "20060102150405"
)

//var timeLocationGTM8 *time.Location

// Now get the current time in East 8.

//func init() {
//	// load time location of East 8 from tzData string.
//	var err error
//	timeLocationGTM8, err = time.LoadLocation("Asia/Shanghai")
//	if err != nil {
//		panic("load ShanghaiTimeLocation failed from TZData, error: " + err.Error())
//	}
//}

// Int10ToString 十位时间戳转字符串
func Int10ToString(timeInt int) string {
	timestamp := int64(timeInt)
	// 将十位时间戳转换为 time.Time 对象
	t := time.Unix(timestamp, 0)
	// 将时间转换为字符串
	return t.Format(TimeLayout)
}
func StringToTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeLayout, timeStr)
}
func TimeToString(timed time.Time) string {
	return timed.Format(TimeLayout)
}
