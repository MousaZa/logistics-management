package postgres

import (
	"context"
	"fmt"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLocationsRepository struct {
	db *pgxpool.Pool
}

func (p PostgresLocationsRepository) AddLocation(ctx context.Context, location *locations.Location) error {
	query := `INSERT INTO locations (location_uuid, name, address, city) VALUES ($1, $2, $3, $4)`
	_, err := p.db.Exec(ctx, query, location.LocationUUID, location.Name, location.Address, location.City)
	if err != nil {
		return fmt.Errorf("unable to insert location: %w", err)
	}
	return nil
}

func (p PostgresLocationsRepository) GetAllLocations(ctx context.Context) ([]*locations.Location, error) {
	query := `SELECT location_uuid, name, address, city, created_at, updated_at FROM locations`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query locations: %w", err)
	}
	defer rows.Close()

	var locationsList []*locations.Location
	for rows.Next() {
		var l locations.Location
		err := rows.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan location: %w", err)
		}
		locationsList = append(locationsList, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over locations: %w", err)
	}

	return locationsList, nil
}

func (p PostgresLocationsRepository) GetLocation(ctx context.Context, locationUUID string) (*locations.Location, error) {
	query := `SELECT location_uuid, name, address, city, created_at, updated_at FROM locations WHERE location_uuid = $1`

	row := p.db.QueryRow(ctx, query, locationUUID)

	var l locations.Location
	err := row.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to scan location: %w", err)
	}

	return &l, nil
}

func (p PostgresLocationsRepository) UpdateLocation(ctx context.Context, locationUUID string, updateFunc func(ctx context.Context, p *locations.Location) (*locations.Location, error)) error {
	query := `SELECT location_uuid, name, address, city, created_at, updated_at FROM locations WHERE location_uuid = $1`

	row := p.db.QueryRow(ctx, query, locationUUID)

	var l locations.Location
	err := row.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to scan location: %w", err)
	}

	updatedLocation, err := updateFunc(ctx, &l)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	updateQuery := `UPDATE locations SET name = $1, address = $2, city = $3, updated_at = NOW() WHERE location_uuid = $4`
	_, err = p.db.Exec(ctx, updateQuery, updatedLocation.Name, updatedLocation.Address, updatedLocation.City, locationUUID)
	if err != nil {
		return fmt.Errorf("unable to update location: %w", err)
	}

	return nil
}

func NewPostgresLocationsRepository(db *pgxpool.Pool) *PostgresLocationsRepository {
	if db == nil {
		panic("missing db")
	}

	return &PostgresLocationsRepository{db: db}
}
