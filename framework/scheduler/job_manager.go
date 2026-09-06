package scheduler

import (
	"errors"
	"general-agent/extension/errorx"
	"general-agent/extension/logz"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type JobManager struct {
	dao         *ConfJobDao
	cronManager *cron.Cron
}

func NewJobManager(dao *ConfJobDao) *JobManager {
	return &JobManager{
		dao:         dao,
		cronManager: cron.New(cron.WithSeconds()),
	}
}

// SaveJob 更新
func (s *JobManager) SaveJob(confJob ConfJob, job cron.Job) error {
	source, err := s.dao.Get(confJob.GetJobID())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	confJob.JobID = confJob.GetJobID()
	err = s.dao.db.Transaction(func(tx *gorm.DB) error {
		logz.InfoNoCtx("add job", "job_id", source.GetJobID())
		//1.将job添加到CronManager
		entryID, err := s.cronManager.AddJob(confJob.Cron, job)
		if err != nil {
			logz.InfoNoCtx("add job err", "job_id", source.GetJobID(), logz.Err(err))
			return err
		} else {
			logz.InfoNoCtx("add job ok", "job_id", source.GetJobID())
		}
		//2.保存job到数据库
		confJob.EntryID = int(entryID)
		if err := s.dao.Save(tx, confJob); err != nil {
			logz.InfoNoCtx("save job err", "job_id", confJob.GetJobID(), "entry_id", confJob.EntryID, logz.Err(err))
			return err
		}
		//3.移除已经存在的任务
		if source.EntryID != 0 {
			logz.InfoNoCtx("remove job", "job_id", source.JobID, "entry_id", source.EntryID)
			s.cronManager.Remove(cron.EntryID(source.EntryID))
		}
		return nil
	})
	return err
}

// Delete 删除
func (s *JobManager) Delete(jobID string) error {
	jobConf, err := s.dao.Get(jobID)
	if err != nil {
		return err
	}
	// 删除job
	s.cronManager.Remove(cron.EntryID(jobConf.EntryID))
	// 更新数据库
	if err := s.dao.Delete(jobID); err != nil {
		return errorx.ErrDelete.WithError(err)
	}

	return nil
}

// List 列表查询
func (s *JobManager) List(module string) ([]ConfJob, error) {
	return s.dao.List(module)
}

// Register 注册任务
var memoryJobs = map[ConfJob]cron.Job{}

func (s *JobManager) Register(jobConf ConfJob, job cron.Job) {
	memoryJobs[jobConf] = job
}
func (s *JobManager) initMemoryJob() {
	for c, job := range memoryJobs {
		logz.InfoNoCtx("初始化定时任务", logz.Any("module", c.Module), logz.Any("job", c.GetJobID()))
		_, err := s.cronManager.AddJob(c.Cron, job)
		if err != nil {
			logz.ErrorNoCtx("初始化定时任务失败", logz.Any("job", c.GetJobID()))
			panic(err)
		}
	}
}

// 从数据库加载
var dbJobs = map[string]LoadJob{}

type LoadJob func(ConfJob) (cron.Job, error)

func (s *JobManager) RegisterFromDB(jobConf string, job LoadJob) {
	dbJobs[jobConf] = job
}

func (s *JobManager) initDBJob() {
	for module, runner := range dbJobs {
		jobConfList, err := s.List(module)
		if err != nil {
			logz.ErrorNoCtx("加载模块任务", logz.Err(err), logz.Any("module", module))
			return
		}
		logz.WarnNoCtx("加载模块任务", logz.Any("module", module), logz.Any("job_size", len(jobConfList)))
		for _, c := range jobConfList {
			job, err := runner(c)
			if err != nil {
				logz.ErrorNoCtx("初始化作业失败", logz.Any("job", c.GetJobID()))
				panic(err)
			}
			logz.InfoNoCtx("初始化定时任务", logz.Any("module", c.Module), logz.Any("job", c.GetJobID()))
			entryID, err := s.cronManager.AddJob(c.Cron, job)
			if err != nil {
				logz.ErrorNoCtx("初始化定时任务失败", logz.Any("job", c.GetJobID()))
				panic(err)
			}
			if err := s.dao.UpdateEntryID(c.JobID, entryID); err != nil {
				logz.ErrorNoCtx("初始化定时任务失败", logz.Any("job", c.GetJobID()))
				panic(err)
			}
		}
	}
}

// Start 初始化定时任务
func (s *JobManager) Start() {
	s.initMemoryJob()
	s.initDBJob()
	s.cronManager.Start()
	logz.InfoNoCtx("定时任务启动")
}

// Stop 停止定时任务
func (s *JobManager) Stop() {
	s.cronManager.Stop()
	logz.InfoNoCtx("定时任务停止")
}
