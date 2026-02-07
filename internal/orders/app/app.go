package app

import "github.com/MousaZa/logistics-management/internal/orders/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	NewOrder command.PlaceOrderHandler
}

type Queries struct {
}
