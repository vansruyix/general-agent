package user

import (
	"gorm.io/gorm"
)

// Repository 是 User 实体的数据访问接口，定义标准 CRUD 操作。
// 实现类 repositoryImpl 通过 gorm.DB 操作 MySQL，接口化使上层 Service 可被单测。
type Repository interface {
	// GetByID 按主键 ID 查询用户，未找到返回 gorm.ErrRecordNotFound。
	GetByID(id string) (*User, error)
	// List 分页查询用户列表，返回数据切片与总数。
	List(offset, limit int) ([]User, int64, error)
	// Create 创建用户，主键冲突返回 gorm.ErrDuplicatedKey。
	Create(user *User) error
	// Update 按 ID 更新用户字段（零值不更新）。
	Update(user *User) error
	// Delete 按 ID 软删除/硬删除用户。
	Delete(id string) error
}

// repositoryImpl 是 Repository 的 gorm 实现。
type repositoryImpl struct {
	db *gorm.DB
}

// NewRepository 创建 Repository 实现。
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// GetByID 按主键 ID 查询用户。
func (r *repositoryImpl) GetByID(id string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 分页查询用户列表。
func (r *repositoryImpl) List(offset, limit int) ([]User, int64, error) {
	var users []User
	var total int64
	if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Create 创建用户。
func (r *repositoryImpl) Create(user *User) error {
	return r.db.Create(user).Error
}

// Update 按 ID 更新用户字段，零值不更新（使用 Updates 而非 Update）。
func (r *repositoryImpl) Update(user *User) error {
	return r.db.Where("id = ?", user.ID).Updates(user).Error
}

// Delete 按 ID 删除用户。
func (r *repositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}