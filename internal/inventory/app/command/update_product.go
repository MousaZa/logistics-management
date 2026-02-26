package command

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/common/decorator"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/sirupsen/logrus"
)

type UpdateProduct struct {
	ProductUUID string
	Name        *string
	Price       *float32
	Weight      *float32
}

type updateProductHandler struct {
	repo products.Repository
}

type UpdateProductHandler decorator.CommandHandler[UpdateProduct]

func NewUpdateProductHandler(repo products.Repository, logger *logrus.Entry) UpdateProductHandler {
	return decorator.ApplyCommandDecorators[UpdateProduct](updateProductHandler{repo: repo}, logger)
}

func (h updateProductHandler) Handle(ctx context.Context, cmd UpdateProduct) error {
	return h.repo.UpdateProduct(ctx, cmd.ProductUUID, func(ctx context.Context, p *products.Product) (*products.Product, error) {
		if cmd.Name != nil && *cmd.Name != "" {
			p.Name = *cmd.Name
		}
		if cmd.Price != nil && *cmd.Price != 0 {
			p.Price = *cmd.Price
		}
		if cmd.Weight != nil && *cmd.Weight != 0 {
			p.Weight = *cmd.Weight
		}
		return p, nil
	})
}
