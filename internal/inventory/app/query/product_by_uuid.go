package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/sirupsen/logrus"
)

type ProductByUUID struct {
	ProductUUID string
}

type productByUUIDHandler struct {
	readModel ProductByUUIDReadModel
}

type ProductByUUIDHandler decorator.QueryHandler[ProductByUUID, *products.Product]

func NewProductByUUIDHandler(
	readModel ProductByUUIDReadModel,
	logger *logrus.Entry,
) ProductByUUIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[ProductByUUID, *products.Product](
		productByUUIDHandler{readModel: readModel},
		logger,
	)

}

type ProductByUUIDReadModel interface {
	GetProduct(ctx context.Context, productUUID string) (*products.Product, error)
}

func (h productByUUIDHandler) Handle(ctx context.Context, pid ProductByUUID) (o *products.Product, err error) {
	return h.readModel.GetProduct(ctx, pid.ProductUUID)
}
