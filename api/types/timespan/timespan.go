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
// Uses a custom Save and Load method for reading/writing to the datastore.
// TODO: add timespan validation to type (end > start).
package timespan

import (
	googlestore "cloud.google.com/go/datastore"
	"github.com/ausocean/openfish/api/types/videotime"
	"github.com/ausocean/openfish/datastore"
)

// TimeSpan is a pair of video timestamps - start time and end time.
type TimeSpan struct {
	Start videotime.VideoTime `json:"start"`
	End   videotime.VideoTime `json:"end"`
}

// Valid tests if a timespan is valid. Start should be less than End.
func (t TimeSpan) Valid() bool {
	return t.Start.Int() <= t.End.Int()
}

func (t *TimeSpan) Load(ps []datastore.Property) error {

	var data struct {
		Start string
		End   string
	}

	if err := googlestore.LoadStruct(&data, ps); err != nil {
		return err
	}

	t.Start, _ = videotime.Parse(data.Start)
	t.End, _ = videotime.Parse(data.End)

	return nil
}

func (t *TimeSpan) Save() ([]datastore.Property, error) {
	return []datastore.Property{
		{
			Name:  "Start",
			Value: t.Start.String(),
		},
		{
			Name:  "End",
			Value: t.End.String(),
		},
	}, nil
}
