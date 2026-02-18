package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(db *pgxpool.Pool) *PostgresOrderRepository {
	if db == nil {
		panic("missing db")
	}

	return &PostgresOrderRepository{db: db}
}

// AddOrder inserts a new order into the database
func (r *PostgresOrderRepository) AddOrder(ctx context.Context, o *orders.Order) error {
	lineItemsJSON, err := json.Marshal(o.LineItems)
	if err != nil {
		return fmt.Errorf("unable to marshal line items: %w", err)
	}

	query := `
		INSERT INTO orders (
			order_uuid, placed_by, order_total, weight, destination, 
			status, line_items, ordered_date, completed_date, delivered_date, shipped_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err = r.db.Exec(ctx, query,
		o.OrderUUID,
		o.PlacedBy,
		o.OrderTotal,
		o.Weight,
		o.Destination,
		o.Status,
		lineItemsJSON,
		o.OrderedDate,
		o.CompletedDate,
		o.DeliveredDate,
		o.ShippedDate,
	)

	if err != nil {
		return fmt.Errorf("unable to insert order: %w", err)
	}

	return nil
}

// GetOrder retrieves a single order by UUID
func (r *PostgresOrderRepository) GetOrder(ctx context.Context, orderUUID string) (*orders.Order, error) {
	query := `
		SELECT 
			order_uuid, placed_by, order_total, weight, destination, 
			status, line_items, ordered_date, completed_date, delivered_date, shipped_date
		FROM orders
		WHERE order_uuid = $1
	`

	var lineItemsJSON []byte
	order := &orders.Order{}

	err := r.db.QueryRow(ctx, query, orderUUID).Scan(
		&order.OrderUUID,
		&order.PlacedBy,
		&order.OrderTotal,
		&order.Weight,
		&order.Destination,
		&order.Status,
		&lineItemsJSON,
		&order.OrderedDate,
		&order.CompletedDate,
		&order.DeliveredDate,
		&order.ShippedDate,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, orders.NotFoundError{OrderUUID: orderUUID}
	}

	if err != nil {
		return nil, fmt.Errorf("unable to get order: %w", err)
	}

	err = json.Unmarshal(lineItemsJSON, &order.LineItems)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal line items: %w", err)
	}

	return order, nil
}

// GetAllOrders retrieves all orders from the database
func (r *PostgresOrderRepository) GetAllOrders(ctx context.Context) ([]*orders.Order, error) {
	query := `
		SELECT 
			order_uuid, placed_by, order_total, weight, destination, 
			status, line_items, ordered_date, completed_date, delivered_date, shipped_date
		FROM orders
		ORDER BY ordered_date DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query orders: %w", err)
	}
	defer rows.Close()

	var result []*orders.Order

	for rows.Next() {
		var lineItemsJSON []byte
		order := &orders.Order{}

		err := rows.Scan(
			&order.OrderUUID,
			&order.PlacedBy,
			&order.OrderTotal,
			&order.Weight,
			&order.Destination,
			&order.Status,
			&lineItemsJSON,
			&order.OrderedDate,
			&order.CompletedDate,
			&order.DeliveredDate,
			&order.ShippedDate,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to scan order: %w", err)
		}

		err = json.Unmarshal(lineItemsJSON, &order.LineItems)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal line items: %w", err)
		}

		result = append(result, order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return result, nil
}

// UpdateOrder updates an existing order using a transaction
func (r *PostgresOrderRepository) UpdateOrder(
	ctx context.Context,
	orderUUID string,
	updateFunc func(ctx context.Context, o *orders.Order) (*orders.Order, error),
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Get existing order for update
	query := `
		SELECT 
			order_uuid, placed_by, order_total, weight, destination, 
			status, line_items, ordered_date, completed_date, delivered_date, shipped_date
		FROM orders
		WHERE order_uuid = $1
		FOR UPDATE
	`

	var lineItemsJSON []byte
	order := &orders.Order{}

	err = tx.QueryRow(ctx, query, orderUUID).Scan(
		&order.OrderUUID,
		&order.PlacedBy,
		&order.OrderTotal,
		&order.Weight,
		&order.Destination,
		&order.Status,
		&lineItemsJSON,
		&order.OrderedDate,
		&order.CompletedDate,
		&order.DeliveredDate,
		&order.ShippedDate,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return orders.NotFoundError{OrderUUID: orderUUID}
	}

	if err != nil {
		return fmt.Errorf("unable to get order: %w", err)
	}

	err = json.Unmarshal(lineItemsJSON, &order.LineItems)
	if err != nil {
		return fmt.Errorf("unable to unmarshal line items: %w", err)
	}

	// Apply update function
	updatedOrder, err := updateFunc(ctx, order)
	if err != nil {
		return err
	}

	// Update order in database
	updatedLineItemsJSON, err := json.Marshal(updatedOrder.LineItems)
	if err != nil {
		return fmt.Errorf("unable to marshal line items: %w", err)
	}

	updateQuery := `
		UPDATE orders SET
			placed_by = $1,
			order_total = $2,
			weight = $3,
			destination = $4,
			status = $5,
			line_items = $6,
			ordered_date = $7,
			completed_date = $8,
			delivered_date = $9,
			shipped_date = $10
		WHERE order_uuid = $11
	`

	_, err = tx.Exec(ctx, updateQuery,
		updatedOrder.PlacedBy,
		updatedOrder.OrderTotal,
		updatedOrder.Weight,
		updatedOrder.Destination,
		updatedOrder.Status,
		updatedLineItemsJSON,
		updatedOrder.OrderedDate,
		updatedOrder.CompletedDate,
		updatedOrder.DeliveredDate,
		updatedOrder.ShippedDate,
		orderUUID,
	)

	if err != nil {
		return fmt.Errorf("unable to update order: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

// NewPostgresConnection creates a new PostgreSQL connection pool
func NewPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
