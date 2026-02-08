package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/errors"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ConfirmOrder struct {
	orderUUID string
}

type ConfirmOrderHandler struct {
	eventBus *cqrs.EventBus
	repo     orders.Repository
}

func NewConfirmOrderHandler(eventBus *cqrs.EventBus, repo orders.Repository) ConfirmOrderHandler {
	return ConfirmOrderHandler{eventBus: eventBus, repo: repo}
}

func (h *ConfirmOrderHandler) Handle(ctx context.Context, cmd *ConfirmOrder) error {
	if err := h.repo.UpdateOrder(ctx, cmd.orderUUID, func(ctx context.Context, o *orders.Order) (*orders.Order, error) {
		err := o.UpdateStatus(orders.Confirmed)
		if err != nil {
			return nil, err
		}
		return o, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-order")
	}
	return nil
}
