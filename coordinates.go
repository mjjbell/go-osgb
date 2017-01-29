package transformation

import (
	"fmt"
	"math"
)

type ETRS89Coordinate struct {
	Lat, Lon, Height float64
}

func NewETRS89Cartesian(x, y, z float64) *ETRS89Coordinate {
	return &ETRS89Coordinate{}
}

func NewETRS89DecimalDegrees(lat, lon, height float64) *ETRS89Coordinate {
	return &ETRS89Coordinate{
		Lat:    lat,
		Lon:    lon,
		Height: height,
	}
}

func NewETRS89DMS(degrees, minutes, seconds, height float64) *ETRS89Coordinate {
	return &ETRS89Coordinate{}
}

type OSGB36Coordinate struct {
	Easting, Northing, Height float64
}

func NewOSGB36(easting, northing, height float64) *OSGB36Coordinate {
	return &OSGB36Coordinate{
		Easting:  easting,
		Northing: northing,
		Height:   height,
	}
}

type geographicCoord struct {
	Lat, Lon, Height float64
}

type cartesianCoord struct {
	X, Y, Z float64
}

type planeCoord struct {
	Easting, Northing float64
}

func LonToRad(degree, minute, seconds float64, direction string) (float64, error) {
	if degree < 0 || degree > 180 {
		return 0, fmt.Errorf("invalid degree %f", degree)
	}
	if minute < 0 || minute > 60 {
		return 0, fmt.Errorf("invalid minute %f", minute)
	}
	if seconds < 0 || seconds > 60 {
		return 0, fmt.Errorf("invalid seconds %f", seconds)
	}
	if direction != "E" && direction != "W" {
		return 0, fmt.Errorf("invalid direction %s", direction)
	}

	rad := (degree + minute/60 + seconds/3600) * (math.Pi / 180)

	if direction == "W" {
		return rad * -1.0, nil
	}
	return rad, nil
}

func LatToRad(degree, minute, seconds float64, direction string) (float64, error) {
	if degree < 0 || degree > 90 {
		return 0, fmt.Errorf("invalid degree %f", degree)
	}
	if minute < 0 || minute > 60 {
		return 0, fmt.Errorf("invalid minute %f", minute)
	}
	if seconds < 0 || seconds > 60 {
		return 0, fmt.Errorf("invalid seconds %f", seconds)
	}
	if direction != "N" && direction != "S" {
		return 0, fmt.Errorf("invalid direction %s", direction)
	}

	rad := (degree + minute/60 + seconds/3600) * (math.Pi / 180)

	if direction == "S" {
		return rad * -1.0, nil
	}
	return rad, nil
}

func RadToDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func DegreesToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}
