package postgres

import (
	"context"

	"github.com/MousaZa/logistics-management/internal/inventory/domain/locations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLocationsRepository struct {
	db *pgxpool.Pool
}

func (p PostgresLocationsRepository) AddLocation(ctx context.Context, location *locations.Location) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresLocationsRepository) GetAllLocations(ctx context.Context) ([]*locations.Location, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresLocationsRepository) GetLocation(ctx context.Context, locationUUID string) (*locations.Location, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresLocationsRepository) UpdateLocation(ctx context.Context, locationUUID string, updateFunc func(ctx context.Context, p *locations.Location) (*locations.Location, error)) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresLocationsRepository(db *pgxpool.Pool) *PostgresLocationsRepository {
	if db == nil {
		panic("missing db")
	}

	return &PostgresLocationsRepository{db: db}
}
