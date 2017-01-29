package transformation

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

type ostn02TestInput struct {
	stationName string
	etrs89X     float64
	etrs89Y     float64
	etrs89Z     float64
}

const (
	test02InputFile  = "testdata/OSTN02_OSGM02Tests_In.txt"
	test02OutputFile = "testdata/OSTN02_OSGM02Tests_Out.txt"
)

func read02InputData() (map[string]ostn02TestInput, error) {
	inputData := map[string]ostn02TestInput{}

	f, err := os.Open(test02InputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn02_osgm02 test input data...")
	// Read header
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		stationName := record[0]
		etrs89X, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		etrs89Y, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		etrs89Z, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		inputData[stationName] = ostn02TestInput{
			stationName: stationName,
			etrs89X:     etrs89X,
			etrs89Y:     etrs89Y,
			etrs89Z:     etrs89Z,
		}
	}
	log.Println("Reading ostn02_osgm02 test input data completed...")

	return inputData, nil
}

type ostn02TestOutput struct {
	stationName    string
	etrs89X        float64
	etrs89Y        float64
	etrs89Z        float64
	etrs89Lat      float64
	etrs89Lon      float64
	etrs89Height   float64
	etrs89Easting  float64
	etrs89Northing float64
	osgb36Easting  float64
	osgb36Northing float64
	osgb36Lat      float64
	osgb36Lon      float64
	odnHeight      float64
	geoidModelID   geoidRegion
}

func read02OutputData() (map[string]ostn02TestOutput, error) {

	outputData := map[string]ostn02TestOutput{}

	f, err := os.Open(test02OutputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn02_osgm02 test output data...")
	// Read header
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		stationName := record[0]
		etrs89X, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		etrs89Y, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		etrs89Z, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		etrs89LatNS := record[4]
		etrs89LatDeg, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		etrs89LatMin, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, err
		}
		etrs89LatSec, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, err
		}
		etrs89Lat, err := LatToRad(etrs89LatDeg, etrs89LatMin, etrs89LatSec, etrs89LatNS)
		if err != nil {
			return nil, err
		}

		etrs89LonEW := record[8]
		etrs89LonDeg, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			return nil, err
		}
		etrs89LonMin, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			return nil, err
		}
		etrs89LonSec, err := strconv.ParseFloat(record[11], 64)
		if err != nil {
			return nil, err
		}
		etrs89Lon, err := LonToRad(etrs89LonDeg, etrs89LonMin, etrs89LonSec, etrs89LonEW)
		if err != nil {
			return nil, err
		}

		etrs89Height, err := strconv.ParseFloat(record[12], 64)
		if err != nil {
			return nil, err
		}
		etrs89Easting, err := strconv.ParseFloat(record[13], 64)
		if err != nil {
			return nil, err
		}
		etrs89Northing, err := strconv.ParseFloat(record[14], 64)
		if err != nil {
			return nil, err
		}

		// If special 'outside' stations, stop parsing here.
		if strings.HasPrefix(stationName, "Outside") {
			outputData[stationName] = ostn02TestOutput{
				stationName:    stationName,
				etrs89X:        etrs89X,
				etrs89Y:        etrs89Y,
				etrs89Z:        etrs89Z,
				etrs89Lat:      etrs89Lat,
				etrs89Lon:      etrs89Lon,
				etrs89Height:   etrs89Height,
				etrs89Easting:  etrs89Easting,
				etrs89Northing: etrs89Northing,
			}
			continue
		}

		osgb36Easting, err := strconv.ParseFloat(record[15], 64)
		if err != nil {
			return nil, err
		}
		osgb36Northing, err := strconv.ParseFloat(record[16], 64)
		if err != nil {
			return nil, err
		}

		osgb36LatNS := record[17]
		osgb36LatDeg, err := strconv.ParseFloat(record[18], 64)
		if err != nil {
			return nil, err
		}
		osgb36LatMin, err := strconv.ParseFloat(record[19], 64)
		if err != nil {
			return nil, err
		}
		osgb36LatSec, err := strconv.ParseFloat(record[20], 64)
		if err != nil {
			return nil, err
		}
		osgb36Lat, err := LatToRad(osgb36LatDeg, osgb36LatMin, osgb36LatSec, osgb36LatNS)
		if err != nil {
			return nil, err
		}

		osgb36LonEW := record[21]
		osgb36LonDeg, err := strconv.ParseFloat(record[22], 64)
		if err != nil {
			return nil, err
		}
		osgb36LonMin, err := strconv.ParseFloat(record[23], 64)
		if err != nil {
			return nil, err
		}
		osgb36LonSec, err := strconv.ParseFloat(record[24], 64)
		if err != nil {
			return nil, err
		}
		osgb36Lon, err := LonToRad(osgb36LonDeg, osgb36LonMin, osgb36LonSec, osgb36LonEW)
		if err != nil {
			return nil, err
		}

		odnHeight, err := strconv.ParseFloat(record[25], 64)
		if err != nil {
			return nil, err
		}

		geoidModelID, err := strconv.ParseUint(record[26], 10, 8)
		if err != nil {
			return nil, err
		}

		outputData[stationName] = ostn02TestOutput{
			stationName:    stationName,
			etrs89X:        etrs89X,
			etrs89Y:        etrs89Y,
			etrs89Z:        etrs89Z,
			etrs89Lat:      etrs89Lat,
			etrs89Lon:      etrs89Lon,
			etrs89Height:   etrs89Height,
			etrs89Easting:  etrs89Easting,
			etrs89Northing: etrs89Northing,
			osgb36Easting:  osgb36Easting,
			osgb36Northing: osgb36Northing,
			osgb36Lat:      osgb36Lat,
			osgb36Lon:      osgb36Lon,
			odnHeight:      odnHeight,
			geoidModelID:   geoidDatumToRegion(geoidModelID),
		}
	}
	log.Println("Reading ostn02_osgm02 test output data completed...")

	return outputData, nil
}

func Test02Data(t *testing.T) {
	inputs, err := read02InputData()
	if err != nil {
		t.Fatal(err)
	}

	outputs, err := read02OutputData()
	if err != nil {
		t.Fatal(err)
	}

	trans, err := NewOSTN02Transformer()
	if err != nil {
		t.Fatal(err)
	}

	for station, input := range inputs {
		output, ok := outputs[station]
		if !ok {
			t.Fatal("missing station in output ", station)
		}
		checkDistance(t, "etrs89 x", output.etrs89X, input.etrs89X)
		checkDistance(t, "etrs89 y", output.etrs89Y, input.etrs89Y)
		checkDistance(t, "etrs89 z", output.etrs89Z, input.etrs89Z)

		geoCoord := grs80Ellipsoid.cartesianToGeographic(&cartesianCoord{
			X: input.etrs89X,
			Y: input.etrs89Y,
			Z: input.etrs89Z,
		})

		checkAngle(t, "etrs89 latitude", output.etrs89Lat, geoCoord.Lat)
		checkAngle(t, "etrs89 longitude", output.etrs89Lon, geoCoord.Lon)
		checkDistance(t, "etrs89 height", output.etrs89Height, geoCoord.Height)

		etrs89Coord := nationalGridProjection.toPlaneCoord(geoCoord.Lat, geoCoord.Lon, grs80Ellipsoid)

		checkDistance(t, "etrs89 east", output.etrs89Easting, etrs89Coord.Easting)
		checkDistance(t, "etrs89 north", output.etrs89Northing, etrs89Coord.Northing)

		osgb36Coord, err := trans.ToNationalGrid(&ETRS89Coordinate{
			Lat:    geoCoord.Lat,
			Lon:    geoCoord.Lon,
			Height: geoCoord.Height,
		})

		if strings.HasPrefix(station, "Outside") {
			if err != ErrPointOutsidePolygon {
				if err != nil {
					t.Errorf("Unexpected error for station %s when performing etrs89 to osgb36/odn: %s", station, err)
				} else {

					t.Errorf("Didn't receive out of polygon error for station %s when performing etrs89 to osgb36/odn", station)
				}
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for station %s: %s", station, err)
			continue
		}

		checkDistance(t, "osgb36 east", output.osgb36Easting, osgb36Coord.Easting)
		checkDistance(t, "osgb36 north", output.osgb36Northing, osgb36Coord.Northing)
		checkDistance(t, "orthometric height", output.odnHeight, osgb36Coord.Height)
		//		checkRegion(t, "geoid region", output.geoidModelID, geoidRegion)

		osgb36Lat, osgb36Lon := nationalGridProjection.fromPlaneCoord(&planeCoord{
			Easting:  osgb36Coord.Easting,
			Northing: osgb36Coord.Northing,
		}, airyEllipsoid)

		checkAngle(t, "osgb36 lat", output.osgb36Lat, osgb36Lat)
		checkAngle(t, "osgb36 lon", output.osgb36Lon, osgb36Lon)
	}
}

func Test02ETRS89ToOSGB36(t *testing.T) {
	etrs89Lat, err := LatToRad(52, 39, 28.8282, "N")
	if err != nil {
		t.Fatal(err)
	}
	etrs89Lon, err := LonToRad(1, 42, 57.8663, "E")
	if err != nil {
		t.Fatal(err)
	}
	etrs89Height := 108.05

	trans, err := NewOSTN02Transformer()
	if err != nil {
		t.Fatal(err)
	}

	expectedEasting := 651409.792
	expectedNorthing := 313177.448
	expectedHeight := 63.806
	//expectedRegion := Region_UK_MAINLAND
	osgb36Coord, err := trans.ToNationalGrid(&ETRS89Coordinate{
		lat:    etrs89Lat,
		lon:    etrs89Lon,
		height: etrs89Height,
	})
	if err != nil {
		t.Errorf("failed to convert etrs89 to osgb36/odn: %s", err)
		t.FailNow()
	}

	checkDistance(t, "national grid east", expectedEasting, osgb36Coord.easting)
	checkDistance(t, "national grid north", expectedNorthing, osgb36Coord.northing)
	checkDistance(t, "national grid height", expectedHeight, osgb36Coord.height)
	//checkRegion(t, "geoid datum ID", expectedRegion, region)
}

func Test02OSGB36ToETRS89(t *testing.T) {
	osgb36Easting := 651409.792
	osgb36Northing := 313177.448
	orthometricHeight := 63.806

	trans, err := NewOSTN02Transformer()
	if err != nil {
		t.Fatal(err)
	}

	expectedETRS89Lat, err := LatToRad(52, 39, 28.8282, "N")
	if err != nil {
		t.Fatal(err)
	}
	expectedETRS89Lon, err := LonToRad(1, 42, 57.8663, "E")
	if err != nil {
		t.Fatal(err)
	}
	expectedETRS89Height := 108.05

	etrs89Coord, err := trans.FromNationalGrid(&OSGB36Coordinate{
		easting:  osgb36Easting,
		northing: osgb36Northing,
		height:   orthometricHeight,
	})

	checkDistance(t, "etrs89 lat", expectedETRS89Lat, etrs89Coord.lat)
	checkDistance(t, "etrs89 lon", expectedETRS89Lon, etrs89Coord.lon)
	checkDistance(t, "etrs89 height", expectedETRS89Height, etrs89Coord.height)
}

func Test02ETRS89ToOSGB36_OutsideTransformationRange(t *testing.T) {
	etrs89Lat := DegreesToRad(0.0)
	etrs89Lon := DegreesToRad(0.0)
	etrs89Height := 0.0

	trans, err := NewOSTN02Transformer()
	if err != nil {
		t.Fatal(err)
	}

	_, err = trans.ToNationalGrid(&ETRS89Coordinate{
		lat:    etrs89Lat,
		lon:    etrs89Lon,
		height: etrs89Height,
	})
	if err == nil {
		t.Errorf("expected error when converting etrs89 coords outside ostn02 transformation range")
	}
}

func Test02OSGB36ToETRS89_OutsideTransformationRange(t *testing.T) {
	osgb36Easting := -651409.792
	osgb36Northing := -313177.448
	odnHeight := 0.0

	trans, err := NewOSTN02Transformer()
	if err != nil {
		t.Fatal(err)
	}
	_, err = trans.FromNationalGrid(&OSGB36Coordinate{
		easting:  osgb36Easting,
		northing: osgb36Northing,
		height:   odnHeight,
	})
	if err == nil {
		t.Errorf("expected error when converting osgb36 coords outside ostn02 transformation range")
	}
}
