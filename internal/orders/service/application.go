package service

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/orders/adapters/postgres"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {

	//ordersRepository := adapters.NewOrderMemoryRepository()

	conn, err := postgres.NewPostgresConnection(ctx)
	if err != nil {
		panic(err)
	}

	ordersRepository := postgres.NewPostgresOrderRepository(conn)

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
