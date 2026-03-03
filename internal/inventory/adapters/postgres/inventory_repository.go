package postgres

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func (p InventoryRepository) GetInventory(ctx context.Context, productUUID, locationUUID string) (*inventory.Inventory, error) {
	query := `SELECT product_uuid, location_uuid, quantity FROM inventory WHERE product_uuid = $1 AND location_uuid = $2`

	row := p.db.QueryRow(ctx, query, productUUID, locationUUID)

	var inv inventory.Inventory
	err := row.Scan(&inv.ProductUUID, &inv.LocationUUID, &inv.Quantity)
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

func (p InventoryRepository) UpdateMultipleInventories(ctx context.Context, inventories ...*inventory.Inventory) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Safe to call even if committed

	query := `
        INSERT INTO inventory (product_uuid, location_uuid, quantity) 
        VALUES ($1, $2, $3)
        ON CONFLICT (product_uuid, location_uuid) 
        DO UPDATE SET quantity = EXCLUDED.quantity;
    `

	for _, item := range inventories {
		_, err := tx.Exec(ctx, query, item.ProductUUID, item.LocationUUID, item.Quantity)
		if err != nil {
			return err // This triggers the rollback
		}
	}

	return tx.Commit(ctx)
}

func (p InventoryRepository) GetLocationProducts(ctx context.Context, locationUUID string) ([]*products.ProductStock, error) {
	query := `SELECT p.product_uuid, p.name, p.price, p.weight, p.created_at, p.updated_at, i.quantity
			  FROM inventory i
			  JOIN products p ON i.product_uuid = p.product_uuid
			  WHERE i.location_uuid = $1`

	rows, err := p.db.Query(ctx, query, locationUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsList []*products.ProductStock
	for rows.Next() {
		var p products.ProductStock
		err := rows.Scan(&p.ProductUUID, &p.Name, &p.Price, &p.Weight, &p.CreatedAt, &p.UpdatedAt, &p.Quantity)
		if err != nil {
			return nil, err
		}
		productsList = append(productsList, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productsList, nil
}

func (p InventoryRepository) GetProductLocations(ctx context.Context, productUUID string) ([]*locations.ProductLocationInventory, error) {
	query := `SELECT l.location_uuid, l.name, l.address, l.city, l.created_at, l.updated_at, i.quantity
			  FROM inventory i
			  JOIN locations l ON i.location_uuid = l.location_uuid
			  WHERE i.product_uuid = $1`

	rows, err := p.db.Query(ctx, query, productUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locationsList []*locations.ProductLocationInventory
	for rows.Next() {
		var l locations.ProductLocationInventory
		err := rows.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.CreatedAt, &l.UpdatedAt, &l.Quantity)
		if err != nil {
			return nil, err
		}
		locationsList = append(locationsList, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locationsList, nil
}

func (p InventoryRepository) AddInventory(ctx context.Context, inventory *inventory.Inventory) error {
	query := `INSERT INTO inventory (location_uuid, product_uuid, quantity, status) VALUES ($1, $2, $3, $4)`

	_, err := p.db.Exec(ctx, query,
		inventory.LocationUUID,
		inventory.ProductUUID,
		inventory.Quantity,
		inventory.Status,
	)

	if err != nil {
		return err
	}
	return nil
}

func NewPostgresInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	if db == nil {
		panic("missing db")
	}

	return &InventoryRepository{db: db}
}
