package postgres

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInventoryRepository struct {
	db *pgxpool.Pool
}

func (p PostgresInventoryRepository) GetLocationProducts(ctx context.Context, locationUUID string) ([]*products.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresInventoryRepository) GetProductLocations(ctx context.Context, productUUID string) ([]*locations.Location, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresInventoryRepository) AddInventory(ctx context.Context, inventory *inventory.Inventory) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresInventoryRepository(db *pgxpool.Pool) *PostgresInventoryRepository {
	if db == nil {
		panic("missing db")
	}

	return &PostgresInventoryRepository{db: db}
}
