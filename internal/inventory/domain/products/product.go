package products

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	AvailableQuantity int
	CreatedAt         time.Time
	ProductUUID       string
	LocationUUID      string
	Name              string
	Price             float32
	UpdatedAt         time.Time
	Weight            float32
}

func NewProduct(name string, price float32, weight float32, locationUUID string, availableQuantity int) (*Product, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	if price == 0 {
		return nil, errors.New("empty price")
	}
	if weight == 0 {
		return nil, errors.New("empty weight")
	}
	if locationUUID == "" {
		return nil, errors.New("empty locationUUID")
	}

	productUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Product{
		ProductUUID:       productUUID.String(),
		Name:              name,
		Price:             price,
		Weight:            weight,
		LocationUUID:      locationUUID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AvailableQuantity: availableQuantity,
	}, nil
}
