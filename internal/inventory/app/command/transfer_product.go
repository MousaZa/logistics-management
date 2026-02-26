package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/sirupsen/logrus"
)

type TransferProduct struct {
	ProductUUID             string
	SourceLocationUUID      string
	DestinationLocationUUID string
	Quantity                int
}

type transferProductHandler struct {
	repo inventory.Repository
}

type TransferProductHandler decorator.CommandHandler[TransferProduct]

func NewTransferProductHandler(repo inventory.Repository, logger *logrus.Entry) TransferProductHandler {
	return decorator.ApplyCommandDecorators[TransferProduct](transferProductHandler{repo: repo}, logger)
}

func (h transferProductHandler) Handle(ctx context.Context, cmd TransferProduct) error {
	sourceStock, err := h.repo.GetInventory(ctx, cmd.ProductUUID, cmd.SourceLocationUUID)
	if err != nil {
		return err
	}

	destStock, err := h.repo.GetInventory(ctx, cmd.ProductUUID, cmd.DestinationLocationUUID)
	if err != nil {
		destStock, err = inventory.NewInventory(cmd.ProductUUID, cmd.DestinationLocationUUID, 0)
		if err != nil {
			return err
		}
	}

	if err := sourceStock.DecreaseQuantity(cmd.Quantity); err != nil {
		return err
	}
	if err := destStock.IncreaseQuantity(cmd.Quantity); err != nil {
		return err
	}

	return h.repo.UpdateMultipleInventories(ctx, sourceStock, destStock)
}
