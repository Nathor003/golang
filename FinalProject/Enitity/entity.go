package enitity

import "time"

type User struct {
	Id        string
	Username  string
	Email     string
	Password  string
	Age       int
	CreateAt  time.Time
	UpdatedAt time.Time
}

type Photo struct {
	Id        string
	Title     string
	Caption   string
	PhotoUrl  string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	Id        string
	UserId    string
	PhotoId   string
	Message   string
	CreateAt  time.Time
	UpdatedAt time.Time
}

type SocialMedia struct {
	Id             string
	Name           string
	SocialMediaUrl string
	UserId         string
}
