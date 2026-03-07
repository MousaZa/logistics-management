package inventory

import "errors"

type Status string

//const (
//	StatusAvailable Status = "available"
//	StatusReserved  Status = "reserved"
//	StatusDamaged   Status = "damaged"
//)

type Inventory struct {
	ProductUUID       string
	LocationUUID      string
	AvailableQuantity int
	ReservedQuantity  int
	DamagedQuantity   int
}

func NewInventory(productUUID string, locationUUID string, quantity int) (*Inventory, error) {
	if productUUID == "" {
		return nil, errors.New("products uuid is empty")
	}
	if locationUUID == "" {
		return nil, errors.New("locations uuid is empty")
	}
	if quantity < 0 {
		return nil, errors.New("quantity is negative")
	}

	return &Inventory{
		ProductUUID:       productUUID,
		LocationUUID:      locationUUID,
		AvailableQuantity: quantity,
		ReservedQuantity:  0,
		DamagedQuantity:   0,
	}, nil
}

func (i *Inventory) DecreaseQuantity(amount int) error {
	if amount <= 0 {
		return errors.New("must decrease by a positive amount")
	}
	if i.AvailableQuantity < amount {
		return errors.New("insufficient stock for transfer")
	}
	i.AvailableQuantity -= amount
	return nil
}

func (i *Inventory) IncreaseQuantity(amount int) error {
	if amount <= 0 {
		return errors.New("must increase by a positive amount")
	}
	i.AvailableQuantity += amount
	return nil
}
