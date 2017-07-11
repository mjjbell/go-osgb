package osgb

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"math"
	"strconv"

	"github.com/fofanov/go-osgb/internal/data"
)

type record struct {
	recordNo        uint32
	etrs89Easting   uint32
	etrs89Northing  uint32
	ostnEastShift   float64
	ostnNorthShift  float64
	ostnGeoidHeight float64
	geoidRegion     geoidRegion
}

func readRecords(translationVectorFile string) ([]record, error) {

	res := []record{}

	data, err := data.Asset(translationVectorFile)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bytes.NewReader(data))
	// Read header
	if _, err := r.Read(); err != nil {
		return nil, err
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		recordNo, err := strconv.ParseUint(rec[0], 10, 32)
		if err != nil {
			return nil, err
		}
		etrs89Easting, err := strconv.ParseUint(rec[1], 10, 32)
		if err != nil {
			return nil, err
		}
		etrs89Northing, err := strconv.ParseUint(rec[2], 10, 32)
		if err != nil {
			return nil, err
		}
		ostnEastShift, err := strconv.ParseFloat(rec[3], 64)
		if err != nil {
			return nil, err
		}
		ostnNorthShift, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			return nil, err
		}
		ostnGeoidHeight, err := strconv.ParseFloat(rec[5], 64)
		if err != nil {
			return nil, err
		}
		geoidDatum, err := strconv.ParseUint(rec[6], 10, 8)
		if err != nil {
			return nil, err
		}

		res = append(res,
			record{
				recordNo:        uint32(recordNo),
				etrs89Easting:   uint32(etrs89Easting),
				etrs89Northing:  uint32(etrs89Northing),
				ostnEastShift:   ostnEastShift,
				ostnNorthShift:  ostnNorthShift,
				ostnGeoidHeight: ostnGeoidHeight,
				geoidRegion:     geoidDatumToRegion(geoidDatum),
			})
	}
	return res, nil
}

func (tr *transformer) getShiftRecord(eastIndex, northIndex uint32) (*record, error) {
	recordIndex := eastIndex + northIndex*nEastIndices
	if recordIndex <= 0 || recordIndex >= uint32(len(tr.records)) {
		return nil, ErrPointOutsidePolygon
	}
	rec := &tr.records[recordIndex]
	if rec.geoidRegion == Region_OUTSIDE_BOUNDARY {
		return nil, ErrPointOutsidePolygon
	}
	if rec.geoidRegion == Region_OUTSIDE_TRANSFORMATION {
		return nil, ErrPointOutsideTransformation
	}
	return rec, nil
}

type shiftRecords struct {
	s2, s3, s0, s1 *record
}

func (tr *transformer) getShiftRecords(etrs89Coord *planeCoord) (*shiftRecords, error) {
	eastIndex := eastingIndex(etrs89Coord.easting)
	northIndex := northingIndex(etrs89Coord.northing)

	bl, err := tr.getShiftRecord(eastIndex, northIndex)
	if err != nil {
		return nil, err
	}
	br, err := tr.getShiftRecord(eastIndex+1, northIndex)
	if err != nil {
		return nil, err
	}
	rt, err := tr.getShiftRecord(eastIndex+1, northIndex+1)
	if err != nil {
		return nil, err
	}
	tl, err := tr.getShiftRecord(eastIndex, northIndex+1)
	if err != nil {
		return nil, err
	}

	return &shiftRecords{
		s0: bl,
		s1: br,
		s2: rt,
		s3: tl,
	}, nil
}

func eastingIndex(easting float64) uint32 {
	return uint32(math.Floor(easting / 1000.0))
}

func northingIndex(northing float64) uint32 {
	return uint32(math.Floor(northing / 1000.0))
}

func geoidDatumToRegion(id uint64) geoidRegion {
	if id < 0 || id > 16 {
		// This is unexpected and we dont't know how to recover from this.
		log.Fatalf("unexpected geoid datum ID %d", id)
	}
	return geoidRegion(id)
}
