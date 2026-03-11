package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/sirupsen/logrus"
)

type AddLocation struct {
	Name    string
	Address string
	City    string
	Lat     float64
	Lon     float64
}

type AddLocationHandler decorator.CommandHandler[AddLocation]

type addLocationHandler struct {
	repo locations.Repository
}

func NewAddLocationHandler(repo locations.Repository, logger *logrus.Entry) AddLocationHandler {
	return decorator.ApplyCommandDecorators[AddLocation](addLocationHandler{repo: repo}, logger)
}

func (h addLocationHandler) Handle(ctx context.Context, cmd AddLocation) error {
	location, err := locations.NewLocation(cmd.Name, cmd.Address, cmd.City, cmd.Lat, cmd.Lon)
	if err != nil {
		return err
	}

	err = h.repo.AddLocation(ctx, location)
	if err != nil {
		return err
	}

	return nil
}
