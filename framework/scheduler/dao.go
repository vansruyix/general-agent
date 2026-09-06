package scheduler

import (
	"general-agent/framework/mysql"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ConfJobDao struct {
	db *mysql.MysqlDB
}

func NewConfJobDao(db *mysql.MysqlDB) *ConfJobDao {

	return &ConfJobDao{db: db}
}

// Create 创建
func (d *ConfJobDao) Create(entity ConfJob) error {
	if err := d.db.Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

// Save 更新
func (d *ConfJobDao) Save(tx *gorm.DB, entity ConfJob) error {
	return tx.Save(&entity).Error
}

// UpdateEntryID 更新
func (d *ConfJobDao) UpdateEntryID(jobID string, entryID cron.EntryID) error {
	if err := d.db.Model(&ConfJob{}).Where("job_id", jobID).Update("entry_id", entryID).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (d *ConfJobDao) Delete(id string) error {
	if err := d.db.Where("job_id=?", id).Delete(&ConfJob{}).Error; err != nil {
		return err
	}
	return nil
}

// Get 根据ID查找
func (d *ConfJobDao) Get(id string) (entity ConfJob, err error) {
	if err = d.db.Where("job_id=?", id).First(&entity).Error; err != nil {
		return ConfJob{}, err
	}
	return entity, nil
}

// List 列表
func (d *ConfJobDao) List(module string) (list []ConfJob, err error) {
	// 指定部分字段
	query := d.db.Model(&ConfJob{})
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if err = query.Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}
