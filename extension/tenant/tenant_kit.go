// Package tenant
// @author: liulun
// @date: 2024/11/13
// @note: 获取租户id
package tenant

import (
	"general-agent/extension/xtime"

	"github.com/gin-gonic/gin"
)

func GetTenantId(ctx *gin.Context) string {
	tenantId := ctx.GetHeader("tenant")
	return tenantId
}
func GetIntervalCount(startTimeStr string, endTimeStr string) int {
	if startTimeStr == "" || endTimeStr == "" {
		return 30
	}
	startTime, err := xtime.StringToTime(startTimeStr)
	if err != nil {
		return 30
	}
	endTime, err := xtime.StringToTime(endTimeStr)
	if err != nil {
		return 30
	}
	minute := endTime.Sub(startTime).Minutes()
	if minute < 1440 {
		return int(minute / 2)
	} else if minute == 1440 {
		return 144
	} else {
		return int(minute / 120)
	}
}
