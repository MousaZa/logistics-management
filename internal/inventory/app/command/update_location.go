package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/sirupsen/logrus"
)

type UpdateLocation struct {
	LocationUUID string
	Name         *string
	Address      *string
	City         *string
}

type updateLocationHandler struct {
	repo locations.Repository
}

type UpdateLocationHandler decorator.CommandHandler[UpdateLocation]

func NewUpdateLocationHandler(repo locations.Repository, logger *logrus.Entry) UpdateLocationHandler {
	return decorator.ApplyCommandDecorators[UpdateLocation](updateLocationHandler{repo: repo}, logger)
}

func (h updateLocationHandler) Handle(ctx context.Context, cmd UpdateLocation) error {
	return h.repo.UpdateLocation(ctx, cmd.LocationUUID, func(ctx context.Context, l *locations.Location) (*locations.Location, error) {
		if cmd.Name != nil && *cmd.Name != "" {
			l.Name = *cmd.Name
		}
		if cmd.Address != nil && *cmd.Address != "" {
			l.Address = *cmd.Address
		}
		if cmd.City != nil && *cmd.City != "" {
			l.City = *cmd.City
		}
		return l, nil
	})
}
