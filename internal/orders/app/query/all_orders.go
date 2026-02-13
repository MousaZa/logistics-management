package query

import (
	"context"
	"time"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/sirupsen/logrus"
)

type AllOrders struct {
	DateFrom time.Time
	DateTo   time.Time
}

type allOrdersHandler struct {
	readModel AllOrdersReadModel
}

type AllOrdersHandler decorator.QueryHandler[AllOrders, []*orders.Order]

func NewAllOrdersHandler(
	readModel AllOrdersReadModel,
	logger *logrus.Entry,
) AllOrdersHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllOrders, []*orders.Order](
		allOrdersHandler{readModel: readModel},
		logger,
	)

}

type AllOrdersReadModel interface {
	GetAllOrders(ctx context.Context) ([]*orders.Order, error)
}

func (h allOrdersHandler) Handle(ctx context.Context, q AllOrders) (o []*orders.Order, err error) {
	return h.readModel.GetAllOrders(ctx)
}
