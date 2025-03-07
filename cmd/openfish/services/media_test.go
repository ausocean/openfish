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

package services_test

import (
	"reflect"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/mediatype"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

func TestCreateImageMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	vs := createTestVideoStream()
	_, err := services.CreateMedia(services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: vs.ID,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}, []byte{})
	if err != nil {
		t.Errorf("Could not create media entity %s", err)
	}
}

func TestCreateVideoMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	vs := createTestVideoStream()
	start := videotime.UncheckedParse("00:00:01.000")
	end := videotime.UncheckedParse("00:00:01.500")
	_, err := services.CreateMedia(services.MediaKey{
		Type:          mediatype.MP4,
		VideoStreamID: vs.ID,
		StartTime:     start,
		EndTime:       &end,
	}, []byte{})
	if err != nil {
		t.Errorf("Could not create media entity %s", err)
	}
}

func TestMediaExists(t *testing.T) {
	setup()

	// Create a new media entity.
	vs := createTestVideoStream()
	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: vs.ID,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	services.CreateMedia(mq, []byte{})

	// Check if the media exists.
	if !services.MediaExists(mq) {
		t.Errorf("Expected media to exist")
	}
}

func TestMediaExistsForNonexistentObject(t *testing.T) {
	setup()

	// Check if the media exists.
	// We expect it to return false.
	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: 1234567890,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	if services.MediaExists(mq) {
		t.Errorf("Did not expect media to exist")
	}
}

func TestGetMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	vs := createTestVideoStream()
	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: vs.ID,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	expected := []byte{1, 2, 3, 4, 5}
	services.CreateMedia(mq, expected)

	bytes, err := services.GetMedia(mq)
	if err != nil {
		t.Errorf("Could not get media object: %s", err)
	}
	if !reflect.DeepEqual(expected, bytes) {
		t.Errorf("Media object does not match expected object, %+v, %+v", bytes, expected)
	}
}

func TestGetMediaForNonexistentObject(t *testing.T) {
	setup()

	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: 1234567890,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	_, err := services.GetMedia(mq)
	if err == nil {
		t.Errorf("GetMediaByID returned non-existing object %s", err)
	}
}

func TestDeleteMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	vs := createTestVideoStream()
	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: vs.ID,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	services.CreateMedia(mq, []byte{})

	// Delete the media entity.
	err := services.DeleteMedia(mq)
	if err != nil {
		t.Errorf("Could not delete media object: %s", err)
	}

	// Check if the media exists.
	if services.MediaExists(mq) {
		t.Errorf("Video stream entity exists after delete")
	}
}

func TestDeleteMediaForNonexistentObject(t *testing.T) {
	setup()

	mq := services.MediaKey{
		Type:          mediatype.JPEG,
		VideoStreamID: 1234567890,
		StartTime:     videotime.UncheckedParse("00:00:01.000"),
		EndTime:       nil,
	}
	err := services.DeleteMedia(mq)
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent media")
	}
}
