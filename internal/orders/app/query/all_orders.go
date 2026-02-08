package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
)

type AllOrders struct{}

type AllOrdersHandler struct {
	readModel AllOrdersReadModel
}

func NewAllOrdersHandler(
	readModel AllOrdersReadModel,
) AllOrdersHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return AllOrdersHandler{readModel: readModel}

}

type AllOrdersReadModel interface {
	GetAllOrders(ctx context.Context) ([]*orders.Order, error)
}

func (h AllOrdersHandler) Handle(ctx context.Context, _ AllOrders) (o []*orders.Order, err error) {
	return h.readModel.GetAllOrders(ctx)
}
