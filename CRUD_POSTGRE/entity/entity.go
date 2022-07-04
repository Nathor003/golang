package entity

import "time"

type User struct {
	Id        int
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
