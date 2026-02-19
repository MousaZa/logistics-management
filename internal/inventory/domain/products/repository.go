package products

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	LocationUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("location '%s' not found", e.LocationUUID)
}

type LocationRepository interface {
	AddProduct(ctx context.Context, location *Product) error
	GetAllProducts(ctx context.Context) ([]*Product, error)
	GetLocationProducts(ctx context.Context, locationUUID string) ([]*Product, error)
	GetProduct(ctx context.Context, locationUUID string) (*Product, error)
	UpdateProduct(
		ctx context.Context,
		productUUID string,
		updateFunc func(
		ctx context.Context,
		o *Product,
	) (*Product, error)) error
}
