/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023-2024, The OpenFish Contributors.

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

// Timespan struct represents a start and end time in a video.
package timespan

import (
	"fmt"
	"strings"

	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// TimeSpan is a pair of video timestamps - start time and end time.
type TimeSpan struct {
	Start videotime.VideoTime
	End   videotime.VideoTime
}

// Valid tests if a timespan is valid. Start should be less than End.
func (t TimeSpan) Valid() bool {
	return t.Start.Int() <= t.End.Int()
}

// String returns a string representation of the timespan in the format "12:01:04.000-12:01:05.000".
func (t TimeSpan) String() string {
	return fmt.Sprintf("%s-%s", t.Start.String(), t.End.String())
}

// Parse takes a string in the format "12:01:04.000-12:01:05.000" and returns a TimeSpan.
// Returns an error if the string is not in the correct format or if the timestamps are invalid.
func Parse(s string) (*TimeSpan, error) {
	str := strings.Split(s, "-")
	if len(str) != 2 {
		return nil, fmt.Errorf("invalid timespan")
	}
	start, err := videotime.Parse(str[0])
	if err != nil {
		return nil, err
	}
	end, err := videotime.Parse(str[0])
	if err != nil {
		return nil, err
	}
	t := TimeSpan{Start: start, End: end}
	return &t, nil
}

// UnmarshalText is used for decoding query params or JSON into a TimeSpan.
func (t *TimeSpan) UnmarshalText(text []byte) error {
	var err error
	t, err = Parse(string(text))
	return err
}

// MarshalText is used for encoding a TimeSpan into JSON or query params.
func (t TimeSpan) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
