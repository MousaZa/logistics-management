package domain

import "time"

type OrderStatus string

const (
	Cancelled OrderStatus = "cancelled"
	Completed OrderStatus = "completed"
	Confirmed OrderStatus = "confirmed"
	Delivered OrderStatus = "delivered"
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
)

type LineItem struct {
	LineTotal   float32
	LineWeight  float32
	ProductName string
	ProductUUID string
	Quantity    int
	UnitPrice   float32
	UnitWeight  float32
}

type Order struct {
	CompletedDate *time.Time
	DeliveredDate *time.Time
	LineItems     []LineItem
	OrderTotal    float32
	OrderUUID     string
	OrderedDate   *time.Time
	ShippedDate   *time.Time
	Status        OrderStatus
	Weight        float32
}
