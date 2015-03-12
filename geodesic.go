// Copyright 2015 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package geodesic provides direct and inverse geodesic calculations.

It is a Go wrapper around the GeographicLib geodesic C library by Charles Karney.
The original library (MIT licensed) can be found at http://geographiclib.sourceforge.net/

*/
package geodesic

// #cgo LDFLAGS: -lm
// #include "geodesic.h"
import "C"

//Point represents a location on earth
type Point struct {
	LatitudeDeg  float64
	LongitudeDeg float64
}

//Invert an azimuth in degree (adds 180° and convert to [0,360[ interval).
func InvertAzimuth(azDeg float64) float64 {
	return azDeg + 180%360
}

//Geodesic contains information about the ellipsoid
type Geodesic struct {
	g C.struct_geod_geodesic
}

//Geodesic contains information about the ellipsoid
func NewGeodesic(f float64, a float64) Geodesic {
	var g Geodesic
	C.geod_init(&g.g, C.double(f), C.double(a))
	return g
}

//WGS84 represents the WGS 84 ellipsoid
var WGS84 = NewGeodesic(6378137, 1/298.257223563)

//Direct computes the point pt2 obtained from pt1 at a distance s12 (in meters) in
//direction of az1 (in degrees, clockwise from north). The returned azimuth az2 is the
//azimuth of the geodesic at pt2 (in degrees).
func (g Geodesic) Direct(pt1 Point, az1 float64, s12 float64) (pt2 Point, az2 float64) {

	var resLat, resLon, resAz C.double

	C.geod_direct(
		&g.g,
		C.double(pt1.LatitudeDeg),
		C.double(pt1.LongitudeDeg),
		C.double(az1),
		C.double(s12),
		&resLat,
		&resLon,
		&resAz)

	pt2.LatitudeDeg = float64(resLat)
	pt2.LongitudeDeg = float64(resLon)
	az2 = float64(resAz)

	return
}

//Inverse computes the azimuths (in degrees, clockwise from north) and the distance (in meters) between two points.
//Azimuths are from one point to the other.
func (g *Geodesic) Inverse(pt1 Point, pt2 Point) (s12, az1, az2 float64) {

	var resS12, resAz1, resAz2 C.double

	C.geod_inverse(
		&g.g,
		C.double(pt1.LatitudeDeg),
		C.double(pt1.LongitudeDeg),
		C.double(pt2.LatitudeDeg),
		C.double(pt2.LongitudeDeg),
		&resS12,
		&resAz1,
		&resAz2)

	s12 = float64(resS12)
	az1 = float64(resAz1)
	az2 = float64(resAz2)

	return
}

//GeodesicLine contains information about a single geodesic line.
type GeodesicLine struct {
	gl C.struct_geod_geodesicline
}

//NewGeodesicLine initialize a GeodesicLine.
func NewGeodesicLine(g Geodesic, origin Point, azDeg float64) GeodesicLine {
	var l GeodesicLine
	C.geod_lineinit(&l.gl, &g.g, C.double(origin.LatitudeDeg), C.double(origin.LongitudeDeg), C.double(azDeg), 0)
	return l
}

//Position computes the point at s12 (meters) from the GeodesicLine origin. The distance can be negativ.
//It is faster than multiple calls to Direct.
//Azimuths is the azimuth of the geodesic line at the resulting point. For short positive distance it
//will be approximately equal to the GeodesicLine azimuth at origin.
func (l GeodesicLine) Position(s12 float64) (pt Point, azDeg float64) {

	var resLat, resLon, resAzDeg C.double

	C.geod_position(
		&l.gl,
		C.double(s12),
		&resLat,
		&resLon,
		&resAzDeg)

	pt.LatitudeDeg = float64(resLat)
	pt.LongitudeDeg = float64(resLon)
	azDeg = float64(resAzDeg)

	return
}
