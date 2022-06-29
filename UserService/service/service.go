package service

import (
	"UserService/entity"
	"fmt"
)

type UserServiceIface interface {
	Register(user *entity.User)
}

type UserSvc struct{}

func NewUserService() UserServiceIface {
	return &UserSvc{}
}

func (u *UserSvc) Register(user *entity.User) {

	fmt.Println(user)
}
