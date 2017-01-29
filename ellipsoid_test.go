package transformation

import (
	"testing"
)

func TestLatLonHeightToCartesian(t *testing.T) {

	lat, err := LatToRad(52, 39, 27.2531, "N")
	if err != nil {
		t.Fatal(err)
	}
	lon, err := LonToRad(1, 43, 4.5177, "E")
	if err != nil {
		t.Fatal(err)
	}

	geoCoord := &geographicCoord{
		Lat:    lat,
		Lon:    lon,
		Height: 24.700,
	}

	cartCoord := airyEllipsoid.geographicToCartesian(geoCoord)
	expectedX := 3874938.850
	expectedY := 116218.624
	expectedZ := 5047168.207

	checkDistance(t, "x", expectedX, cartCoord.X)
	checkDistance(t, "y", expectedY, cartCoord.Y)
	checkDistance(t, "z", expectedZ, cartCoord.Z)
}

func TestCartesianToLatLonHeight(t *testing.T) {

	cartesianCoord := &cartesianCoord{
		X: 3874938.850,
		Y: 116218.624,
		Z: 5047168.207,
	}

	geoCoord := airyEllipsoid.cartesianToGeographic(cartesianCoord)

	expectedLat, err := LatToRad(52, 39, 27.2531, "N")
	if err != nil {
		t.Fatal(err)
	}
	expectedLon, err := LonToRad(1, 43, 4.5177, "E")
	if err != nil {
		t.Fatal(err)
	}
	expectedHeight := 24.700

	checkAngle(t, "latitude", expectedLat, geoCoord.Lat)
	checkAngle(t, "longitude", expectedLon, geoCoord.Lon)
	checkDistance(t, "height", expectedHeight, geoCoord.Height)
}
