package types

type APIResponse struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewAPIResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
	}
}

func NewAPIResponseWithData(code int, message string, data any) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type UserInfo struct {
	ID        int    `json:"id"`
	RoleName  string `json:"role"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CartItemResponse struct {
	ProductID   int     `json:"productId"`
	ProductName string  `json:"productName"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type CartResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}

type OrderDetailResponse struct {
	ProductID   int     `json:"productId"`
	ProductName string  `json:"productName"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    int     `json:"quantity"`
	Discount    float64 `json:"discount"`
}
