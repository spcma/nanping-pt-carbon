package valueobject

import (
	"errors"
	"fmt"
	"math"
)

// Coordinate GPS coordinate value object
type Coordinate struct {
	Latitude   float64 // latitude
	Longtitude float64 // longitude
}

// NewCoordinate creates a coordinate
func NewCoordinate(latitude, longtitude float64) (*Coordinate, error) {
	if latitude < -90 || latitude > 90 {
		return nil, errors.New("latitude range: -90 to 90")
	}
	if longtitude < -180 || longtitude > 180 {
		return nil, errors.New("longitude range: -180 to 180")
	}

	return &Coordinate{
		Latitude:   latitude,
		Longtitude: longtitude,
	}, nil
}

// DistanceTo calculates distance to another coordinate (Haversine formula)
func (c *Coordinate) DistanceTo(other *Coordinate) float64 {
	if other == nil {
		return 0
	}

	const R = 6371e3 // Earth radius in meters

	lat1 := c.Latitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	deltaLat := (other.Latitude - c.Latitude) * math.Pi / 180
	deltaLon := (other.Longtitude - c.Longtitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	cVal := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * cVal
}

// String string representation
func (c *Coordinate) String() string {
	return fmt.Sprintf("(%.6f,%.6f)", c.Latitude, c.Longtitude)
}

// Equal checks if equal
func (c *Coordinate) Equal(other *Coordinate) bool {
	if other == nil {
		return false
	}
	return c.Latitude == other.Latitude && c.Longtitude == other.Longtitude
}
