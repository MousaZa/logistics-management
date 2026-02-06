package domain

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	OrderUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order '%s' not found", e.OrderUUID)
}

type Repository interface {
	AddOrder(ctx context.Context, o *Order) error
	GetOrder(ctx context.Context, orderUUID string) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(
		ctx context.Context,
		o *Order,
		updateFunc func(ctx context.Context, tr *Order) (*Order, error),
	) error
}
