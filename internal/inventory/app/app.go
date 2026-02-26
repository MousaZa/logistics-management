package app

import (
	"github.com/MousaZa/logistics-management/internal/inventory/app/command"
	"github.com/MousaZa/logistics-management/internal/inventory/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	// Locations
	AddLocation    command.AddLocationHandler
	UpdateLocation command.UpdateLocationHandler

	// Products
	AddProduct    command.AddProductHandler
	UpdateProduct command.UpdateProductHandler

	// Inventory
	AddInventory command.AddInventoryHandler
}

type Queries struct {
	// Locations
	LocationByUUID query.LocationByUUIDHandler
	AllLocations   query.AllLocationsHandler

	// Products
	ProductByUUID query.ProductByUUIDHandler
	AllProducts   query.AllProductsHandler

	// Inventory
	LocationProducts query.LocationProductsHandler
	ProductLocations query.ProductLocationsHandler
}
