package ports

import (
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
	//TODO implement me
	panic("implement me")
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}
