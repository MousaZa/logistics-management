package inventory

import "errors"

type Inventory struct {
	ProductUUID  string
	LocationUUID string
	Quantity     int
}

func NewInventory(productUUID string, locationUUID string, quantity int) (*Inventory, error) {
	if productUUID == "" {
		return nil, errors.New("product uuid is empty")
	}
	if locationUUID == "" {
		return nil, errors.New("location uuid is empty")
	}
	if quantity < 0 {
		return nil, errors.New("quantity is negative")
	}

	return &Inventory{
		ProductUUID:  productUUID,
		LocationUUID: locationUUID,
		Quantity:     quantity,
	}, nil
}
