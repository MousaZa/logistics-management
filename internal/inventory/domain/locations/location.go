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
	Longitude    float32
	Latitude     float32
	UpdatedAt    time.Time
}

func NewLocation(name string, address string, city string, longitude float32, latitude float32) (*Location, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	if address == "" {
		return nil, errors.New("empty address")
	}
	if city == "" {
		return nil, errors.New("empty city")
	}
	if longitude == 0 {
		return nil, errors.New("empty longitude")
	}
	if latitude == 0 {
		return nil, errors.New("empty latitude")
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
		Longitude:    longitude,
		Latitude:     latitude,
	}, nil
}
