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
	v := el.semiMajorAxis / math.Sqrt(1.0-eSq*math.Sin(c.Lat)*math.Sin(c.Lat))
	x := (v + c.Height) * math.Cos(c.Lat) * math.Cos(c.Lon)
	y := (v + c.Height) * math.Cos(c.Lat) * math.Sin(c.Lon)
	z := ((1-eSq)*v + c.Height) * math.Sin(c.Lat)
	return &cartesianCoord{
		X: x,
		Y: y,
		Z: z,
	}
}

func (el *ellipsoid) cartesianToGeographic(c *cartesianCoord) *geographicCoord {
	lon := math.Atan(c.Y / c.X)
	p := math.Sqrt(c.X*c.X + c.Y*c.Y)
	eSq := el.eccentricity()
	lat := math.Atan(c.Z / (p * (1 - eSq)))
	var v float64

	// Iteratively reach new latitude value
	for {
		v = el.semiMajorAxis / math.Sqrt(1.0-eSq*math.Sin(lat)*math.Sin(lat))
		newLat := math.Atan((c.Z + eSq*v*math.Sin(lat)) / p)
		const epsilon = 0.00000000001
		if math.Abs(newLat-lat) < epsilon {
			break
		}
		lat = newLat
	}
	height := p/math.Cos(lat) - v
	return &geographicCoord{
		Lat:    lat,
		Lon:    lon,
		Height: height,
	}
}
