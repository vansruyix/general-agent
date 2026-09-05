package user

import (
	"errors"
	"time"

	"general-agent/internal/errs"

	"gorm.io/gorm"
)

// Service 是 User 的业务逻辑层，封装数据校验、错误转换等业务规则。
// 依赖 Repository 接口而非具体实现，便于单测。
type Service struct {
	repo Repository
}

// NewService 创建 Service 实例。
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID 按 ID 查询用户，记录不存在时返回 errs.ErrNotFound。
func (s *Service) GetByID(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// List 分页查询用户列表，pageNum/pageSize 有默认值保护（默认 1/10，最大 100）。
func (s *Service) List(pageNum, pageSize int) ([]User, int64, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

// Create 创建用户，用户名重复时返回 errs.ErrUserExists。
func (s *Service) Create(req *CreateUserReq) (*User, error) {
	user := &User{
		ID:         generateID(),
		Username:   req.Username,
		Password:   req.Password,
		Name:       req.Name,
		RoleID:     req.RoleID,
		Status:     1,
		CreateTime: time.Now().Unix(),
	}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.ErrUserExists
		}
		return nil, err
	}
	return user, nil
}

// Update 按 ID 更新用户，仅更新非零值字段。记录不存在时返回 errs.ErrNotFound。
func (s *Service) Update(id string, req *UpdateUserReq) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete 按 ID 删除用户。
func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

// generateID 生成用户 ID（时间戳 + 后缀），生产环境应替换为 UUID 或雪花算法。
func generateID() string {
	return time.Now().Format("20060102150405") + "000000"
}