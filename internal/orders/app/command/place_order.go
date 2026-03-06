package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/sirupsen/logrus"
)

type PlaceOrder struct {
	LineItems   []orders.LineItem
	PlacedBy    string
	Weight      float32
	OrderTotal  float32
	Destination string
}

type PlaceOrderHandler decorator.CommandHandler[PlaceOrder]

type placeOrderHandler struct {
	eventBus *cqrs.EventBus
	repo     orders.Repository
}

func NewPlaceOrderHandler(eventBus *cqrs.EventBus, repo orders.Repository, logger *logrus.Entry) PlaceOrderHandler {
	return decorator.ApplyCommandDecorators[PlaceOrder](placeOrderHandler{eventBus: eventBus, repo: repo}, logger)
}

func (h placeOrderHandler) Handle(ctx context.Context, cmd PlaceOrder) error {
	order, err := orders.NewOrder(cmd.PlacedBy, cmd.LineItems, cmd.Weight, cmd.OrderTotal, cmd.Destination)
	if err != nil {
		return err
	}

	err = h.repo.AddOrder(ctx, order)
	if err != nil {
		return err
	}

	lineItems := orders.EventLineItemsFromLineItems(cmd.LineItems)

	err = h.eventBus.Publish(ctx, &orders.OrderPlacedEvent{OrderUUID: order.OrderUUID, LineItems: lineItems})
	if err != nil {
		return err
	}

	return nil
}
