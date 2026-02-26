package service

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/adapters/postgres"
	"github.com/MousaZa/logistics-management/internal/inventory/app"
	"github.com/MousaZa/logistics-management/internal/inventory/app/command"
	"github.com/MousaZa/logistics-management/internal/inventory/app/query"
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
		},
	}
}
