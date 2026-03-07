package postgres

import (
	"context"
	"fmt"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/inventory"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func (p InventoryRepository) UpdateInventory(ctx context.Context, locationUUID string, productUUID string, updateFunc func(ctx context.Context, i *inventory.Inventory) (*inventory.Inventory, error)) error {
	query := `SELECT product_uuid, location_uuid, qty_available, qty_damaged, qty_reserved FROM inventory WHERE product_uuid = $1 AND location_uuid = $2`

	row := p.db.QueryRow(ctx, query, productUUID, locationUUID)

	var i inventory.Inventory
	err := row.Scan(&i.ProductUUID, &i.LocationUUID, &i.AvailableQuantity, &i.DamagedQuantity, &i.ReservedQuantity)
	if err != nil {
		return fmt.Errorf("unable to scan inventory: %w", err)
	}

	updatedInventory, err := updateFunc(ctx, &i)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	updateQuery := `UPDATE inventory SET qty_available = $1, qty_reserved=$2, qty_damaged=$3 WHERE product_uuid = $4 and location_uuid = $5`
	_, err = p.db.Exec(ctx, updateQuery, updatedInventory.AvailableQuantity, updatedInventory.ReservedQuantity, updatedInventory.DamagedQuantity, productUUID, locationUUID)
	if err != nil {
		return fmt.Errorf("unable to update inventory: %w", err)
	}

	return nil
}

func (p InventoryRepository) GetInventory(ctx context.Context, productUUID, locationUUID string) (*inventory.Inventory, error) {
	query := `SELECT product_uuid, location_uuid, qty_available, qty_damaged, qty_reserved FROM inventory WHERE product_uuid = $1 AND location_uuid = $2`

	row := p.db.QueryRow(ctx, query, productUUID, locationUUID)

	var inv inventory.Inventory
	err := row.Scan(&inv.ProductUUID, &inv.LocationUUID, &inv.AvailableQuantity, &inv.DamagedQuantity, &inv.ReservedQuantity)
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
        INSERT INTO inventory (product_uuid, location_uuid, qty_available) 
        VALUES ($1, $2, $3)
        ON CONFLICT (product_uuid, location_uuid) 
        DO UPDATE SET quantity = EXCLUDED.quantity;
    `

	for _, item := range inventories {
		_, err := tx.Exec(ctx, query, item.ProductUUID, item.LocationUUID, item.AvailableQuantity)
		if err != nil {
			return err // This triggers the rollback
		}
	}

	return tx.Commit(ctx)
}

func (p InventoryRepository) GetLocationProducts(ctx context.Context, locationUUID string) ([]*products.ProductStock, error) {
	query := `SELECT p.product_uuid, p.name, p.price, p.weight, p.created_at, p.updated_at, i.qty_available
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
	query := `SELECT l.location_uuid, l.name, l.address, l.city, l.created_at, l.updated_at, i.qty_available, i.qty_damaged, i.qty_reserved
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
		err := rows.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.CreatedAt, &l.UpdatedAt, &l.AvailableQuantity, &l.DamagedQuantity, &l.ReservedQuantity)
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
	query := `INSERT INTO inventory (location_uuid, product_uuid, qty_available, qty_damaged, qty_reserved) VALUES ($1, $2, $3, $4, $5)`

	_, err := p.db.Exec(ctx, query,
		inventory.LocationUUID,
		inventory.ProductUUID,
		inventory.AvailableQuantity,
		inventory.DamagedQuantity,
		inventory.ReservedQuantity,
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
