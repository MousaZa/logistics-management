package events

import (
	"context"
	"errors"

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

	count := 0
	for _, li := range e.LineItems {
		inv, err := h.repo.GetInventoriesByProduct(ctx, li.ProductUUID)
		if err != nil {
			return err
		}

		for _, i := range inv { // TODO improve the mechanism
			err := i.ReserveQuantity(li.Quantity)
			if err == nil {
				err = h.repo.UpdateInventory(ctx, i.LocationUUID, i.ProductUUID, func(ctx context.Context, i *inventory.Inventory) (*inventory.Inventory, error) {
					err = i.ReserveQuantity(li.Quantity)
					if err != nil {
						return nil, err
					}
					return i, nil
				})
				if err != nil {
					return err
				}
				count++
				break
			}
		}
	}
	if count == len(e.LineItems) { // TODO inform the orders service
		return nil
	} else {
		return errors.New("failed to reserve inventory")
	}
}
