package idw

import (
	"errors"
	"math"
)

// GeoPoint represents a geometric point on earth with long-lat and 4 environmental values
type GeoPoint struct {
	Latitude    float64
	Longitude   float64
	WindSpeed   float64 // km/h
	Temperature float64 // °C
	Humidity    float64 // %
	Rainfall    float64 // mm
}

// InterpolateAll estimates windSpeed, temperature, humidity, rainfall
// of target through existed points based on IDW
func InterpolateAll(lat, long float64, dataPoints []GeoPoint, power float64) (GeoPoint, error) {
	if len(dataPoints) == 0 {
		return GeoPoint{}, errors.New("no data points provided")
	}

	// Check for exact location match (distance == 0)
	for _, dp := range dataPoints {
		if dp.Latitude == lat && dp.Longitude == long {
			return dp, nil
		}
	}

	var (
		numWind, numTemp, numHumid, numRain float64
		denom                               float64
	)

	for _, dp := range dataPoints {
		dist := haversineDistance(dp.Latitude, dp.Longitude, lat, long)

		weight := 1.0 / math.Pow(dist, power)

		numWind += weight * dp.WindSpeed
		numTemp += weight * dp.Temperature
		numHumid += weight * dp.Humidity
		numRain += weight * dp.Rainfall
		denom += weight
	}

	if denom == 0 {
		return GeoPoint{}, errors.New("sum of weights is zero")
	}

	return GeoPoint{
		Latitude:    lat,
		Longitude:   long,
		WindSpeed:   numWind / denom,
		Temperature: numTemp / denom,
		Humidity:    numHumid / denom,
		Rainfall:    numRain / denom,
	}, nil
}

// haversineDistance calculates the great-circle distance
// between two points based on Haversine formula (in km)
func haversineDistance(lat1, long1, lat2, long2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := degreesToRadians(lat2 - lat1)
	dLong := degreesToRadians(long2 - long1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLong/2)*math.Sin(dLong/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
