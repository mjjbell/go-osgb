package transformation

import "testing"

func TestLatLonToEastNort(t *testing.T) {
	lat, err := dmsToDecimal(52, 39, 27.2531, north)
	if err != nil {
		t.Fatal(err)
	}
	latRadians := degreesToRadians(lat)
	lon, err := dmsToDecimal(1, 43, 4.5177, east)
	if err != nil {
		t.Fatal(err)
	}
	lonRadians := degreesToRadians(lon)
	expectedEast := 651409.903
	expectedNorth := 313177.270

	coord := nationalGridProjection.toPlaneCoord(latRadians, lonRadians, airyEllipsoid)

	checkDistance(t, "east", expectedEast, coord.easting)
	checkDistance(t, "north", expectedNorth, coord.northing)
}

func TestEastNortToLatLon(t *testing.T) {
	easting := 651409.903
	northing := 313177.270
	expectedLat, err := dmsToDecimal(52, 39, 27.2531, north)
	if err != nil {
		t.Fatal(err)
	}
	expectedLatRadians := degreesToRadians(expectedLat)
	expectedLon, err := dmsToDecimal(1, 43, 4.5177, east)
	if err != nil {
		t.Fatal(err)
	}
	expectedLonRadians := degreesToRadians(expectedLon)

	coord := &planeCoord{
		easting:  easting,
		northing: northing,
	}

	lat, lon := nationalGridProjection.fromPlaneCoord(coord, airyEllipsoid)

	checkAngle(t, "latitude", expectedLatRadians, lat)
	checkAngle(t, "longitude", expectedLonRadians, lon)
}
