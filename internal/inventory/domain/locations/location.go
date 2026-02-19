package locations

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Address      string
	City         string
	CreatedAt    time.Time
	LocationUUID string
	Name         string
	UpdatedAt    time.Time
}

func NewLocation(name string, address string, city string) (*Location, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	if address == "" {
		return nil, errors.New("empty address")
	}
	if city == "" {
		return nil, errors.New("empty city")
	}

	locationUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Location{
		LocationUUID: locationUUID.String(),
		Name:         name,
		Address:      address,
		City:         city,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
