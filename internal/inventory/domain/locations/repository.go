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

type LocationRepository interface {
	Save(ctx context.Context, location *Location) error
	FindAll(ctx context.Context) ([]*Location, error)
	FindByUUID(ctx context.Context, locationUUID string) (*Location, error)
	Update(
		ctx context.Context,
		locationUUID string,
		updateFunc func(
			ctx context.Context,
			o *Location,
		) (*Location, error)) error
}
