package service

import (
	"github.com/MousaZa/logistics-management/internal/inventory/app"
	"github.com/MousaZa/logistics-management/internal/inventory/app/command"
	"github.com/MousaZa/logistics-management/internal/inventory/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication() *app.Application {

	logger := logrus.NewEntry(logrus.StandardLogger())

	return &app.Application{
		Commands: app.Commands{
			AddLocation: command.NewAddLocationHandler(nil, logger),
			AddProduct:  command.NewAddProductHandler(nil, logger),
		},
		Queries: app.Queries{
			LocationByUUID:   query.NewLocationByUUIDHandler(nil, logger),
			AllLocations:     query.NewAllLocationsHandler(nil, logger),
			ProductByUUID:    query.NewProductByUUIDHandler(nil, logger),
			AllProducts:      query.NewAllProductsHandler(nil, logger),
			LocationProducts: query.NewLocationProductsHandler(nil, logger),
		},
	}
}
