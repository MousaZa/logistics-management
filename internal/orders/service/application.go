package service

import (
	"context"

	kafka_adapter "github.com/MousaZa/logistics-management/internal/common/adapters/kafka"
	"github.com/MousaZa/logistics-management/internal/orders/adapters/postgres"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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

	watermillLogger := watermill.NewStdLogger(true, true)

	publisher, err := kafka_adapter.NewPublisher([]string{"kafka:29092"}, watermillLogger)
	if err != nil {
		panic(err)
	}

	eventBus, err := cqrs.NewEventBusWithConfig(publisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return "events." + params.EventName, nil
		},
		Marshaler: cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		Logger:    watermillLogger,
	})

	return app.Application{
		Commands: app.Commands{
			PlaceOrder:   command.NewPlaceOrderHandler(eventBus, ordersRepository, logger),
			ConfirmOrder: command.NewConfirmOrderHandler(nil, ordersRepository, logger),
			CancelOrder:  command.NewCancelOrderHandler(eventBus, ordersRepository, logger),
		},
		Queries: app.Queries{
			AllOrders: query.NewAllOrdersHandler(ordersRepository, logger),
			OrderById: query.NewOrderByIdHandler(ordersRepository, logger),
		},
	}
}
