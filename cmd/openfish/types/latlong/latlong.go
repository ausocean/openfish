/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2025, The OpenFish Contributors.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions are met:

  1. Redistributions of source code must retain the above copyright notice, this
     list of conditions and the following disclaimer.

  2. Redistributions in binary form must reproduce the above copyright notice,
     this list of conditions and the following disclaimer in the documentation
     and/or other materials provided with the distribution.

  3. Neither the name of The Australian Ocean Lab Ltd. ("AusOcean")
     nor the names of its contributors may be used to endorse or promote
     products derived from this software without specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
  SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
  CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
  OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// latlong extends Google datastore's GeoPoint with support for parsing and serialising to a string.
package latlong

import (
	"fmt"
	"strconv"
	"strings"

	googlestore "cloud.google.com/go/datastore"
)

// VideoTime is a time in hours, minutes and seconds.
type LatLong struct {
	googlestore.GeoPoint
}

// String returns a string representation of the latlong in the format "-37.12345678,145.12345678".
func (l LatLong) String() string {
	return fmt.Sprintf("%.8f,%.8f", l.Lat, l.Lng)
}

// Parse converts a string to a LatLong, or throws an error.
func Parse(str string) (LatLong, error) {
	errMsg := "invalid location string: %w"

	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return LatLong{}, fmt.Errorf(errMsg, "string split failed")
	}
	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return LatLong{}, fmt.Errorf(errMsg, err)
	}
	long, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return LatLong{}, fmt.Errorf(errMsg, err)
	}
	return New(lat, long)
}

// UncheckedParse converts a string to a LatLong, or panics.
//
// This should only really be used when str is a string literal / constant,
// where the input is not dynamic. Prefer Parse when handling dynamic
// values, because it can return an error.
func UncheckedParse(str string) LatLong {
	vt, err := Parse(str)
	if err != nil {
		panic(err.Error())
	}
	return vt
}

// New creates a new LatLong.
func New(lat float64, long float64) (LatLong, error) {
	if lat < -90 || lat > 90 {
		return LatLong{}, fmt.Errorf("invalid value, latitude is not -90 ≤ %0.8f ≤ 90", lat)
	}
	if long < -180 || long > 180 {
		return LatLong{}, fmt.Errorf("invalid value, longitude is not -180 ≤ %0.8f ≤ 180", long)
	}

	return LatLong{GeoPoint: googlestore.GeoPoint{Lat: lat, Lng: long}}, nil
}

// UnmarshalText is used for decoding query params or JSON into a LatLong.
func (l *LatLong) UnmarshalText(text []byte) error {
	var err error
	*l, err = Parse(string(text))
	return err
}

// MarshalText is used for encoding a LatLong into JSON or query params.
func (l LatLong) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
