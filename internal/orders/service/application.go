package service

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/orders/adapters"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {

	ordersRepository := adapters.NewOrderMemoryRepository()

	logger := logrus.NewEntry(logrus.StandardLogger())

	return app.Application{
		Commands: app.Commands{
			PlaceOrder:   command.NewPlaceOrderHandler(nil, ordersRepository, logger),
			ConfirmOrder: command.NewConfirmOrderHandler(nil, ordersRepository, logger),
		},
		Queries: app.Queries{
			AllOrders: query.NewAllOrdersHandler(ordersRepository, logger),
			OrderById: query.NewOrderByIdHandler(ordersRepository, logger),
		},
	}
}
