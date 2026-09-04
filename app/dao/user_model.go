package dao

import (
	"general-agent/app/entity"

	"gorm.io/gorm"
)

var Table = "user"

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (u *UserModel) GetById(id int) (user *entity.User, err error) {
	if u.db.Model(&entity.User{}).Where("id = ?", id).Find(user).Error != nil {
		return nil, err
	}
	return user, nil
}
