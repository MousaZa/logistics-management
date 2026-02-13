package ports

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MousaZa/logistics-management/internal/common/server/httperr"
	"github.com/MousaZa/logistics-management/internal/orders/app"
	"github.com/MousaZa/logistics-management/internal/orders/app/command"
	"github.com/MousaZa/logistics-management/internal/orders/app/query"
	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	data, err := h.app.Queries.AllOrders.Handle(r.Context(), query.AllOrders{DateFrom: params.DateFrom, DateTo: params.DateTo})
	if err != nil {
		httperr.InternalError("unable-to-get", err, w, r)
	}
	or, err := orderDataToResponse(data)
	if err != nil {
		httperr.InternalError("unable-to-get", err, w, r)
	}

	render.Respond(w, r, or)
}

func orderDataToResponse(data []*orders.Order) (Orders, error) {
	var or Orders
	for _, order := range data {
		oid, err := uuid.FromBytes([]byte(order.OrderUUID))
		if err != nil {
			return Orders{}, err
		}
		placedBy, err := uuid.FromBytes([]byte(order.PlacedBy))
		var lis []LineItem
		for _, l := range order.LineItems {
			pid, err := uuid.FromBytes([]byte(l.ProductUUID))
			if err != nil {
				return Orders{}, err
			}
			li := LineItem{
				LineTotal:   l.LineTotal,
				Quantity:    l.Quantity,
				LineWeight:  l.LineWeight,
				ProductName: l.ProductName,
				ProductUUID: pid,
				UnitPrice:   l.UnitPrice,
				UnitWeight:  l.UnitWeight,
			}
			lis = append(lis, li)
		}
		or.Orders = append(or.Orders, Order{
			OrderUUID:     oid,
			PlacedBy:      placedBy,
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
		httperr.BadRequest("unable-to-read", err, w, r)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httperr.InternalError("unable-to-close", err, w, r)
		}
	}(r.Body)
	err = json.Unmarshal(d, &o)
	if err != nil {
		httperr.BadRequest("unable-to-read", err, w, r)
	}
	err = h.app.Commands.PlaceOrder.Handle(r.Context(), command.PlaceOrder{OrderTotal: o.OrderTotal, PlacedBy: o.PlacedBy.String(), Weight: o.Weight, Destination: o.Destination})
	if err != nil {
		httperr.InternalError("unable-to-create", err, w, r)
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
