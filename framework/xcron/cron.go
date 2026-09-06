// Package xcron
// @author: fengyi
// @date: 2024/4/8
// @note:
package xcron

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Timer struct {
	TriggerTime string `json:"trigger_time"` // 触发时间 23:59
	Type        string `json:"type"`         // day 每天 week 每周 month 每月
	Week        []int  `json:"week"`         // type=week 不为空, 表示每周的星期几
	Month       []int  `json:"month"`        // type=month 不为空，表示每月的几号
}

func (timer *Timer) ToCron() string {
	timeParts := strings.Split(timer.TriggerTime, ":")
	minute, _ := strconv.Atoi(timeParts[1])
	hour, _ := strconv.Atoi(timeParts[0])

	var dayOfMonth, month, dayOfWeek string

	switch timer.Type {
	case "day":
		dayOfMonth = "*"
		month = "*"
		dayOfWeek = "*"
	case "week":
		// 星期日用0表示
		for i, _ := range timer.Week {
			if timer.Week[i] == 7 {
				timer.Week[i] = 0
			}
		}
		dayOfWeek = intSliceToString(timer.Week)
		dayOfMonth = "*"
		month = "*"
	case "month":
		dayOfMonth = intSliceToString(timer.Month)
		month = "*"
		dayOfWeek = "*"
	default:
		return ""
	}

	return fmt.Sprintf("0 %d %d %s %s %s", minute, hour, dayOfMonth, month, dayOfWeek)
}

func intSliceToString(slice []int) string {
	var strSlice []string
	for _, num := range slice {
		strSlice = append(strSlice, strconv.Itoa(num))
	}
	return strings.Join(strSlice, ",")
}
func (b *Timer) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}

	return json.Unmarshal(bytes, b)
}
func (j *Timer) Value() (driver.Value, error) {
	return json.Marshal(j)
}
