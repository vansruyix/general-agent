// Package scheduler
// @author: fengyi
// @date: 2024/6/7
// @note:
package scheduler

import "fmt"

type ConfJob struct {
	Module  string `json:"module"`                   // 模块
	JobID   string `json:"job_id" gorm:"primarykey"` // 作业ID
	Cron    string `json:"cron"`                     // cron表达式
	EntryID int    `json:"entry_id"`                 // cron实例ID
	Param   string `json:"param"`                    // 任务参数
}

func NewConfJob(module string, jobID string, cron string, param string) ConfJob {
	return ConfJob{Module: module, JobID: jobID, Cron: cron, Param: param}
}

func (c ConfJob) GetJobID() string {
	return fmt.Sprintf("%s_%s", c.Module, c.JobID)
}

func GetJobID(module, jobID string) string {
	return fmt.Sprintf("%s_%s", module, jobID)
}
func (ConfJob) TableName() string {
	return "tgp_job"
}
