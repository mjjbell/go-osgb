package osgb

import (
	"errors"
	"math"
)

var (
	// ErrPointOutsidePolygon indicates the position is too far offshore to be transformed.
	ErrPointOutsidePolygon = errors.New("point outside polygon")
	// ErrPointOutsideTransformation indicates the position is completely outside the grid transformation extent
	ErrPointOutsideTransformation = errors.New("point outside transformation limits")
)

const (
	nEastIndices            = 701
	translationVectorFile02 = "data/OSTN02_OSGM02_GB.txt"
	translationVectorFile15 = "data/OSTN15_OSGM15_GB.txt"
)

type geoidRegion uint8

const (
	Region_OUTSIDE_BOUNDARY       geoidRegion = 0  // 02
	Region_UK_MAINLAND            geoidRegion = 1  // 02,15
	Region_SCILLY_ISLES           geoidRegion = 2  // 02,15
	Region_ISLE_OF_MAN            geoidRegion = 3  // 02,15
	Region_OUTER_HEBRIDES         geoidRegion = 4  // 02,15
	Region_ST_KILDA               geoidRegion = 5  // 02
	Region_SHETLAND_ISLES         geoidRegion = 6  // 02,15
	Region_ORKNEY_ISLES           geoidRegion = 7  // 02,15
	Region_FAIR_ISLE              geoidRegion = 8  // 02
	Region_FLANNAN_ISLES          geoidRegion = 9  // 02
	Region_NORTH_RONA             geoidRegion = 10 // 02
	Region_SULE_SKERRY            geoidRegion = 11 // 02
	Region_FOULA                  geoidRegion = 12 // 02
	Region_REPUBLIC_OF_IRELAND    geoidRegion = 13 // 02
	Region_NORTHERN_IRELAND       geoidRegion = 14 // 02
	Region_OFFSHORE               geoidRegion = 15 // 15
	Region_OUTSIDE_TRANSFORMATION geoidRegion = 16 // 15
)

// CoordinateTransformer is used to convert between OSGB36/ODN and ETRS89 geodetic datums
type CoordinateTransformer interface {
	// ToNationalGrid coverts a coordinate position from ETRS89 to OSGB36/ODN
	ToNationalGrid(c *ETRS89Coordinate) (*OSGB36Coordinate, error)
	// FromNationalGrid coverts a coordinate position from OSGB36/ODN to ETRS89
	FromNationalGrid(c *OSGB36Coordinate) (*ETRS89Coordinate, error)
}

type transformer struct {
	records []record
}

func (tr *transformer) ToNationalGrid(c *ETRS89Coordinate) (*OSGB36Coordinate, error) {
	latRadians := degreesToRadians(c.Lat)
	lonRadians := degreesToRadians(c.Lon)
	etrs89PlaneCoord := nationalGridProjection.toPlaneCoord(latRadians, lonRadians, grs80Ellipsoid)
	osgb36Coord, odnHeight, _, err := tr.toOSGB36(&planeCoord{
		easting:  etrs89PlaneCoord.easting,
		northing: etrs89PlaneCoord.northing,
	}, c.Height)
	if err != nil {
		return nil, err
	}
	return &OSGB36Coordinate{
		Easting:  osgb36Coord.easting,
		Northing: osgb36Coord.northing,
		Height:   odnHeight,
	}, nil
}

func (tr *transformer) FromNationalGrid(c *OSGB36Coordinate) (*ETRS89Coordinate, error) {
	etrs89Coord, etrs89Height, err := tr.fromOSGB36(&planeCoord{
		easting:  c.Easting,
		northing: c.Northing,
	}, c.Height)
	if err != nil {
		return nil, err
	}

	etrs89Lat, etrs89Lon := nationalGridProjection.fromPlaneCoord(etrs89Coord, grs80Ellipsoid)
	degreeLat := radiansToDegrees(etrs89Lat)
	degreeLon := radiansToDegrees(etrs89Lon)

	return &ETRS89Coordinate{
		Lat:    degreeLat,
		Lon:    degreeLon,
		Height: etrs89Height,
	}, nil
}

func nearestGeoidRegion(etrs89Coord *planeCoord, rs *shiftRecords) geoidRegion {

	dx := etrs89Coord.easting - float64(rs.s0.etrs89Easting)
	t := dx / 1000.0
	dy := etrs89Coord.northing - float64(rs.s0.etrs89Northing)
	u := dy / 1000.0

	if t <= 0.5 && u <= 0.5 {
		return rs.s0.geoidRegion
	} else if u <= 0.5 {
		return rs.s1.geoidRegion
	} else if t > 0.5 {
		return rs.s2.geoidRegion
	}
	return rs.s3.geoidRegion
}

func (tr *transformer) toOSGB36(etrs89Coord *planeCoord, etrs89Height float64) (*planeCoord, float64, geoidRegion, error) {

	rs, err := tr.getShiftRecords(etrs89Coord)
	if err != nil {
		return nil, 0, Region_FOULA, err
	}

	dx := etrs89Coord.easting - float64(rs.s0.etrs89Easting)
	t := dx / 1000.0
	it := 1 - t
	dy := etrs89Coord.northing - float64(rs.s0.etrs89Northing)
	u := dy / 1000.0
	iu := 1 - u
	shiftEast := it*iu*rs.s0.ostnEastShift +
		t*iu*rs.s1.ostnEastShift +
		t*u*rs.s2.ostnEastShift +
		it*u*rs.s3.ostnEastShift

	shiftNorth := it*iu*rs.s0.ostnNorthShift +
		t*iu*rs.s1.ostnNorthShift +
		t*u*rs.s2.ostnNorthShift +
		it*u*rs.s3.ostnNorthShift

	geoidHeight := it*iu*rs.s0.ostnGeoidHeight +
		t*iu*rs.s1.ostnGeoidHeight +
		t*u*rs.s2.ostnGeoidHeight +
		it*u*rs.s3.ostnGeoidHeight

	geoidRegion := nearestGeoidRegion(etrs89Coord, rs)

	return &planeCoord{
		easting:  etrs89Coord.easting + shiftEast,
		northing: etrs89Coord.northing + shiftNorth,
	}, etrs89Height - geoidHeight, geoidRegion, nil
}

func (tr *transformer) fromOSGB36(osgb36Coord *planeCoord, odnHeight float64) (*planeCoord, float64, error) {

	etrs89Coord := &planeCoord{
		easting:  osgb36Coord.easting,
		northing: osgb36Coord.northing,
	}
	etrs89Height := odnHeight

	// Iteatively find the map coordinate shift.
	for {

		rs, err := tr.getShiftRecords(etrs89Coord)
		if err != nil {
			return nil, 0, err
		}

		dx := etrs89Coord.easting - float64(rs.s0.etrs89Easting)
		t := dx / 1000.0
		it := 1 - t
		dy := etrs89Coord.northing - float64(rs.s0.etrs89Northing)
		u := dy / 1000.0
		iu := 1 - u
		shiftEast := it*iu*rs.s0.ostnEastShift +
			t*iu*rs.s1.ostnEastShift +
			t*u*rs.s2.ostnEastShift +
			it*u*rs.s3.ostnEastShift

		shiftNorth := it*iu*rs.s0.ostnNorthShift +
			t*iu*rs.s1.ostnNorthShift +
			t*u*rs.s2.ostnNorthShift +
			it*u*rs.s3.ostnNorthShift

		geoidHeight := it*iu*rs.s0.ostnGeoidHeight +
			t*iu*rs.s1.ostnGeoidHeight +
			t*u*rs.s2.ostnGeoidHeight +
			it*u*rs.s3.ostnGeoidHeight

		const epsilon = 0.0001

		newEasting := osgb36Coord.easting - shiftEast
		newNorthing := osgb36Coord.northing - shiftNorth
		newHeight := odnHeight + geoidHeight
		if math.Abs(etrs89Coord.easting-newEasting) <= epsilon &&
			math.Abs(etrs89Coord.northing-newNorthing) <= epsilon &&
			math.Abs(etrs89Height-newHeight) <= epsilon {
			break
		}
		etrs89Coord.easting = newEasting
		etrs89Coord.northing = newNorthing
		etrs89Height = newHeight
	}

	return etrs89Coord, etrs89Height, nil
}

// NewOSTN02Transformer returns a transformer that uses OSTN02/OSGM02
func NewOSTN02Transformer() (CoordinateTransformer, error) {
	records, err := readRecords(translationVectorFile02)
	if err != nil {
		return nil, err
	}
	return &transformer{
		records: records,
	}, nil
}

// NewOSTN15Transformer returns a transformer that uses OSTN15/OSGM15
func NewOSTN15Transformer() (CoordinateTransformer, error) {
	records, err := readRecords(translationVectorFile15)
	if err != nil {
		return nil, err
	}
	return &transformer{
		records: records,
	}, nil
}
