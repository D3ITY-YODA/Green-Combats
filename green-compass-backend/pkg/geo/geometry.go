package geo

import (
	"fmt"
	"math"
)

type Point struct {
	Lat float64
	Lon float64
}

func NewPoint(lat, lon float64) (Point, error) {
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return Point{}, fmt.Errorf("latitude %v is not a finite number", lat)
	}
	if math.IsNaN(lon) || math.IsInf(lon, 0) {
		return Point{}, fmt.Errorf("longitude %v is not a finite number", lon)
	}
	if lat < -90 || lat > 90 {
		return Point{}, fmt.Errorf("latitude %v out of range [-90, 90]", lat)
	}
	if lon < -180 || lon > 180 {
		return Point{}, fmt.Errorf("longitude %v out of range [-180, 180]", lon)
	}
	return Point{Lat: lat, Lon: lon}, nil
}

const earthRadiusMeters = 6371008.8

func MetersBetween(a, b Point) float64 {
	lat1 := radians(a.Lat)
	lat2 := radians(b.Lat)
	dLat := radians(b.Lat - a.Lat)
	dLon := radians(b.Lon - a.Lon)

	sinDLat := math.Sin(dLat / 2)
	sinDLon := math.Sin(dLon / 2)

	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLon*sinDLon
	return 2 * earthRadiusMeters * math.Asin(math.Min(1, math.Sqrt(h)))
}

func radians(deg float64) float64 {
	return deg * math.Pi / 180
}
