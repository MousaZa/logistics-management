package locations

import "time"

type ProductLocationInventory struct {
	Address           string
	City              string
	CreatedAt         time.Time
	LocationUUID      string
	Name              string
	UpdatedAt         time.Time
	AvailableQuantity int
	DamagedQuantity   int
	ReservedQuantity  int
}
