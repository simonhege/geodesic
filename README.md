# Geodesic

Package geodesic provides direct and inverse geodesic calculations on any ellipsoid, especially on WGS 84.

It is a Go wrapper around the GeographicLib geodesic C library by Charles Karney.
The original library (MIT licensed) can be found at http://geographiclib.sourceforge.net/

## Install

	go get github.com/xeonx/geodesic

## Docs

<http://godoc.org/github.com/xeonx/geodesic>
	
## Tests

`go test` is used for testing.

Run partial test suite:

	go test github.com/xeonx/geodesic -short
	
Run full test suite:

	go test github.com/xeonx/geodesic
	
Run full test suite, with given error thresholds:

	go test github.com/xeonx/geodesic -deltaPos=5e-8 -deltaAz=0.02 -deltaS=1e-7
	
Run benchmarks:

	go test github.com/xeonx/geodesic -short -bench .

## License

This code is licensed under the MIT license. See [LICENSE](https://github.com/xeonx/geodesic/blob/master/LICENSE).

Files from GeographicLib are Copyright (c) Charles Karney (2012-2014) <charles@karney.com> and licensed under the MIT/X11 License.