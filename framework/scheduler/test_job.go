// Package scheduler
// @author: fengyi
// @date: 2024/8/14
// @note:
package scheduler

import (
	"general-agent/extension/logz"
)

type TestJob struct {
	Message string
}

func (job *TestJob) Run() {
	logz.InfoNoCtx("测试任务：" + job.Message)
}

func registerJob(jobManager *JobManager) {
	// 自定义任务实例
	//jobManager.RegisterFromDB("test", func(c ConfJob) (cron.Job, error) {
	//	runner := &TestJob{"db job"}
	//	return runner, nil
	//})

	//conf := ConfJob{Module: "test",
	//	JobID: "test memory jobs",
	//	Cron:  "0/30 * * * * ?"}
	//jobManager.Register(conf, &TestJob{"memory job"})
}
