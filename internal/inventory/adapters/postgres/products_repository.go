package postgres

import (
	"context"
	"fmt"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/products"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductsRepository(db *pgxpool.Pool) *ProductsRepository {
	if db == nil {
		panic("missing db")
	}

	return &ProductsRepository{db: db}
}

func (p ProductsRepository) AddProduct(ctx context.Context, product *products.Product) error {
	query := `INSERT INTO products (product_uuid, name, price, weight) VALUES ($1, $2, $3, $4)`

	_, err := p.db.Exec(ctx, query,
		product.ProductUUID,
		product.Name,
		product.Price,
		product.Weight,
	)

	if err != nil {
		return fmt.Errorf("unable to insert product: %w", err)
	}
	return nil
}

func (p ProductsRepository) GetAllProducts(ctx context.Context) ([]*products.Product, error) {
	query := `SELECT product_uuid, name, price, weight, created_at, updated_at FROM products`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query products: %w", err)
	}
	defer rows.Close()

	var productsList []*products.Product
	for rows.Next() {
		var p products.Product
		err := rows.Scan(&p.ProductUUID, &p.Name, &p.Price, &p.Weight, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan product: %w", err)
		}
		productsList = append(productsList, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %w", err)
	}

	return productsList, nil
}

func (p ProductsRepository) GetProduct(ctx context.Context, productUUID string) (*products.Product, error) {
	query := `SELECT product_uuid, name, price, weight, created_at, updated_at FROM products WHERE product_uuid = $1`

	var pr products.Product
	err := p.db.QueryRow(ctx, query, productUUID).Scan(&pr.ProductUUID, &pr.Name, &pr.Price, &pr.Weight, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to query product: %w", err)
	}

	return &pr, nil
}

func (p ProductsRepository) UpdateProduct(ctx context.Context, productUUID string, updateFunc func(ctx context.Context, p *products.Product) (*products.Product, error)) error {
	query := `SELECT product_uuid, name, price, weight, created_at, updated_at FROM products WHERE product_uuid = $1`

	var pr products.Product
	err := p.db.QueryRow(ctx, query, productUUID).Scan(&pr.ProductUUID, &pr.Name, &pr.Price, &pr.Weight, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to query product: %w", err)
	}

	updatedProduct, err := updateFunc(ctx, &pr)
	if err != nil {
		return fmt.Errorf("unable to update product: %w", err)
	}

	updateQuery := `UPDATE products SET name = $1, price = $2, weight = $3, updated_at = $4 WHERE product_uuid = $5`

	_, err = p.db.Exec(ctx, updateQuery, updatedProduct.Name, updatedProduct.Price, updatedProduct.Weight, productUUID)
	if err != nil {
		return fmt.Errorf("unable to update product: %w", err)
	}

	return nil
}
