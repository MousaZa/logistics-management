package products

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	LocationUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Product '%s' not found", e.LocationUUID)
}

type Repository interface {
	AddProduct(ctx context.Context, product *Product) error
	GetAllProducts(ctx context.Context) ([]*Product, error)
	GetProduct(ctx context.Context, productUUID string) (*Product, error)
	UpdateProduct(
		ctx context.Context,
		productUUID string,
		updateFunc func(
			ctx context.Context,
			p *Product,
		) (*Product, error)) error
}
