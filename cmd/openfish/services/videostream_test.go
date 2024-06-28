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

package services_test

import (
	"testing"
	"time"

	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"
)

// Constants.
var _8am = time.Date(2023, time.January, 1, 8, 0, 0, 0, time.UTC)
var _9am = time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC)
var _4pm = time.Date(2023, time.January, 1, 16, 0, 0, 0, time.UTC)
var _5pm = time.Date(2023, time.January, 1, 17, 0, 0, 0, time.UTC)

func TestCreateVideoStream(t *testing.T) {
	setup()

	// Create a new video stream entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	_, err := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})
	if err != nil {
		t.Errorf("Could not create video stream entity %s", err)
	}
}

func TestVideoStreamExists(t *testing.T) {
	setup()

	// Create a new video stream entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	id, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})

	// Check if the video stream exists.
	if !services.VideoStreamExists(int64(id)) {
		t.Errorf("Expected video stream to exist")
	}
}

func TestVideoStreamExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the video stream exists.
	// We expect it to return false.
	if services.VideoStreamExists(int64(123456789)) {
		t.Errorf("Did not expect video stream to exist")
	}
}

func TestGetVideoStreamByID(t *testing.T) {
	setup()

	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	id, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})

	videoStream, err := services.GetVideoStreamByID(int64(id))
	if err != nil {
		t.Errorf("Could not get video stream entity %s", err)
	}
	if videoStream.CaptureSource != int64(cs) && videoStream.StartTime != _8am && *videoStream.EndTime != _4pm && videoStream.StreamUrl != "http://youtube.com/watch?v=abc123" {
		t.Errorf("Video stream entity does not match created entity")
	}
}

func TestGetVideoStreamByIDForNonexistentEntity(t *testing.T) {
	setup()

	videoStream, err := services.GetVideoStreamByID(int64(123456789))
	if videoStream != nil && err == nil {
		t.Errorf("GetVideoStreamByID returned non-existing entity %s", err)
	}
}

// TODO: Write tests for GetVideoStreams. Test limit, offset and filtering.

func TestUpdateVideoStream(t *testing.T) {
	setup()
	// Create a new video stream entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	id, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})

	// Update the url.
	url := "http://youtube.com/watch?v=xyz789"
	err := services.UpdateVideoStream(int64(id), &url, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ := services.GetVideoStreamByID(int64(id))
	if videoStream.StreamUrl != url {
		t.Errorf("URL did not update, expected %s, actual %s", url, videoStream.StreamUrl)
	}

	// Update the capture source.
	cs2, _ := services.CreateCaptureSource("Stony Point camera 2", 0.0, 0.0, "RPI camera", nil)
	csnew := int64(cs2)
	err = services.UpdateVideoStream(int64(id), nil, &csnew, nil, nil, nil)
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ = services.GetVideoStreamByID(int64(id))
	if videoStream.CaptureSource != csnew {
		t.Errorf("Capture source did not update, expected %d, actual %d", csnew, videoStream.CaptureSource)
	}

	// Update stream times.
	err = services.UpdateVideoStream(int64(id), nil, nil, &_9am, &_5pm, nil)
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ = services.GetVideoStreamByID(int64(id))
	if !videoStream.StartTime.Equal(_9am) || !videoStream.EndTime.Equal(_5pm) {
		t.Errorf("Start / end times did not update, expected %s %s, actual %s %s",
			_9am.Format(time.RFC3339),
			_5pm.Format(time.RFC3339),
			videoStream.StartTime.Format(time.RFC3339),
			videoStream.EndTime.Format(time.RFC3339),
		)
	}

}

func TestUpdateVideoStreamForNonExistentEntity(t *testing.T) {
	setup()

	err := services.UpdateVideoStream(int64(123456789), nil, nil, nil, nil, nil)
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent video stream")
	}
}

func TestDeleteVideoStream(t *testing.T) {
	setup()

	// Create a new video stream entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	id, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})

	// Delete the video stream entity.
	err := services.DeleteVideoStream(int64(id))
	if err != nil {
		t.Errorf("Could not delete video stream entity %d: %s", id, err)
	}

	// Check if the video stream exists.
	if services.VideoStreamExists(int64(id)) {
		t.Errorf("Video stream entity exists after delete")
	}
}

func TestDeleteVideoStreamForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteVideoStream(int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent video stream")
	}
}

func TestDeleteVideoStreamWithAssociatedAnnotations(t *testing.T) {
	// TODO: Run test in CI when issue is fixed.
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	setup()

	// Create a new video stream entity and an annotation that references it.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	id, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []string{})
	services.CreateAnnotation(id,
		entities.TimeSpan{Start: _8am, End: _9am},
		nil, "scott@ausocean.org",
		map[string]string{"species": "Sepia Apama"})

	err := services.DeleteVideoStream(int64(id))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting video stream with associated annotation")
	}
}
