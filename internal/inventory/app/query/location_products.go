package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/sirupsen/logrus"
)

type LocationProducts struct {
	LocationUUID string
}

type locationProductsHandler struct {
	locationProducts LocationProductsReadModel
	logger           *logrus.Entry
}

type LocationProductsHandler decorator.QueryHandler[LocationProducts, []*products.ProductStock]

func NewLocationProductsHandler(
	locationProducts LocationProductsReadModel,
	logger *logrus.Entry,
) LocationProductsHandler {
	if locationProducts == nil {
		panic("nil locationProducts")
	}

	return decorator.ApplyQueryDecorators[LocationProducts, []*products.ProductStock](
		locationProductsHandler{locationProducts: locationProducts, logger: logger},
		logger,
	)
}

func (h locationProductsHandler) Handle(ctx context.Context, q LocationProducts) ([]*products.ProductStock, error) {
	return h.locationProducts.GetLocationProducts(ctx, q.LocationUUID)
}

type LocationProductsReadModel interface {
	GetLocationProducts(ctx context.Context, LocationUUID string) ([]*products.ProductStock, error)
}
