package transformation

import (
	"math"
)

type projection struct {
	scaleFactor        float64
	geodeticTrueOrigin geographicCoord
	mapTrueOrigin      planeCoord
}

var (
	nationalGridProjection = &projection{
		scaleFactor: 0.9996012717,
		geodeticTrueOrigin: geographicCoord{
			lat: degreesToRadians(49.0),
			lon: degreesToRadians(-2.0),
		},
		mapTrueOrigin: planeCoord{
			easting:  400000,
			northing: -100000,
		},
	}
)

func (proj *projection) toPlaneCoord(φ, λ float64, el *ellipsoid) *planeCoord {
	// a - semi0major axis (metres)
	a := el.semiMajorAxis
	// b - semi-minor axis (metres)
	b := el.semiMinorAxis
	// e^2 - ellipsoid squared eccentricity constant.
	e2 := el.eccentricity()
	// N0 – northing of true origin;
	n0 := proj.mapTrueOrigin.northing
	// E0 – easting of true origin;
	e0 := proj.mapTrueOrigin.easting
	// F0 – scale factor on central meridian;
	f0 := proj.scaleFactor
	// φ0 – latitude of true origin; and
	φ0 := proj.geodeticTrueOrigin.lat
	// λ0 – longitude of true origin and central meridian.
	λ0 := proj.geodeticTrueOrigin.lon

	// (B2) n = a−b/a+b
	n := (a - b) / (a + b)

	//tan^2(φ)
	t2φ := math.Tan(φ) * math.Tan(φ)
	//tan^4(φ)
	t4φ := t2φ * t2φ
	// sinφ
	sφ := math.Sin(φ)
	// sin^2(φ)
	s2φ := sφ * sφ
	// cosφ
	cφ := math.Cos(φ)
	// cos^3φ
	c3φ := cφ * cφ * cφ
	// cos^5φ
	c5φ := c3φ * cφ * cφ
	// n^2
	n2 := n * n
	// n^3
	n3 := n2 * n

	// (B3) ν = aF0 (1 − e^2 sin^2(φ)) ^ −0.5
	ν := a * f0 * math.Pow(1-e2*s2φ, -0.5)
	// (B4) ρ=aF0(1−e^2)(1−e^2 sin^2(φ)) ^−1.5
	ρ := a * f0 * (1 - e2) * math.Pow(1-e2*s2φ, -1.5)
	// (B5) η2 = ν/ρ − 1
	η2 := ν/ρ - 1
	// φ - φ0
	dφ0 := φ - φ0
	// φ + φ0
	aφ0 := φ + φ0

	//Ma = (1+n+(5/4)n^2 +(5/4)n^3)(φ−φ0)
	ma := (1 + n + 5.0*n2/4.0 + 5.0*n3/4.0) * dφ0
	//Mb = (3n+3n^2 + (21/8)n^3)sin(φ−φ0)cos(φ+φ0)
	mb := (3*n + 3*n2 + (21.0/8.0)*n3) * math.Sin(dφ0) * math.Cos(aφ0)
	//Mc = ((15/8)n^2 + (15/8)n^3)sin(2(φ−φ0))cos(2(φ+φ0))
	mc := ((15.0/8.0)*n2 + (15.0/8.0)*n3) * math.Sin(2*dφ0) * math.Cos(2*(aφ0))
	//Md = (35/24)n^3 sin(3(φ −φ0))cos(3(φ +φ0))
	md := (35.0 / 24.0) * n3 * math.Sin(3*dφ0) * math.Cos(3*(aφ0))

	// (B6) M = bF0 (Ma - Mb + Mc - Md)
	m := b * f0 * (ma + mc - (mb + md))

	//I = M + N0
	si := m + n0

	//II = (ν/2) sinφ cosφ
	sii := (ν / 2) * sφ * cφ

	//III=(ν/24)sin(φ)cos^3(φ)(5−tan^2(φ)+9η^2)
	siii := (ν / 24) * sφ * c3φ * (5 - t2φ + 9*η2)

	//IIIA= (ν/720)sin(φ)cos^5(φ)(61−58tan^2(φ)+tan^4(φ))
	siiia := (ν / 720) * sφ * c5φ * (61 - 58*t2φ + t4φ)

	//IV =ν cosφ
	siv := ν * cφ

	//V=(v/6)cos^3(φ)((v/ρ)−tan^2(φ))
	sv := (ν / 6) * c3φ * ((ν / ρ) - t2φ)

	//VI= (ν/120) cos^5(φ)(5−18tan^2(φ)+tan^4(φ)+14η^2 −58tan^2(φη^2))
	svi := (ν / 120) * c5φ * (5 - 18*t2φ + t4φ + 14*η2 - 58*t2φ*η2)

	// λ - λ0
	dλ0 := λ - λ0
	// (λ - λ0)^2
	d2λ0 := dλ0 * dλ0
	// (λ - λ0)^3
	d3λ0 := d2λ0 * dλ0
	// (λ - λ0)^4
	d4λ0 := d2λ0 * d2λ0
	// (λ - λ0)^5
	d5λ0 := d3λ0 * d2λ0
	// (λ - λ0)^6
	d6λ0 := d3λ0 * d3λ0

	// (B7) N =I+II(λ−λ0)^2 +III(λ−λ0)^4 +IIIA(λ−λ0)^6
	northing := si + sii*d2λ0 + siii*d4λ0 + siiia*d6λ0

	// (B8) E = E0 +IV(λ−λ0)+V(λ−λ0)^3 +VI(λ−λ0)^5
	easting := e0 + siv*dλ0 + sv*d3λ0 + svi*d5λ0

	return &planeCoord{
		easting:  easting,
		northing: northing,
	}
}

func (proj *projection) fromPlaneCoord(coord *planeCoord, el *ellipsoid) (float64, float64) {
	// a - semi0major axis (metres)
	a := el.semiMajorAxis
	// b - semi-minor axis (metres)
	b := el.semiMinorAxis
	// e^2 - ellipsoid squared eccentricity constant.
	e2 := el.eccentricity()
	// N0 – northing of true origin;
	n0 := proj.mapTrueOrigin.northing
	// E0 – easting of true origin;
	e0 := proj.mapTrueOrigin.easting
	// F0 – scale factor on central meridian;
	f0 := proj.scaleFactor
	// φ0 – latitude of true origin; and
	φ0 := proj.geodeticTrueOrigin.lat
	// λ0 – longitude of true origin and central meridian.
	λ0 := proj.geodeticTrueOrigin.lon
	// (B2) n = a−b/a+b
	n := (a - b) / (a + b)

	φ := φ0
	m := 0.0

	for {

		// (C2) φnew = (N-N0-M)/(aF0) +φ′
		φ = φ + (coord.northing-(n0+m))/(a*f0)

		// n^2
		n2 := n * n
		// n^3
		n3 := n2 * n

		// φ - φ0
		dφ0 := φ - φ0
		// φ + φ0
		aφ0 := φ + φ0

		//Ma = (1+n+(5/4)n^2 +(5/4)n^3)(φ−φ0)
		ma := (1 + n + 5.0*n2/4.0 + 5.0*n3/4.0) * dφ0
		//Mb = (3n+3n^2 + (21/8)n^3)sin(φ−φ0)cos(φ+φ0)
		mb := (3*n + 3*n2 + (21.0/8.0)*n3) * math.Sin(dφ0) * math.Cos(aφ0)
		//Mc = ((15/8)n^2 + (15/8)n^3)sin(2(φ−φ0))cos(2(φ+φ0))
		mc := ((15.0/8.0)*n2 + (15.0/8.0)*n3) * math.Sin(2*dφ0) * math.Cos(2*(aφ0))
		//Md = (35/24)n^3 sin(3(φ −φ0))cos(3(φ +φ0))
		md := (35.0 / 24.0) * n3 * math.Sin(3*dφ0) * math.Cos(3*(aφ0))

		// (B6) M = bF0 (Ma - Mb + Mc - Md)
		m = b * f0 * (ma + mc - (mb + md))
		if coord.northing-(n0+m) < 0.00001 {
			break
		}
	}

	// sinφ
	sφ := math.Sin(φ)
	// cosφ
	cφ := math.Cos(φ)
	// sin^2(φ)
	s2φ := sφ * sφ
	// (B3) ν = aF0 (1 − e^2 sin^2(φ)) ^ −0.5
	ν := a * f0 * math.Pow(1-e2*s2φ, -0.5)
	// (B4) ρ=aF0(1−e^2)(1−e^2 sin^2(φ)) ^−1.5
	ρ := a * f0 * (1 - e2) * math.Pow(1-e2*s2φ, -1.5)
	// (B5) η2 = ν/ρ − 1
	η2 := ν/ρ - 1

	// tanφ
	tφ := math.Tan(φ)
	// tan^2φ
	t2φ := tφ * tφ
	// tan^4φ
	t4φ := t2φ * t2φ
	// tan^4φ
	t6φ := t4φ * t2φ
	// v^3
	ν3 := ν * ν * ν
	// v^5
	ν5 := ν3 * ν * ν
	// v^7
	ν7 := ν3 * ν3 * ν

	// VII= tanφ′/2ρν
	svii := tφ / (2 * ρ * ν)
	// VIII= tan(φ′)/24ρν^3 (5+3tan^2(φ)′+η2 −9tan^2(φ′)η2)
	sviii := (tφ / (24 * ρ * ν3)) * (5 + 3*t2φ + η2 - 9*t2φ*η2)
	// IX= (tan(φ′)/720ρν^5) (61+90tan^2(φ′)+45tan^4(φ′))
	six := (tφ / (720 * ρ * ν5)) * (61 + 90*t2φ + η2 + 45*t4φ)
	// X = sec(φ′)/ν
	sx := 1 / (ν * cφ)
	// XI= (sec(φ′)/6ν^3)(v/p)+2tan^2(φ′)
	sxi := (1 / (6 * ν3 * cφ)) * ((ν / ρ) + 2*t2φ)
	//XII= (sec(φ′)/120ν^5)(5+28tan^2(φ′)+24tan^4(φ′))
	sxii := (1 / (120 * ν5 * cφ)) * (5 + 28*t2φ + 24*t4φ)
	//XIIA= (sec(φ′)/5040ν^7)(61+662tan^2(φ′)+1320tan^4(φ′)+720tan^6(φ′))
	sxiia := (1 / (5040 * ν7 * cφ)) * (61 + 662*t2φ + 1320*t4φ + 720*t6φ)

	// E−E0
	de := coord.easting - e0
	// (E−E0)^2
	d2e := de * de
	// (E−E0)^3
	d3e := d2e * de
	// (E−E0)^4
	d4e := d2e * d2e
	// (E−E0)^5
	d5e := d3e * d2e
	// (E−E0)^6
	d6e := d3e * d3e
	// (E−E0)^7
	d7e := d4e * d3e
	// (C3) φ =φ′−VII(E−E0)^2 +VIII(E−E0)^4 −IX(E−E0)^6
	φ = φ - svii*d2e + sviii*d4e - six*d6e
	// (C4) λ=λ0 +X(E−E0)−XI(E−E0)^3 +XII(E−E0)^5 −XIIA(E−E0)^7
	λ := λ0 + sx*de - sxi*d3e + sxii*d5e - sxiia*d7e

	return φ, λ
}
