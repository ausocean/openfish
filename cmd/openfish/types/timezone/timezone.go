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

// timezone extends time.Location so it can be serialised and deserialised
package timezone

import (
	"time"
)

type TimeZone struct {
	time.Location
}

// Parse parses a string as a TimeZone.
func Parse(str string) (TimeZone, error) {
	l, err := time.LoadLocation(str)
	return TimeZone{Location: *l}, err
}

// UncheckedParse converts a string to a TimeZone, or panics.
//
// This should only really be used when str is a string literal / constant,
// where the input is not dynamic. Prefer Parse when handling dynamic
// values, because it can return an error.
func UncheckedParse(str string) TimeZone {
	vt, err := Parse(str)
	if err != nil {
		panic(err.Error())
	}
	return vt
}

// UnmarshalText is used for decoding query params or JSON into a TimeZone.
func (tz *TimeZone) UnmarshalText(text []byte) error {
	var err error
	*tz, err = Parse(string(text))
	return err
}

// MarshalText is used for encoding a TimeZone into JSON or query params.
func (tz TimeZone) MarshalText() ([]byte, error) {
	return []byte(tz.String()), nil
}
