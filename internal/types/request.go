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
	Name    string  `json:"name" validate:"required"`
	Website *string `json:"website" validate:"omitempty,url"`
}

type UpdateBrandRequest struct {
	Name    string  `json:"name" validate:"required"`
	Website *string `json:"website" validate:"omitempty,url"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type CreateProductRequest struct {
	CategoryID  int     `json:"categoryId" validate:"gt=0"`
	BrandID     int     `json:"brandId" validate:"gt=0"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"imageUrl" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
	CategoryID  int     `json:"categoryId" validate:"gt=0"`
	BrandID     int     `json:"brandId" validate:"gt=0"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"imageUrl" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}
