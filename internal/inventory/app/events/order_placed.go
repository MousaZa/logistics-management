package events

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
)

type LineItem struct {
	ProductUUID string
	Quantity    int
}

type OrderPlacedEvent struct {
	OrderUUID string
	LineItems []LineItem
}

type ExternalOrderPlacedHandler struct {
	repo inventory.Repository
}

func NewExternalOrderPlacedHandler(repo inventory.Repository) ExternalOrderPlacedHandler {
	return ExternalOrderPlacedHandler{repo: repo}
}

func (h ExternalOrderPlacedHandler) HandlerName() string {
	return "inventory_reserve_on_order_placed"
}

func (h ExternalOrderPlacedHandler) NewEvent() interface{} {
	return &OrderPlacedEvent{}
}

func (h ExternalOrderPlacedHandler) Handle(ctx context.Context, event any) error {
	e := event.(*OrderPlacedEvent)
	println("HEEEEEE", (*e).LineItems[0].ProductUUID)
	return nil
}
