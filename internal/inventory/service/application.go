package service

import (
	"context"
	"fmt"
	"time"

	kafka_adapter "github.com/MousaZa/logistics-management/internal/common/adapters/kafka"
	"github.com/MousaZa/logistics-management/internal/inventory/adapters/postgres"
	"github.com/MousaZa/logistics-management/internal/inventory/app"
	"github.com/MousaZa/logistics-management/internal/inventory/app/command"
	"github.com/MousaZa/logistics-management/internal/inventory/app/events"
	"github.com/MousaZa/logistics-management/internal/inventory/app/query"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {

	conn, err := postgres.NewPostgresConnection(ctx)
	if err != nil {
		panic(err)
	}

	productsRepository := postgres.NewPostgresProductsRepository(conn)
	locationsRepository := postgres.NewPostgresLocationsRepository(conn)
	inventoryRepository := postgres.NewPostgresInventoryRepository(conn)

	logger := logrus.NewEntry(logrus.StandardLogger())

	watermillLogger := watermill.NewStdLogger(true, true)

	sub, err := kafka_adapter.NewSubscriber([]string{"kafka:29092"}, "inventory_service_group", watermillLogger)
	if err != nil {
		panic(err)
	}

	// 2. Setup the Router
	// The router acts as the engine that pulls messages from the subscriber
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)

	// Add standard middleware (retries, poison queue, logging)
	router.AddMiddleware(
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
		}.Middleware,
		middleware.Recoverer,
	)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			// How do we generate the topic name to listen to?
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				// This must match the topic the Order service published to!
				return "events.OrderPlacedEvent", nil
			},
			SubscriberConstructor: func(e cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return sub, nil
			},
			// How do we know which event is which?
			Marshaler: cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
			Logger:    watermillLogger,
		},
	)

	// 4. Register the App Handler with the Processor
	handler := events.NewExternalOrderPlacedHandler(inventoryRepository)
	err = eventProcessor.AddHandlers(handler)
	if err != nil {
		panic(err)
	}

	// 5. Start the engine (this blocks and runs forever)
	fmt.Println("Starting Inventory Event Router...")
	go func() {
		err := router.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	return app.Application{
		Commands: app.Commands{
			AddLocation:    command.NewAddLocationHandler(locationsRepository, logger),
			UpdateLocation: command.NewUpdateLocationHandler(locationsRepository, logger),

			AddProduct:    command.NewAddProductHandler(productsRepository, logger),
			UpdateProduct: command.NewUpdateProductHandler(productsRepository, logger),

			AddInventory:    command.NewAddInventoryHandler(inventoryRepository, logger),
			TransferProduct: command.NewTransferProductHandler(inventoryRepository, logger),
		},
		Queries: app.Queries{
			LocationByUUID: query.NewLocationByUUIDHandler(locationsRepository, logger),
			AllLocations:   query.NewAllLocationsHandler(locationsRepository, logger),

			ProductByUUID: query.NewProductByUUIDHandler(productsRepository, logger),
			AllProducts:   query.NewAllProductsHandler(productsRepository, logger),

			LocationProducts: query.NewLocationProductsHandler(inventoryRepository, logger),
			ProductLocations: query.NewProductLocationsHandler(inventoryRepository, logger),
		},
	}
}
