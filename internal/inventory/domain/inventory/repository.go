package inventory

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
)

type Repository interface {
	GetLocationProducts(ctx context.Context, locationUUID string) ([]*products.ProductStock, error)
	GetProductLocations(ctx context.Context, productUUID string) ([]*locations.ProductLocationInventory, error)
	GetInventory(ctx context.Context, productUUID, locationUUID string) (*Inventory, error)

	AddInventory(ctx context.Context, inventory *Inventory) error
	UpdateMultipleInventories(ctx context.Context, inventories ...*Inventory) error
	UpdateInventory(
		ctx context.Context,
		locationUUID string,
		productUUID string,
		updateFunc func(
			ctx context.Context,
			i *Inventory,
		) (*Inventory, error)) error
}
