package app

import (
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	PlaceOrder   command.PlaceOrderHandler
	ConfirmOrder command.ConfirmOrderHandler
	CancelOrder  command.CancelOrderHandler
}

type Queries struct {
	AllOrders query.AllOrdersHandler
	OrderById query.OrderByIdHandler
}
