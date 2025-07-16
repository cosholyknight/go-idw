package idw

import (
	"testing"
)

func TestInterpolateAll(t *testing.T) {
	points := []GeoPoint{
		{Latitude: 10.0, Longitude: 105.0, WindSpeed: 15, Temperature: 30, Humidity: 60, Rainfall: 2},
		{Latitude: 10.1, Longitude: 105.1, WindSpeed: 10, Temperature: 32, Humidity: 62, Rainfall: 4},
		{Latitude: 9.9, Longitude: 104.9, WindSpeed: 12, Temperature: 29, Humidity: 63, Rainfall: 3},
	}

	targetLat := 10.05
	targetLng := 105.05
	power := 2.0

	result, err := InterpolateAll(targetLat, targetLng, points, power)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check result is between range of input points
	if result.WindSpeed < 10 || result.WindSpeed > 15 {
		t.Errorf("WindSpeed out of expected range: got %.2f", result.WindSpeed)
	}
	if result.Temperature < 29 || result.Temperature > 32 {
		t.Errorf("Temperature out of expected range: got %.2f", result.Temperature)
	}
	if result.Humidity < 60 || result.Humidity > 63 {
		t.Errorf("Humidity out of expected range: got %.2f", result.Humidity)
	}
	if result.Rainfall < 2 || result.Rainfall > 4 {
		t.Errorf("Rainfall out of expected range: got %.2f", result.Rainfall)
	}
}

func TestInterpolateExactMatch(t *testing.T) {
	points := []GeoPoint{
		{Latitude: 10.0, Longitude: 105.0, WindSpeed: 15, Temperature: 30, Humidity: 60, Rainfall: 2},
	}
	targetLat := 10.0
	targetLng := 105.0

	result, err := InterpolateAll(targetLat, targetLng, points, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != points[0] {
		t.Errorf("Expected exact match, got %+v", result)
	}
}
