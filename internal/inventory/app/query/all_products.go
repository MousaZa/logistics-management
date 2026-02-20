package query

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/sirupsen/logrus"
)

type AllProducts struct{}

type allProductsHandler struct {
	readModel AllProductsReadModel
	logger    *logrus.Entry
}

type AllProductsHandler decorator.QueryHandler[AllProducts, []*products.Product]

func NewAllProductsHandler(
	readModel AllProductsReadModel,
	logger *logrus.Entry,
) AllProductsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllProducts, []*products.Product](
		allProductsHandler{readModel: readModel, logger: logger},
		logger,
	)
}

func (h allProductsHandler) Handle(ctx context.Context, _ AllProducts) ([]*products.Product, error) {
	return h.readModel.GetAllProducts(ctx)
}

type AllProductsReadModel interface {
	GetAllProducts(ctx context.Context) ([]*products.Product, error)
}
