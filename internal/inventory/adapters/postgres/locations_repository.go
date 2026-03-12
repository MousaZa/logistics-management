package postgres

import (
	"context"
	"fmt"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationsRepository struct {
	db *pgxpool.Pool
}

func (p LocationsRepository) AddLocation(ctx context.Context, location *locations.Location) error {
	query := `INSERT INTO locations (location_uuid, name, address, city, coordinates) VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography)`
	_, err := p.db.Exec(ctx, query, location.LocationUUID, location.Name, location.Address, location.City, location.Longitude, location.Latitude)
	if err != nil {
		return fmt.Errorf("unable to insert locations: %w", err)
	}
	return nil
}

func (p LocationsRepository) GetAllLocations(ctx context.Context) ([]*locations.Location, error) {
	query := `SELECT location_uuid, name, address, city,ST_Y(coordinates::geometry) AS longitude, 
    ST_X(coordinates::geometry) AS latitude, created_at, updated_at FROM locations`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query locations: %w", err)
	}
	defer rows.Close()

	var locationsList []*locations.Location
	for rows.Next() {
		var l locations.Location
		err := rows.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.Longitude, &l.Latitude, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan locations: %w", err)
		}
		locationsList = append(locationsList, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over locations: %w", err)
	}

	return locationsList, nil
}

func (p LocationsRepository) GetLocation(ctx context.Context, locationUUID string) (*locations.Location, error) {
	query := `SELECT location_uuid, name, address, city,ST_Y(coordinates::geometry) AS longitude, 
    ST_X(coordinates::geometry) AS latitude, created_at, updated_at FROM locations WHERE location_uuid = $1`

	row := p.db.QueryRow(ctx, query, locationUUID)

	var l locations.Location
	err := row.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.Longitude, &l.Latitude, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to scan locations: %w", err)
	}

	return &l, nil
}

func (p LocationsRepository) UpdateLocation(ctx context.Context, locationUUID string, updateFunc func(ctx context.Context, p *locations.Location) (*locations.Location, error)) error {
	query := `SELECT location_uuid, name, address, city,ST_Y(coordinates::geometry) AS longitude, 
    ST_X(coordinates::geometry) AS latitude, created_at, updated_at FROM locations WHERE location_uuid = $1`

	row := p.db.QueryRow(ctx, query, locationUUID)

	var l locations.Location
	err := row.Scan(&l.LocationUUID, &l.Name, &l.Address, &l.City, &l.Longitude, &l.Latitude, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to scan locations: %w", err)
	}

	updatedLocation, err := updateFunc(ctx, &l)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	updateQuery := `UPDATE locations SET name = $1, address = $2, city = $3, coordinates = ST_SetSRID(ST_MakePoint($4, $5), 4326)::geography, updated_at = NOW() WHERE location_uuid = $6`
	_, err = p.db.Exec(ctx, updateQuery, updatedLocation.Name, updatedLocation.Address, updatedLocation.City, updatedLocation.Longitude, updatedLocation.Latitude, locationUUID)
	if err != nil {
		return fmt.Errorf("unable to update locations: %w", err)
	}

	return nil
}

func NewPostgresLocationsRepository(db *pgxpool.Pool) *LocationsRepository {
	if db == nil {
		panic("missing db")
	}

	return &LocationsRepository{db: db}
}
