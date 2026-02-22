package inventory

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
)

type Repository interface {
	GetLocationProducts(ctx context.Context, locationUUID string) ([]*products.Product, error)
	GetProductLocations(ctx context.Context, productUUID string) ([]*locations.Location, error)
}
