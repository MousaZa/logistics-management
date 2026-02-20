package locations

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	LocationUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("location '%s' not found", e.LocationUUID)
}

type Repository interface {
	AddLocation(ctx context.Context, location *Location) error
	GetAllLocations(ctx context.Context) ([]*Location, error)
	GetLocation(ctx context.Context, locationUUID string) (*Location, error)
	UpdateLocation(
		ctx context.Context,
		locationUUID string,
		updateFunc func(
			ctx context.Context,
			o *Location,
		) (*Location, error)) error
}
