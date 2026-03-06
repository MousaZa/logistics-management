package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/common/errors"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/sirupsen/logrus"
)

type CancelOrder struct {
	orderUUID string
}

type CancelOrderHandler decorator.CommandHandler[CancelOrder]

type cancelOrderHandler struct {
	eventBus *cqrs.EventBus
	repo     orders.Repository
}

func NewCancelOrderHandler(eventBus *cqrs.EventBus, repo orders.Repository, logger *logrus.Entry) CancelOrderHandler {
	return decorator.ApplyCommandDecorators[CancelOrder](cancelOrderHandler{eventBus: eventBus, repo: repo}, logger)
}

func (h cancelOrderHandler) Handle(ctx context.Context, cmd CancelOrder) error {
	if err := h.repo.UpdateOrder(ctx, cmd.orderUUID, func(ctx context.Context, o *orders.Order) (*orders.Order, error) {
		err := o.UpdateStatus(orders.Cancelled)
		if err != nil {
			return nil, err
		}
		return o, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-order")
	}
	return nil
}
