package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/sirupsen/logrus"
)

type AddInventory struct {
	ProductUUID  string
	LocationUUID string
	Quantity     int
}

type addInventoryHandler struct {
	repo inventory.Repository
}

type AddInventoryHandler decorator.CommandHandler[AddInventory]

func NewAddInventoryHandler(repo inventory.Repository, logger *logrus.Entry) AddInventoryHandler {
	return decorator.ApplyCommandDecorators(addInventoryHandler{repo: repo}, logger)
}

func (h addInventoryHandler) Handle(ctx context.Context, cmd AddInventory) error {
	inv, err := inventory.NewInventory(cmd.ProductUUID, cmd.LocationUUID, cmd.Quantity)
	if err != nil {
		return err
	}

	err = h.repo.AddInventory(ctx, inv)
	if err != nil {
		return err
	}

	return nil
}
