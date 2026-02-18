include .env
export

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh orders internal/orders/ports ports
	@./scripts/openapi-http.sh inventory internal/inventory/ports ports

.PHONY: proto
proto:
	@./scripts/proto.sh orders
