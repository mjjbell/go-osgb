package osgb

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

type ostn15OSGBToETRSTestInput struct {
	pointID           string
	osgbEasting       float64
	osgbNorthing      float64
	orthometricHeight float64
}

type ostn15ETRSToOSGBTestInput struct {
	pointID      string
	etrs89Lat    float64
	etrs89Lon    float64
	etrs89Height float64
}

const (
	test15ETRSToOSGBInputFile  = "testdata/OSTN15_OSGM15_TestInput_ETRStoOSGB.txt"
	test15OSGBToETRSInputFile  = "testdata/OSTN15_OSGM15_TestInput_OSGBtoETRS.txt"
	test15ETRSToOSGBOutputFile = "testdata/OSTN15_OSGM15_TestOutput_ETRStoOSGB.txt"
	test15OSGBToETRSOutputFile = "testdata/OSTN15_OSGM15_TestOutput_OSGBtoETRS.txt"
)

func read15ETRSToOSGBInputData() (map[string]ostn15ETRSToOSGBTestInput, error) {
	inputData := map[string]ostn15ETRSToOSGBTestInput{}

	f, err := os.Open(test15ETRSToOSGBInputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn15_osgm15 test input data...")
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

		pointID := record[0]
		etrs89Lat, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		etrs89Lon, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		etrs89Height, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		inputData[pointID] = ostn15ETRSToOSGBTestInput{
			pointID:      pointID,
			etrs89Lat:    etrs89Lat,
			etrs89Lon:    etrs89Lon,
			etrs89Height: etrs89Height,
		}
	}
	log.Println("Reading ostn15_osgm15 test input data completed...")

	return inputData, nil
}

func read15OSGBToETRSInputData() (map[string]ostn15OSGBToETRSTestInput, error) {
	inputData := map[string]ostn15OSGBToETRSTestInput{}

	f, err := os.Open(test15OSGBToETRSInputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn15_osgm15 test input data...")
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

		pointID := record[0]
		osgb36Easting, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		osgb36Northing, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		orthometricHeight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		inputData[pointID] = ostn15OSGBToETRSTestInput{
			pointID:           pointID,
			osgbEasting:       osgb36Easting,
			osgbNorthing:      osgb36Northing,
			orthometricHeight: orthometricHeight,
		}
	}
	log.Println("Reading ostn15_osgm15 test input data completed...")

	return inputData, nil
}

type ostn15ETRSToOSGBTestOutput struct {
	pointID        string
	osgb36Easting  float64
	osgb36Northing float64
	odnHeight      float64
	geoidModelID   geoidRegion
}

type ostn15OSGBToETRSTestOutput struct {
	pointID      string
	etrs89Lat    float64
	etrs89Lon    float64
	etrs89Height float64
}

func read15ETRSToOSGBOutputData() (map[string]ostn15ETRSToOSGBTestOutput, error) {

	outputData := map[string]ostn15ETRSToOSGBTestOutput{}

	f, err := os.Open(test15ETRSToOSGBOutputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn15_osgm15 test output data...")
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

		pointID := record[0]
		osgb36Easting, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		osgb36Northing, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		odnHeight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		geoidModelID, err := strconv.ParseUint(record[4], 10, 8)
		if err != nil {
			return nil, err
		}

		outputData[pointID] = ostn15ETRSToOSGBTestOutput{
			pointID:        pointID,
			osgb36Easting:  osgb36Easting,
			osgb36Northing: osgb36Northing,
			odnHeight:      odnHeight,
			geoidModelID:   geoidDatumToRegion(geoidModelID),
		}
	}
	log.Println("Reading ostn15_osgm15 test output data completed...")

	return outputData, nil
}

func read15OSGBToETRSOutputData() (map[string]ostn15OSGBToETRSTestOutput, error) {

	outputData := map[string]ostn15OSGBToETRSTestOutput{}

	f, err := os.Open(test15OSGBToETRSOutputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	log.Println("Reading ostn15_osgm15 test output data...")
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
		if iteration := record[1]; iteration != "RESULT" {
			continue
		}

		pointID := record[0]
		etrs89Lat, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		etrs89Lon, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		etrs89Height, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		outputData[pointID] = ostn15OSGBToETRSTestOutput{
			pointID:      pointID,
			etrs89Lat:    etrs89Lat,
			etrs89Lon:    etrs89Lon,
			etrs89Height: etrs89Height,
		}
	}
	log.Println("Reading ostn15_osgm15 test output data completed...")

	return outputData, nil
}

func Test15ETRS89ToOSGB36Data(t *testing.T) {
	inputs, err := read15ETRSToOSGBInputData()
	if err != nil {
		t.Fatal(err)
	}

	outputs, err := read15ETRSToOSGBOutputData()
	if err != nil {
		t.Fatal(err)
	}

	trans, err := NewOSTN15Transformer()
	if err != nil {
		t.Fatal(err)
	}

	for pointID, input := range inputs {
		output, ok := outputs[pointID]
		if !ok {
			t.Fatal("missing point ID in output ", pointID)
		}

		osgb36Coord, err := trans.ToNationalGrid(&ETRS89Coordinate{
			Lat:    input.etrs89Lat,
			Lon:    input.etrs89Lon,
			Height: input.etrs89Height,
		})

		if err != nil {
			t.Errorf("Unexpected error for point ID %s: %s", pointID, err)
			continue
		}

		checkDistance(t, "osgb36 east", output.osgb36Easting, osgb36Coord.Easting)
		checkDistance(t, "osgb36 north", output.osgb36Northing, osgb36Coord.Northing)
		checkDistance(t, "orthometric region", output.odnHeight, osgb36Coord.Height)
	}
}

func Test15OSGB36ToETRS89Data(t *testing.T) {
	inputs, err := read15OSGBToETRSInputData()
	if err != nil {
		t.Fatal(err)
	}

	outputs, err := read15OSGBToETRSOutputData()
	if err != nil {
		t.Fatal(err)
	}

	trans, err := NewOSTN15Transformer()
	if err != nil {
		t.Fatal(err)
	}

	for pointID, input := range inputs {
		output, ok := outputs[pointID]
		if !ok {
			t.Fatal("missing point ID in output ", pointID)
		}

		etrs89Coord, err := trans.FromNationalGrid(&OSGB36Coordinate{
			Easting:  input.osgbEasting,
			Northing: input.osgbNorthing,
			Height:   input.orthometricHeight,
		})

		if err != nil {
			t.Errorf("Unexpected error for point ID %s: %s", pointID, err)
			continue
		}

		checkAngle(t, "etrs89 lat", output.etrs89Lat, etrs89Coord.Lat)
		checkAngle(t, "etrs89 lon", output.etrs89Lon, etrs89Coord.Lon)
		checkDistance(t, "etrs89 height", output.etrs89Height, etrs89Coord.Height)
	}
}

func Test15ETRS89ToOSGB36(t *testing.T) {
	etrs89Lat, err := dmsToDecimal(52, 39, 28.8282, north)
	if err != nil {
		t.Fatal(err)
	}
	etrs89Lon, err := dmsToDecimal(1, 42, 57.8663, east)
	if err != nil {
		t.Fatal(err)
	}
	etrs89Height := 108.05

	trans, err := NewOSTN15Transformer()
	if err != nil {
		t.Fatal(err)
	}

	expectedEasting := 651409.804
	expectedNorthing := 313177.450
	expectedHeight := 63.822
	//expectedRegion := Region_UK_MAINLAND
	osgb36Coord, err := trans.ToNationalGrid(&ETRS89Coordinate{
		Lat:    etrs89Lat,
		Lon:    etrs89Lon,
		Height: etrs89Height,
	})
	if err != nil {
		t.Errorf("failed to convert etrs89 to osgb36/odn: %s", err)
		t.FailNow()
	}

	checkDistance(t, "national grid east", expectedEasting, osgb36Coord.Easting)
	checkDistance(t, "national grid north", expectedNorthing, osgb36Coord.Northing)
	checkDistance(t, "national grid height", expectedHeight, osgb36Coord.Height)
}

func Test15OSGB36ToETRS89(t *testing.T) {
	osgbEasting := 651409.804
	osgbNorthing := 313177.450
	orthometricHeight := 63.822

	trans, err := NewOSTN15Transformer()
	if err != nil {
		t.Fatal(err)
	}

	expectedETRS89Lat, err := dmsToDecimal(52, 39, 28.8282, north)
	if err != nil {
		t.Fatal(err)
	}
	expectedETRS89Lon, err := dmsToDecimal(1, 42, 57.8663, east)
	if err != nil {
		t.Fatal(err)
	}
	expectedETRS89Height := 108.05

	etrs89Coord, err := trans.FromNationalGrid(&OSGB36Coordinate{
		Easting:  osgbEasting,
		Northing: osgbNorthing,
		Height:   orthometricHeight,
	})

	checkDistance(t, "etrs89 lat", expectedETRS89Lat, etrs89Coord.Lat)
	checkDistance(t, "etrs89 lon", expectedETRS89Lon, etrs89Coord.Lon)
	checkDistance(t, "etrs89 height", expectedETRS89Height, etrs89Coord.Height)
}

func Test15ETRS89ToOSGB36_OutsideTransformationRange(t *testing.T) {
	testData := []struct {
		lat float64
		lon float64
	}{
		// Coordinate is off the national grid.
		{
			lat: 0.0,
			lon: 0.0,
		},
		// Dublin. On the grid but out of transformation range.
		{
			lat: 53.3498,
			lon: -6.2603,
		},
	}

	for _, c := range testData {
		etrs89Lat := c.lat
		etrs89Lon := c.lon
		etrs89Height := 0.0

		trans, err := NewOSTN15Transformer()
		if err != nil {
			t.Fatal(err)
		}

		_, err = trans.ToNationalGrid(&ETRS89Coordinate{
			Lat:    etrs89Lat,
			Lon:    etrs89Lon,
			Height: etrs89Height,
		})
		if err == nil {
			t.Errorf("expected error when converting etrs89 coords outside ostn15 transformation range")
		}
	}
}

func Test15OSGB36ToETRS89_OutsideTransformationRange(t *testing.T) {
	osgb36Easting := -651409.792
	osgb36Northing := -313177.448
	odnHeight := 0.0

	osgb36Coord := &OSGB36Coordinate{
		Easting:  osgb36Easting,
		Northing: osgb36Northing,
		Height:   odnHeight,
	}

	trans, err := NewOSTN15Transformer()
	if err != nil {
		t.Fatal(err)
	}
	_, err = trans.FromNationalGrid(osgb36Coord)
	if err == nil {
		t.Errorf("expected error when converting osgb36 coords outside ostn15 transformation range")
	}
}
