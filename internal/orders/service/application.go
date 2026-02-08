package service

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/orders/adapters"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
)

func NewApplication(_ context.Context) app.Application {

	ordersRepository := adapters.NewOrderMemoryRepository()

	return app.Application{
		Commands: app.Commands{
			PlaceOrder:   command.NewPlaceOrderHandler(nil, ordersRepository),
			ConfirmOrder: command.NewConfirmOrderHandler(nil, ordersRepository),
		},
		Queries: app.Queries{
			AllOrders: query.NewAllOrdersHandler(ordersRepository),
		},
	}
}
