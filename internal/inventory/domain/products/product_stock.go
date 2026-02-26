package products

import (
	"time"
)

type ProductStock struct {
	CreatedAt   time.Time
	ProductUUID string
	Name        string
	Price       float32
	UpdatedAt   time.Time
	Weight      float32
	Quantity    int
}
