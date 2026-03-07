package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/sirupsen/logrus"
)

type ReportDamaged struct {
	ProductUUID  string
	LocationUUID string
	Reason       string
	Quantity     int
}

type reportDamagedHandler struct {
	repo     inventory.Repository
	eventBus *cqrs.EventBus
}

type ReportDamagedHandler decorator.CommandHandler[ReportDamaged]

func NewReportDamagedHandler(repo inventory.Repository, eventBus *cqrs.EventBus, logger *logrus.Entry) ReportDamagedHandler {
	return decorator.ApplyCommandDecorators(reportDamagedHandler{repo: repo, eventBus: eventBus}, logger)
}

func (h reportDamagedHandler) Handle(ctx context.Context, cmd ReportDamaged) error {
	return h.repo.UpdateInventory(ctx, cmd.LocationUUID, cmd.ProductUUID, func(ctx context.Context, i *inventory.Inventory) (*inventory.Inventory, error) {
		err := i.MarkQuantityAsDamaged(cmd.Quantity)
		if err != nil {
			return nil, err
		}
		return i, nil
	})
}
