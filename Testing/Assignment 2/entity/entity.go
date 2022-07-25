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

type Item struct {
	LineItemID  int    `json:"lineItemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Order struct {
	CustomerName string `json:"customerName"`
	OrderAt      string `json:"orderAt"`
	Barang       []Item `json:"items"`
}
