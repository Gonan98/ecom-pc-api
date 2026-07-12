package types

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

type CartItemRequest struct {
	ProductID int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"gt=0"`
}

type CreateBrandRequest struct {
	Name    string `json:"name" validate:"required"`
	Website string `json:"website"`
}

type UpdateBrandRequest struct {
	CreateBrandRequest
}
