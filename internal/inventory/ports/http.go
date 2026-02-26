package ports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MousaZa/logistics-management/internal/common/server/httperr"
	"github.com/MousaZa/logistics-management/internal/inventory/app"
	"github.com/MousaZa/logistics-management/internal/inventory/app/command"
	"github.com/MousaZa/logistics-management/internal/inventory/app/query"
	"github.com/go-chi/render"
	"github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) GetLocationByUUID(w http.ResponseWriter, r *http.Request, locationUUID types.UUID) {
	uuids := locationUUID.String()

	l, err := h.app.Queries.LocationByUUID.Handle(r.Context(), query.LocationByUUID{LocationUUID: uuids})
	if err != nil {
		httperr.InternalError("unable-to-get-location", err, w, r)
		return
	}
	if l == nil {
		httperr.InternalError("unable-to-get-location", err, w, r)
		return
	}

	resp := Location{
		LocationUUID: &l.LocationUUID,
		Name:         l.Name,
		Address:      l.Address,
		City:         l.City,
		CreatedAt:    &l.CreatedAt,
		UpdatedAt:    &l.UpdatedAt,
	}

	render.Respond(w, r, resp)
}

func (h HttpServer) UpdateLocation(w http.ResponseWriter, r *http.Request, locationUUID types.UUID) {
	uuids := locationUUID.String()
	loc := Location{}

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &loc)
	if err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return
	}

	err = h.app.Commands.UpdateLocation.Handle(r.Context(), command.UpdateLocation{LocationUUID: uuids, Name: &loc.Name, Address: &loc.Address, City: &loc.City})
	if err != nil {
		httperr.InternalError("unable-to-update-location", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AddProductsToLocationRequest struct {
	ProductUUID types.UUID `json:"product_uuid"`
	Quantity    int        `json:"quantity"`
}

func (h HttpServer) AddProductsToLocation(w http.ResponseWriter, r *http.Request, locationUUID types.UUID) {
	req := &AddProductsToLocationRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return
	}

	err = h.app.Commands.AddInventory.Handle(r.Context(), command.AddInventory{LocationUUID: locationUUID.String(), ProductUUID: req.ProductUUID.String(), Quantity: req.Quantity})
	if err != nil {
		httperr.InternalError("unable-to-add-products-to-location", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetProductLocations(w http.ResponseWriter, r *http.Request, productUUID types.UUID) {
	locs, err := h.app.Queries.ProductLocations.Handle(r.Context(), query.ProductLocations{ProductUUID: productUUID.String()})
	if err != nil {
		httperr.InternalError("unable-to-get-product-locations", err, w, r)
		return
	}

	resp := make([]Location, len(locs))
	for i, loc := range locs {
		resp[i] = Location{
			LocationUUID: &loc.LocationUUID,
			Name:         loc.Name,
			Address:      loc.Address,
			City:         loc.City,
			CreatedAt:    &loc.CreatedAt,
			UpdatedAt:    &loc.UpdatedAt,
		}
	}

	render.Respond(w, r, resp)
}

func (h HttpServer) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.app.Queries.AllProducts.Handle(r.Context(), query.AllProducts{})
	if err != nil {
		httperr.InternalError("unable-to-get-products", err, w, r)
		return
	}

	resp := make([]Product, len(products))
	for i, p := range products {
		resp[i] = Product{
			ProductUUID: &p.ProductUUID,
			Name:        p.Name,
			Price:       p.Price,
			Weight:      p.Weight,
			CreatedAt:   &p.CreatedAt,
			UpdatedAt:   &p.UpdatedAt,
		}
	}

	render.Respond(w, r, resp)
}

func (h HttpServer) CreateProduct(w http.ResponseWriter, r *http.Request) {
	prod := &Product{}
	err := json.NewDecoder(r.Body).Decode(prod)
	if err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		fmt.Println(err)
		return
	}

	err = h.app.Commands.AddProduct.Handle(r.Context(), command.AddProduct{Name: prod.Name, Price: prod.Price, Weight: prod.Weight})
	if err != nil {
		httperr.InternalError("unable-to-create-product", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetProductByUUID(w http.ResponseWriter, r *http.Request, productUUID types.UUID) {
	product, err := h.app.Queries.ProductByUUID.Handle(r.Context(), query.ProductByUUID{ProductUUID: productUUID.String()})
	if err != nil {
		httperr.InternalError("unable-to-get-product", err, w, r)
		return
	}

	resp := Product{
		ProductUUID: &product.ProductUUID,
		Name:        product.Name,
		Price:       product.Price,
		Weight:      product.Weight,
		CreatedAt:   &product.CreatedAt,
		UpdatedAt:   &product.UpdatedAt,
	}

	render.Respond(w, r, resp)
}

func (h HttpServer) UpdateProduct(w http.ResponseWriter, r *http.Request, productUUID types.UUID) {
	prod := Product{}

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &prod)
	if err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return
	}

	err = h.app.Commands.UpdateProduct.Handle(r.Context(), command.UpdateProduct{ProductUUID: productUUID.String(), Name: &prod.Name, Price: &prod.Price, Weight: &prod.Weight})
	if err != nil {
		httperr.InternalError("unable-to-update-product", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetLocations(w http.ResponseWriter, r *http.Request) {
	l, err := h.app.Queries.AllLocations.Handle(r.Context(), query.AllLocations{})
	if err != nil {
		httperr.InternalError("unable-to-get-locations", err, w, r)
		return
	}
	if l == nil {
		httperr.InternalError("nil-locations", nil, w, r)
		return
	}
	resp := make([]Location, len(l))
	for i, loc := range l {
		resp[i] = Location{
			LocationUUID: &loc.LocationUUID,
			Name:         loc.Name,
			Address:      loc.Address,
			City:         loc.City,
			CreatedAt:    &loc.CreatedAt,
			UpdatedAt:    &loc.UpdatedAt,
		}
	}

	render.Respond(w, r, resp)
}

func (h HttpServer) CreateLocation(w http.ResponseWriter, r *http.Request) {
	loc := &Location{}
	err := render.Decode(r, loc)
	if err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return
	}

	err = h.app.Commands.AddLocation.Handle(r.Context(), command.AddLocation{Name: loc.Name, Address: loc.Address, City: loc.City})
	if err != nil {
		httperr.InternalError("unable-to-create-location", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	render.Respond(w, r, nil)
}

func (h HttpServer) TransferProducts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) GetLocationContents(w http.ResponseWriter, r *http.Request, locationUUID types.UUID) {
	prods, err := h.app.Queries.LocationProducts.Handle(r.Context(), query.LocationProducts{LocationUUID: locationUUID.String()})
	if err != nil {
		httperr.InternalError("unable-to-get-location-contents", err, w, r)
		return
	}

	resp := make([]Product, len(prods))
	for i, p := range prods {
		q := 999 // random quantity for now
		resp[i] = Product{
			ProductUUID: &p.ProductUUID,
			Name:        p.Name,
			Price:       p.Price,
			Weight:      p.Weight,
			CreatedAt:   &p.CreatedAt,
			UpdatedAt:   &p.UpdatedAt,
			Quantity:    &q, //TODO add quantity to product struct
		}
	}

	render.Respond(w, r, resp)
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}
