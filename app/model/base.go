package model

import (
	"general-agent/extension/pagination"
	"general-agent/extension/xtime"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint64         `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt xtime.JsonTime `json:"created_at,omitempty" gorm:"->;<-:create" swaggerignore:"true"` //创建时间
	UpdatedAt xtime.JsonTime `json:"updated_at,omitempty" gorm:"<-" swaggerignore:"true"`           //更新时间
}

type BaseModel struct {
	ID        string         `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt xtime.JsonTime `json:"created_at,omitempty" gorm:"->;<-:create" swaggerignore:"true"` //创建时间
	UpdatedAt xtime.JsonTime `json:"updated_at,omitempty" gorm:"<-" swaggerignore:"true"`           //更新时间
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = xtime.Now()
	b.UpdatedAt = xtime.Now()
	return nil
}
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = xtime.Now()
	return nil
}

type IdReq struct {
	ID string `json:"id" form:"id" uri:"id" binding:"required"` // id
}

type IdNameReq struct {
	Id   string `json:"id" form:"id" `
	Name string `json:"name" form:"name" binding:"required_if=Id 0"`
}

type IdsReq struct {
	Ids []string `json:"ids" form:"ids"` // ids
}

type FindReq struct {
	pagination.LimitOffsetPagination
}
