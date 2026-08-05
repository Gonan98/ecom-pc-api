package repository

type DetailWithProduct struct {
	OrderID     int
	ProductID   int
	ProductName string
	UnitPrice   float64
	Quantity    int
	Discount    float64
}
