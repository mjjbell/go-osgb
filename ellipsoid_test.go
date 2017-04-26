package transformation

import (
	"testing"
)

func TestLatLonHeightToCartesian(t *testing.T) {

	lat, err := dmsToDecimal(52, 39, 27.2531, north)
	if err != nil {
		t.Fatal(err)
	}
	lon, err := dmsToDecimal(1, 43, 4.5177, east)
	if err != nil {
		t.Fatal(err)
	}

	latRadians := degreesToRadians(lat)
	lonRadians := degreesToRadians(lon)

	geoCoord := &geographicCoord{
		lat:    latRadians,
		lon:    lonRadians,
		height: 24.700,
	}

	cartCoord := airyEllipsoid.geographicToCartesian(geoCoord)
	expectedX := 3874938.850
	expectedY := 116218.624
	expectedZ := 5047168.207

	checkDistance(t, "x", expectedX, cartCoord.x)
	checkDistance(t, "y", expectedY, cartCoord.y)
	checkDistance(t, "z", expectedZ, cartCoord.z)
}

func TestCartesianToLatLonHeight(t *testing.T) {

	cartesianCoord := &cartesianCoord{
		x: 3874938.850,
		y: 116218.624,
		z: 5047168.207,
	}

	geoCoord := airyEllipsoid.cartesianToGeographic(cartesianCoord)

	expectedLat, err := dmsToDecimal(52, 39, 27.2531, north)
	if err != nil {
		t.Fatal(err)
	}
	expectedLon, err := dmsToDecimal(1, 43, 4.5177, east)
	if err != nil {
		t.Fatal(err)
	}
	expectedHeight := 24.700
	expectedLatRadians := degreesToRadians(expectedLat)
	expectedLonRadians := degreesToRadians(expectedLon)

	checkAngle(t, "latitude", expectedLatRadians, geoCoord.lat)
	checkAngle(t, "longitude", expectedLonRadians, geoCoord.lon)
	checkDistance(t, "height", expectedHeight, geoCoord.height)
}
