package geo

import (
	"math"
)

const (
	// EarthRadiusKm is the radius of the Earth in kilometers
	EarthRadiusKm = 6371.0
)

// HaversineDistance calculates the great-circle distance between two points
// on the Earth's surface using the Haversine formula.
// Returns distance in kilometers.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1Rad := degToRad(lat1)
	lon1Rad := degToRad(lon1)
	lat2Rad := degToRad(lat2)
	lon2Rad := degToRad(lon2)

	// Calculate differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	distance := EarthRadiusKm * c

	return distance
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// radToDeg converts radians to degrees
func radToDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}