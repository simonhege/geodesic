// Copyright 2015 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package geodesic

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"math"
	"os"
	"testing"
)

var (
	deltaAz  = flag.Float64("deltaAz", 1e-2, "Maximum delta allowed on azimuths (degree)")
	deltaS   = flag.Float64("deltaS", 1e-7, "Maximum delta allowed on distances (m)")
	deltaPos = flag.Float64("deltaPos", 1e-8, "Maximum delta allowed on latitudes or longitudes (degree)")
)

func approxEq(a, b float64, delta float64) bool {
	return math.Abs(a-b) < delta
}

type d struct {
	lat1 float64
	lon1 float64
	azi1 float64
	lat2 float64
	lon2 float64
	azi2 float64
	s12  float64
	a12  float64
	m12  float64
	S12  float64
}

var testData = []d{}

func TestMain(m *testing.M) {
	flag.Parse()

	filepath := "testdata/GeodTest.dat.gz"
	if testing.Short() {
		filepath = "testdata/GeodTest-short.dat.gz"
	}

	//Load test data file
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	unziper, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer unziper.Close()
	reader := bufio.NewReader(unziper)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		var d d
		_, err := fmt.Sscan(scanner.Text(), &d.lat1, &d.lon1, &d.azi1, &d.lat2, &d.lon2, &d.azi2, &d.s12, &d.a12, &d.m12, &d.S12)
		if err != nil {
			panic(err)
		}
		testData = append(testData, d)
	}

	//Run test suite
	os.Exit(m.Run())
}

func TestDirect(t *testing.T) {

	for i, tt := range testData {

		outPt2, outAz2 := WGS84.Direct(Point{tt.lat1, tt.lon1}, tt.azi1, tt.s12)

		if !approxEq(outPt2.LatitudeDeg, tt.lat2, *deltaPos) {
			t.Errorf("%d: Direct(%v, %v, %v).LatitudeDeg => %v, want %v", i, tt.lat1, tt.lon1, tt.azi1, outPt2.LatitudeDeg, tt.lat2)
		}
		if !approxEq(outPt2.LongitudeDeg, tt.lon2, *deltaPos) {
			t.Errorf("%d: Direct(%v, %v, %v).LongitudeDeg => %v, want %v", i, tt.lat1, tt.lon1, tt.azi1, outPt2.LongitudeDeg, tt.lon2)
		}
		if !approxEq(outAz2, tt.azi2, *deltaAz) {
			t.Errorf("%d: Direct(%v, %v, %v).Azimuth => %v, want %v", i, tt.lat1, tt.lon1, tt.azi1, outAz2, tt.azi2)
		}
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tt := testData[i%len(testData)]
		WGS84.Direct(Point{tt.lat1, tt.lon1}, tt.azi1, tt.s12)
	}
}

func TestInverse(t *testing.T) {

	for i, tt := range testData {

		outS12, outAz1, outAz2 := WGS84.Inverse(Point{tt.lat1, tt.lon1}, Point{tt.lat2, tt.lon2})

		if !approxEq(outS12, tt.s12, *deltaS) {
			t.Errorf("%d: Inverse(%v, %v, %v, %v).s12 => %v, want %v", i, tt.lat1, tt.lon1, tt.lat2, tt.lon2, outS12, tt.s12)
		}
		if !approxEq(outAz1, tt.azi1, *deltaAz) {
			t.Errorf("%d: Inverse(%v, %v, %v, %v).azi1 => %v, want %v", i, tt.lat1, tt.lon1, tt.lat2, tt.lon2, outAz1, tt.azi1)
		}
		if !approxEq(outAz2, tt.azi2, *deltaAz) {
			t.Errorf("%d: Inverse(%v, %v, %v, %v).azi2 => %v, want %v", i, tt.lat1, tt.lon1, tt.lat2, tt.lon2, outAz2, tt.azi2)
		}
	}
}
func BenchmarkInverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tt := testData[i%len(testData)]
		WGS84.Inverse(Point{tt.lat1, tt.lon1}, Point{tt.lat2, tt.lon2})
	}
}
