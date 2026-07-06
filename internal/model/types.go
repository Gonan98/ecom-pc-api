package model

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
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Website string `json:"website"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"categoryId"`
	BrandID     int     `json:"brandId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"imageUrl"`
	Price       float64 `json:"price"`
	Stock       float64 `json:"stock"`
}
