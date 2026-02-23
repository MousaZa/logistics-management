package main

import (
	"context"
	"net/http"

	"github.com/MousaZa/logistics-management/internal/common/logs"
	"github.com/MousaZa/logistics-management/internal/common/server"
	"github.com/MousaZa/logistics-management/internal/inventory/ports"
	"github.com/MousaZa/logistics-management/internal/inventory/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHttpServer(application),
			router,
		)
	})
}
