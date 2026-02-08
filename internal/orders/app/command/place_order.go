package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type PlaceOrder struct {
	LineItems   []orders.LineItem
	OrderUUID   string
	PlacedBy    string
	Weight      float32
	OrderTotal  float32
	Destination string
}

type OrderPlacedEvent struct {
	OrderUUID     string
	ProductsUUIDs []string
}

type PlaceOrderHandler struct {
	eventBus *cqrs.EventBus
	repo     orders.Repository
}

func NewPlaceOrderHandler(eventBus *cqrs.EventBus, repo orders.Repository) PlaceOrderHandler {
	return PlaceOrderHandler{eventBus: eventBus, repo: repo}
}

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd *PlaceOrder) error {
	order, err := orders.NewOrder(cmd.OrderUUID, cmd.PlacedBy, cmd.LineItems, cmd.Weight, cmd.OrderTotal, cmd.Destination)
	if err != nil {
		return err
	}

	err = h.repo.AddOrder(ctx, order)
	if err != nil {
		return err
	}

	productsUUIDs := productsUUIDsFromLineItems(cmd.LineItems)

	err = h.eventBus.Publish(ctx, &OrderPlacedEvent{OrderUUID: cmd.OrderUUID, ProductsUUIDs: productsUUIDs})
	if err != nil {
		return err
	}

	return nil
}

func productsUUIDsFromLineItems(lineItems []orders.LineItem) []string {
	productUUIDs := make([]string, len(lineItems))
	for i, lineItem := range lineItems {
		productUUIDs[i] = lineItem.ProductUUID
	}
	return productUUIDs
}
