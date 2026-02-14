package ports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MousaZa/logistics-management/internal/common/server/httperr"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/go-chi/render"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	if h.app.Queries.AllOrders == nil {
		httperr.InternalError("query-handler-not-initialized", nil, w, r)
		return
	}
	data, err := h.app.Queries.AllOrders.Handle(r.Context(), query.AllOrders{}) //TODO add params
	if err != nil {
		fmt.Println(err)
		httperr.InternalError("unable-to-get", err, w, r)
		return
	}
	if data == nil {
		render.Respond(w, r, []Order{})
		return
	}
	or, err := orderDataToResponse(data)
	if err != nil {
		fmt.Println(err)
		httperr.InternalError("unable-to-convert", err, w, r)
		return
	}
	render.Respond(w, r, or.Orders)
}

func orderDataToResponse(data []*orders.Order) (Orders, error) {
	var or Orders
	if len(data) == 0 {
		return or, nil
	}
	for _, order := range data {

		var lis []LineItem
		for _, l := range order.LineItems {

			li := LineItem{
				LineTotal:   l.LineTotal,
				Quantity:    l.Quantity,
				LineWeight:  l.LineWeight,
				ProductName: l.ProductName,
				ProductUUID: &l.ProductUUID,
				UnitPrice:   l.UnitPrice,
				UnitWeight:  l.UnitWeight,
			}
			lis = append(lis, li)
		}
		or.Orders = append(or.Orders, Order{
			OrderUUID:     &order.OrderUUID,
			PlacedBy:      order.PlacedBy,
			Destination:   order.Destination,
			OrderedDate:   &order.OrderedDate,
			OrderTotal:    order.OrderTotal,
			CompletedDate: &order.CompletedDate,
			DeliveredDate: &order.DeliveredDate,
			LineItems:     lis,
			Status:        OrderStatus(order.Status),
			Weight:        order.Weight,
		})
	}
	return or, nil
}

func (h HttpServer) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	d, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("here1")
		httperr.BadRequest("unable-to-read", err, w, r)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httperr.InternalError("unable-to-close", err, w, r)
			return
		}
	}(r.Body)
	err = json.Unmarshal(d, &o)
	if err != nil {
		fmt.Println(err)
		httperr.BadRequest("unable-to-marshal", err, w, r)
		return
	}

	var lineItems []orders.LineItem
	for _, li := range o.LineItems {
		lineItems = append(lineItems, orders.LineItem{
			ProductUUID: *li.ProductUUID,
			ProductName: li.ProductName,
			Quantity:    li.Quantity,
			UnitPrice:   li.UnitPrice,
			UnitWeight:  li.UnitWeight,
			LineTotal:   li.LineTotal,
			LineWeight:  li.LineWeight,
		})
	}
	err = h.app.Commands.PlaceOrder.Handle(r.Context(), command.PlaceOrder{
		LineItems:   lineItems,
		OrderTotal:  o.OrderTotal,
		PlacedBy:    o.PlacedBy,
		Weight:      o.Weight,
		Destination: o.Destination,
	})
	if err != nil {
		fmt.Println("here3")
		httperr.InternalError("unable-to-create", err, w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetOrderById(w http.ResponseWriter, r *http.Request, orderUUID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) CancelOrder(w http.ResponseWriter, r *http.Request, orderUUID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}
