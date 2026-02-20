package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/sirupsen/logrus"
)

type LocationByUUID struct {
	LocationUUID string
}

type locationByUUIDHandler struct {
	readModel LocationByUUIDReadModel
}

type LocationByUUIDHandler decorator.QueryHandler[LocationByUUID, *locations.Location]

func NewLocationByUUIDHandler(
	readModel LocationByUUIDReadModel,
	logger *logrus.Entry,
) LocationByUUIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[LocationByUUID, *locations.Location](
		locationByUUIDHandler{readModel: readModel},
		logger,
	)

}

type LocationByUUIDReadModel interface {
	GetLocation(ctx context.Context, locationUUID string) (*locations.Location, error)
}

func (h locationByUUIDHandler) Handle(ctx context.Context, lid LocationByUUID) (o *locations.Location, err error) {
	return h.readModel.GetLocation(ctx, lid.LocationUUID)
}
