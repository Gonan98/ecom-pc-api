package repository

type CartItemWithProduct struct {
	CartID      int
	ProductID   int
	ProductName string
	Quantity    int
	Discount    float64
	UnitPrice   float64
}
