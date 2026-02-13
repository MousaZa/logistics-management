package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/sirupsen/logrus"
)

type OrderById struct {
	OrderId string
}

type orderByIdHandler struct {
	readModel OrderByIdReadModel
}

type OrderByIdHandler decorator.QueryHandler[OrderById, *orders.Order]

func NewOrderByIdHandler(
	readModel OrderByIdReadModel,
	logger *logrus.Entry,
) OrderByIdHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[OrderById, *orders.Order](
		orderByIdHandler{readModel: readModel},
		logger,
	)

}

type OrderByIdReadModel interface {
	GetOrder(ctx context.Context, orderUUID string) (*orders.Order, error)
}

func (h orderByIdHandler) Handle(ctx context.Context, oid OrderById) (o *orders.Order, err error) {
	return h.readModel.GetOrder(ctx, oid.OrderId)
}
