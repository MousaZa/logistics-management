package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/sirupsen/logrus"
)

type ProductLocations struct {
	ProductUUID string
}

type productLocationsHandler struct {
	productLocations ProductLocationsReadModel
	logger           *logrus.Entry
}

type ProductLocationsHandler decorator.QueryHandler[ProductLocations, []*locations.Location]

func NewProductLocationsHandler(
	productLocations ProductLocationsReadModel,
	logger *logrus.Entry,
) ProductLocationsHandler {
	if productLocations == nil {
		panic("nil productLocations")
	}

	return decorator.ApplyQueryDecorators[ProductLocations, []*locations.Location](
		productLocationsHandler{productLocations: productLocations, logger: logger},
		logger,
	)
}

func (h productLocationsHandler) Handle(ctx context.Context, q ProductLocations) ([]*locations.Location, error) {
	return h.productLocations.GetProductLocations(ctx, q.ProductUUID)
}

type ProductLocationsReadModel interface {
	GetProductLocations(ctx context.Context, productUUID string) ([]*locations.Location, error)
}
