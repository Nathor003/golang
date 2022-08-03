package service

import (
	"CRUD_POSTGRE_GITIGNORE/entity"
	"fmt"
)

type UserServiceIface interface {
	Register(wuser *entity.User)
}

type UserSvc struct{}

func NewUserService() UserServiceIface {
	return &UserSvc{}
}

func (u *UserSvc) Register(user *entity.User) {

	fmt.Println(user)
	// data := entity.User{
	// 	Username: "budi123",
	// 	Email:    "budi123@gmail.com",
	// 	Password: "password123",
	// 	Age:      9,
	// }

	// u, err := json.Marshal(data)

	// if err != nil {
	// 	fmt.Fprint("error")
	// }

	// fmt.Fprint(w, u)
}
