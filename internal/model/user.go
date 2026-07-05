package model

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	IsActive  bool
	RoleID    int
}

type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=5,max=16"`
}

type LogUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=16"`
}
