/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2024, The OpenFish Contributors.

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

// videotime provides a type VideoTime that represents the time in a video - hours, minutes and seconds.
// Unlike the time package, there is no concept of dates, and no timezones. VideoTimes are displayed in
// the format 12:01:04.000, and can be serialized/deserialized to and from JSON in this format also.
package videotime

import (
	"fmt"
)

// VideoTime is a time in hours, minutes and seconds.
type VideoTime struct {
	value int64
}

// String converts a VideoTime to a string.
func (t VideoTime) String() string {

	// Milliseconds are what remain.
	ms := t.value

	// Calculate hours.
	h := ms / (60 * 60 * 1000)
	ms -= h * (60 * 60 * 1000)

	// Calculate minutes.
	m := ms / (60 * 1000)
	ms -= m * (60 * 1000)

	// Calculate seconds
	s := ms / 1000
	ms -= s * 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// Int converts a VideoTime to an int64 (seconds).
func (t VideoTime) Int() int64 {
	return t.value
}

// Parse converts a string to a VideoTime, or throws an error.
func Parse(str string) (VideoTime, error) {
	var h, m, s, ms int

	_, err := fmt.Sscanf(str, "%02d:%02d:%02d.%03d", &h, &m, &s, &ms)
	if err != nil {
		return VideoTime{}, fmt.Errorf("invalid format: %s, %v", str, err)
	}

	return New(h, m, s, ms)
}

// New creates a new VideoTime.
func New(h int, m int, s int, ms int) (VideoTime, error) {
	if ms < 0 || ms >= 1000 {
		return VideoTime{}, fmt.Errorf("invalid value, milliseconds is not 0 ≤ %d ≤ 999", s)
	}
	if s < 0 || s >= 60 {
		return VideoTime{}, fmt.Errorf("invalid value, seconds is not 0 ≤ %d ≤ 59", s)
	}
	if m < 0 || m >= 60 {
		return VideoTime{}, fmt.Errorf("invalid value, minutes is not 0 ≤ %d ≤ 59", m)
	}
	return VideoTime{value: int64(h*60*60*1000 + m*60*1000 + s*1000 + ms)}, nil
}

func FromInt(i int64) VideoTime {
	return VideoTime{value: i}
}

// UnmarshalText is used for decoding query params or JSON into a VideoTime.
func (t *VideoTime) UnmarshalText(text []byte) error {
	var err error
	*t, err = Parse(string(text))
	return err
}

// MarshalText is used for encoding a VideoTime into JSON or query params.
func (t VideoTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
