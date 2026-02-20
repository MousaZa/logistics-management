package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/sirupsen/logrus"
)

type AllLocations struct{}

type allLocationsHandler struct {
	readModel AllLocationsReadModel
}

type AllLocationsHandler decorator.QueryHandler[AllLocations, []*locations.Location]

func NewAllLocationsHandler(
	readModel AllLocationsReadModel,
	logger *logrus.Entry,
) AllLocationsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllLocations, []*locations.Location](
		allLocationsHandler{readModel: readModel},
		logger,
	)
}

func (h allLocationsHandler) Handle(ctx context.Context, _ AllLocations) ([]*locations.Location, error) {
	return h.readModel.GetAllLocations(ctx)
}

type AllLocationsReadModel interface {
	GetAllLocations(ctx context.Context) ([]*locations.Location, error)
}
