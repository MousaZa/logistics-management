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
	LineTotal   float32 `json:"line_total"`
	LineWeight  float32 `json:"line_weight"`
	ProductName string  `json:"product_name"`
	ProductUUID string  `json:"product_uuid"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float32 `json:"unit_price"`
	UnitWeight  float32 `json:"unit_weight"`
}

type Order struct {
	CompletedDate time.Time
	DeliveredDate time.Time
	LineItems     []LineItem `db:"line_items" json:"line_items"`
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
