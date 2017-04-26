package transformation

import (
	"math"
)

type ellipsoid struct {
	semiMajorAxis float64
	semiMinorAxis float64
}

var (
	airyEllipsoid = &ellipsoid{
		semiMajorAxis: 6377563.396,
		semiMinorAxis: 6356256.909,
	}
	grs80Ellipsoid = &ellipsoid{
		semiMajorAxis: 6378137.000,
		semiMinorAxis: 6356752.3141,
	}
)

func (el *ellipsoid) eccentricity() float64 {
	aSq := el.semiMajorAxis * el.semiMajorAxis
	bSq := el.semiMinorAxis * el.semiMinorAxis
	return (aSq - bSq) / aSq
}

func (el *ellipsoid) geographicToCartesian(c *geographicCoord) *cartesianCoord {
	eSq := el.eccentricity()
	v := el.semiMajorAxis / math.Sqrt(1.0-eSq*math.Sin(c.lat)*math.Sin(c.lat))
	x := (v + c.height) * math.Cos(c.lat) * math.Cos(c.lon)
	y := (v + c.height) * math.Cos(c.lat) * math.Sin(c.lon)
	z := ((1-eSq)*v + c.height) * math.Sin(c.lat)
	return &cartesianCoord{
		x: x,
		y: y,
		z: z,
	}
}

func (el *ellipsoid) cartesianToGeographic(c *cartesianCoord) *geographicCoord {
	lon := math.Atan(c.y / c.x)
	p := math.Sqrt(c.x*c.x + c.y*c.y)
	eSq := el.eccentricity()
	lat := math.Atan(c.z / (p * (1 - eSq)))
	var v float64

	// Iteratively reach new latitude value
	for {
		v = el.semiMajorAxis / math.Sqrt(1.0-eSq*math.Sin(lat)*math.Sin(lat))
		newLat := math.Atan((c.z + eSq*v*math.Sin(lat)) / p)
		const epsilon = 0.00000000001
		if math.Abs(newLat-lat) < epsilon {
			break
		}
		lat = newLat
	}
	height := p/math.Cos(lat) - v
	return &geographicCoord{
		lat:    lat,
		lon:    lon,
		height: height,
	}
}
