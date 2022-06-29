package main

import (
	"UserService/entity"
	"UserService/service"
)

func main() {

	userSvc := service.NewUserService()

	userSvc.Register(&entity.User{
		Username: "budi123",
		Email:    "budi123@gmail.com",
		Password: "password123",
		Age:      9,
	})

}
