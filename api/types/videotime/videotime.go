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

// videotime provides a type VideoTime that represents the time in a video - hours, minutes and seconds.
// Unlike the time package, there is no concept of dates, and no timezones. VideoTimes are displayed in
// the format 12:01:04, and can be serialized/deserialized to and from JSON in this format also.
package videotime

import (
	"fmt"
	"strconv"
	"strings"
)

// VideoTime is a time in hours, minutes and seconds.
type VideoTime struct {
	value int64
}

func (t VideoTime) String() string {
	// Seconds are what remain.
	s := t.value

	// Calculate hours.
	h := s / (60 * 60)
	s -= h * (60 * 60)

	// Calculate minutes.
	m := s / 60
	s -= m * 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (t VideoTime) Int() int64 {
	return t.value
}

func Parse(s string) (VideoTime, error) {
	// Split on ':' characters.
	str := strings.Split(s, ":")
	if len(str) != 3 {
		return VideoTime{}, fmt.Errorf("invalid format: %s, must have three parts: hh:mm:ss", s)
	}

	// Parse as integers.
	var nums [3]int
	var err error
	for i := 0; i < 3; i++ {
		nums[i], err = strconv.Atoi(str[i])
		if err != nil {
			return VideoTime{}, fmt.Errorf("invalid format: %s, must only have numbers between ':' separator", s)
		}
	}

	return New(nums[0], nums[1], nums[2])
}

func New(h int, m int, s int) (VideoTime, error) {
	if s < 0 || s >= 60 {
		return VideoTime{}, fmt.Errorf("invalid value, seconds is not 0 ≤ %d ≤ 59", s)
	}
	if m < 0 || m >= 60 {
		return VideoTime{}, fmt.Errorf("invalid value, minutes is not 0 ≤ %d ≤ 59", m)
	}
	return VideoTime{value: int64(h*60*60 + m*60 + s)}, nil
}

func (t *VideoTime) UnmarshalText(text []byte) error {
	var err error
	*t, err = Parse(string(text))
	return err
}

func (t VideoTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
