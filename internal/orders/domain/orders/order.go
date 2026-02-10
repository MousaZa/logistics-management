package orders

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

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
	CompletedDate time.Time
	DeliveredDate time.Time
	LineItems     []LineItem
	OrderTotal    float32
	OrderUUID     string
	PlacedBy      string
	OrderedDate   time.Time
	ShippedDate   time.Time
	Status        OrderStatus
	Weight        float32
	Destination   string
}

func NewOrder(placedBy string, lineItems []LineItem, orderTotal float32, weight float32, destination string) (*Order, error) {
	if placedBy == "" {
		return nil, errors.New("empty placedBy")
	}
	if len(lineItems) == 0 {
		return nil, errors.New("empty lineItems")
	}
	if orderTotal == 0 {
		return nil, errors.New("empty orderTotal")
	}
	if weight == 0 {
		return nil, errors.New("empty orderWeight")
	}
	if destination == "" {
		return nil, errors.New("empty destination")
	}
	orderUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Order{
		OrderUUID:   orderUUID.String(),
		PlacedBy:    placedBy,
		LineItems:   lineItems,
		OrderedDate: time.Now(),
		Weight:      weight,
		OrderTotal:  orderTotal,
		Status:      Pending,
		Destination: destination,
	}, nil
}

func (o *Order) UpdateStatus(s OrderStatus) error {
	if o.Status == s {
		return errors.New("order already with the same status")
	}
	o.Status = s
	return nil
}
