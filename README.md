go-osgb
============

[![GoDoc](https://godoc.org/github.com/fofanov/go-osgb?status.svg)](https://godoc.org/github.com/fofanov/go-osgb)

A Go library for performing accurate conversions between OS National Grid and GPS coordinates. (More technically, coordinate transformations between the OSGB36/ODN and ETRS89 geodetic datums)

Why would I need this?
------------
Geodetic transformations that assume perfect geometry lose precision on imperfect geoids (e.g. The Helmert transformation for ETRS89 to OSGB36/ODN has errors up to 5m). This library implements the definitive transformation developed by Ordnance Survey for converting between National Grid and GPS coordinates, resulting in very accurate conversions (average < 0.1m error).

Features
------------
  - Supports OSTN02/OSGM02 and OSTN15/OSGM15 transformations
  - All transformation parameters are included in the library. No need to load additional files!
  - Fully tested against the conversion samples provided in the Ordnance Survey developer resources

Installation
------------
    go get github.com/fofanov/go-osgb

Usage
------------
Converting from National Grid to GPS
```go
    import (
        "log"

        osgb "github.com/fofanov/go-osgb"
    )

    func main() {
        trans, err := osgb.NewOSTN15Transformer()
        if err != nil {
            log.Fatal(err)
        }
        easting := 400001.4
        northing := 305001.4
        height := 5.0
        nationalGridCoord := osgb.NewOSGB36Coord(easting, northing, height)
        gpsCoord, err := trans.FromNationalGrid(nationalGridCoord)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%#v\n", gpsCoord)
        // &osgb.ETRS89Coordinate{Lon:-2.001408181413446, Lat:52.64276984554203, Height:55.34172940576156}
    }
```

Converting from GPS to National Grid
```go
    import (
        "log"

        osgb "github.com/fofanov/go-osgb"
    )

    func main() {
        trans, err := osgb.NewOSTN15Transformer()
        if err != nil {
            log.Fatal(err)
        }

        lon := -0.1262
        lat := 51.5080
        height := 10.5
        gpsCoord := osgb.NewETRS89Coord(lon, lat height)
        nationalGridCoord, err := trans.ToNationalGrid(gpsCoord)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%#v\n", nationalGridCoord)
        // &osgb.OSGB36Coordinate{Easting:530136.3274244666, Northing:180449.45428526515, Height:-35.04663654446814}
    }
```

Coordinate Units
------------
National Grid eastings, northings and ODN height are all in metres.
GPS longitude and latitude are in decimal degrees. Height is in metres.

Transformation Limits
------------
The transformation is only accurately defined for onshore positions of British Islands.

When using OSTNO2, transformations attempted for positions greater than 10km offshore will return an `ErrPointOutsidePolygon` error.

OSTN15 will not return an error for offshore transformations, but precision is severely degraded, so usage is not recommended. However, straying outside the extents of the 700x1250km transformation grid completely will lead to an `ErrPointOutsideTransformation` error.

I want to know more about the transformation
------------
The full details can be found in the [developers section](https://www.ordnancesurvey.co.uk/business-and-government/help-and-support/navigation-technology/os-net/formats-for-developers.html) of the Ordnance Survey website.


Roadmap
------------
-  [ ] Expose geoid regions for ODN heights
-  [ ] Support transformations on the Irish mainland

License
------------
This library is released under the BSD license, as are the transformation models provided by Ordnance Survey. Full details can be found in the LICENSE file.




