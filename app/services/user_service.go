package services

import "general-agent/app/dao"

type UserService struct {
	dao *dao.UserModel
}

func NewUserService(dao *dao.UserModel) *UserService {
	return &UserService{
		dao: dao,
	}
}
