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
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	_, err := services.CreateMedia(services.MediaContents{
		Type:              mediatype.JPEG,
		VideoStreamSource: vs,
		StartTime:         videotime.UncheckedParse("00:00:01.000"),
		EndTime:           nil,
		Bytes:             []byte{},
	})
	if err != nil {
		t.Errorf("Could not create media entity %s", err)
	}
}

func TestCreateVideoMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	end := videotime.UncheckedParse("00:00:01.500")
	_, err := services.CreateMedia(services.MediaContents{
		Type:              mediatype.MP4,
		VideoStreamSource: vs,
		StartTime:         videotime.UncheckedParse("00:00:01.000"),
		EndTime:           &end,
		Bytes:             []byte{},
	})
	if err != nil {
		t.Errorf("Could not create media entity %s", err)
	}
}

func TestMediaExists(t *testing.T) {
	setup()

	// Create a new media entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	id, _ := services.CreateMedia(services.MediaContents{
		Type:              mediatype.JPEG,
		VideoStreamSource: vs,
		StartTime:         videotime.UncheckedParse("00:00:01.000"),
		EndTime:           nil,
		Bytes:             []byte{},
	})

	// Check if the media exists.
	if !services.MediaExists(id) {
		t.Errorf("Expected media to exist")
	}
}

func TestMediaExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the media exists.
	// We expect it to return false.
	if services.MediaExists(int64(123456789)) {
		t.Errorf("Did not expect media to exist")
	}
}

func TestGetMediaByID(t *testing.T) {
	setup()

	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	contents := services.MediaContents{
		Type:              mediatype.JPEG,
		VideoStreamSource: vs,
		StartTime:         videotime.UncheckedParse("00:00:01.000"),
		EndTime:           nil,
		Bytes:             []byte{},
	}
	id, _ := services.CreateMedia(contents)

	media, err := services.GetMediaByID(id)
	if err != nil {
		t.Errorf("Could not get media entity %s", err)
	}
	if !reflect.DeepEqual(media.MediaContents, contents) {
		t.Errorf("Video stream entity does not match created entity, %+v, %+v", media.MediaContents, contents)
	}
}

func TestGetMediaByIDForNonexistentEntity(t *testing.T) {
	setup()

	videoStream, err := services.GetMediaByID(int64(123456789))
	if videoStream != nil && err == nil {
		t.Errorf("GetMediaByID returned non-existing entity %s", err)
	}
}

func TestDeleteMedia(t *testing.T) {
	setup()

	// Create a new media entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	id, _ := services.CreateMedia(services.MediaContents{
		Type:              mediatype.JPEG,
		VideoStreamSource: vs,
		StartTime:         videotime.UncheckedParse("00:00:01.000"),
		EndTime:           nil,
		Bytes:             []byte{},
	})

	// Delete the media entity.
	err := services.DeleteMedia(int64(id))
	if err != nil {
		t.Errorf("Could not delete media entity %d: %s", id, err)
	}

	// Check if the media exists.
	if services.MediaExists(int64(id)) {
		t.Errorf("Video stream entity exists after delete")
	}
}

func TestDeleteMediaForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteMedia(int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent media")
	}
}
