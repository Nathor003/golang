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

	// if user, err := userSvc.Register(&entity.User{
	// 	Username: "budi123",
	// 	Email:    "budi123@gmail.com",
	// 	Password: "password123",
	// 	Age:      9,
	// }); err != nil {
	// 	fmt.Printf("Error when register user: %+v", err)
	// 	return
	// } else {
	// 	fmt.Printf("Success register user: %+v", user)
	// }
}
