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

package videotime_test

import (
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

func TestString(t *testing.T) {
	vt, _ := videotime.New(12, 34, 56)
	expected := "12:34:56"

	if vt.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, vt.String())
	}
}

func TestParse(t *testing.T) {
	str := "12:34:56"
	expected, _ := videotime.New(12, 34, 56)

	vt, err := videotime.Parse(str)
	if err != nil {
		t.Errorf("Error parsing time: %v", err)
	}

	if vt != expected {
		t.Errorf("Expected value %s, but got %s", expected.String(), vt.String())
	}
}

func TestUnmarshalText(t *testing.T) {
	str := "12:34:56"
	expected, _ := videotime.New(12, 34, 56)

	vt := videotime.VideoTime{}

	err := vt.UnmarshalText([]byte(str))
	if err != nil {
		t.Errorf("Error unmarshalling text: %v", err)
	}

	if vt != expected {
		t.Errorf("Expected value %s, but got %s", expected.String(), vt.String())
	}
}

func TestMarshalText(t *testing.T) {
	vt, _ := videotime.New(12, 34, 56)
	expected := "12:34:56"

	data, err := vt.MarshalText()
	if err != nil {
		t.Errorf("Error marshalling text: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected text %s, but got %s", expected, string(data))
	}
}
