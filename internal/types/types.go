package types

import "time"

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID           int    `json:"id"`
	RoleID       int    `json:"roleId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type Brand struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Website *string `json:"website"`
}

type Category struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Product struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"categoryId"`
	BrandID     int     `json:"brandId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"imageUrl"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Status    string    `json:"status"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderDetail struct {
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Discount  float64 `json:"discount"`
}

type Cart struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
}

type CartItem struct {
	CartID    int `json:"cartId"`
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}
