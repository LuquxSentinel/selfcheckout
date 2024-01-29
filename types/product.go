package types

import "time"

type Product struct {
	ID          string
	BrandName   string
	ProductName string
	ExpiryDate  time.Time
	Price       float64
}

type BusketProduct struct {
	ID        string
	UnitPrice float64
	UnitSize  int
}
