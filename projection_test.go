package transformation

import "testing"

func TestLatLonToEastNort(t *testing.T) {
	lat, err := LatToRad(52, 39, 27.2531, "N")
	if err != nil {
		t.Fatal(err)
	}
	lon, err := LonToRad(1, 43, 4.5177, "E")
	if err != nil {
		t.Fatal(err)
	}
	expectedEast := 651409.903
	expectedNorth := 313177.270

	coord := nationalGridProjection.toPlaneCoord(lat, lon, airyEllipsoid)

	checkDistance(t, "east", expectedEast, coord.Easting)
	checkDistance(t, "north", expectedNorth, coord.Northing)
}

func TestEastNortToLatLon(t *testing.T) {
	east := 651409.903
	north := 313177.270
	expectedLat, err := LatToRad(52, 39, 27.2531, "N")
	if err != nil {
		t.Fatal(err)
	}
	expectedLon, err := LonToRad(1, 43, 4.5177, "E")
	if err != nil {
		t.Fatal(err)
	}

	coord := &planeCoord{
		Easting:  east,
		Northing: north,
	}

	lat, lon := nationalGridProjection.fromPlaneCoord(coord, airyEllipsoid)

	checkAngle(t, "latitude", expectedLat, lat)
	checkAngle(t, "longitude", expectedLon, lon)
}
