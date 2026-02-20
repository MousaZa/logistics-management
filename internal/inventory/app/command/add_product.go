package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/sirupsen/logrus"
)

type AddProduct struct {
	Name              string
	Price             float32
	Weight            float32
	LocationUUID      string
	AvailableQuantity int
}

type AddProductHandler decorator.CommandHandler[AddProduct]

type addProductHandler struct {
	repo products.Repository
}

func NewAddProductHandler(repo products.Repository, logger *logrus.Entry) AddProductHandler {
	return decorator.ApplyCommandDecorators[AddProduct](addProductHandler{repo: repo}, logger)
}

func (h addProductHandler) Handle(ctx context.Context, cmd AddProduct) error {
	product, err := products.NewProduct(cmd.Name, cmd.Price, cmd.Weight, cmd.LocationUUID, cmd.AvailableQuantity)
	if err != nil {
		return err
	}

	err = h.repo.AddProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
