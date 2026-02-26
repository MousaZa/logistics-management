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

func (i *Inventory) DecreaseQuantity(amount int) error {
	if amount <= 0 {
		return errors.New("must decrease by a positive amount")
	}
	if i.Quantity < amount {
		return errors.New("insufficient stock for transfer")
	}
	i.Quantity -= amount
	return nil
}

func (i *Inventory) IncreaseQuantity(amount int) error {
	if amount <= 0 {
		return errors.New("must increase by a positive amount")
	}
	i.Quantity += amount
	return nil
}
