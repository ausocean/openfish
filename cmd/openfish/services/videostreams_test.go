/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023-2025, The OpenFish Contributors.

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
	"time"

	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/ausocean/openfish/cmd/openfish/types/timezone"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// Constants.
var _8am = time.Date(2023, time.January, 1, 8, 0, 0, 0, time.UTC)
var _9am = time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC)
var _4pm = time.Date(2023, time.January, 1, 16, 0, 0, 0, time.UTC)
var _5pm = time.Date(2023, time.January, 1, 17, 0, 0, 0, time.UTC)

func createTestVideoStream() services.VideoStream {
	cs := createTestCaptureSource()
	vs, _ := services.CreateVideoStream(services.VideoStreamContents{
		StartTime:     _8am,
		EndTime:       &_4pm,
		AnnotatorList: []int64{},
		BaseVideoStreamFields: services.BaseVideoStreamFields{
			TimeZone:      timezone.UncheckedParse("Australia/Adelaide"),
			StreamURL:     "http://youtube.com/watch?v=abc123",
			CaptureSource: cs.ID,
		},
	})
	return *vs
}

func TestCreateVideoStream(t *testing.T) {
	setup()

	// Create a new video stream entity.
	cs := createTestCaptureSource()
	_, err := services.CreateVideoStream(services.VideoStreamContents{
		StartTime:     _8am,
		EndTime:       &_4pm,
		AnnotatorList: []int64{},
		BaseVideoStreamFields: services.BaseVideoStreamFields{
			TimeZone:      timezone.UncheckedParse("Australia/Adelaide"),
			StreamURL:     "http://youtube.com/watch?v=abc123",
			CaptureSource: cs.ID,
		},
	})
	if err != nil {
		t.Errorf("Could not create video stream entity %s", err)
	}
}

func TestVideoStreamExists(t *testing.T) {
	setup()

	// Create a new video stream entity.
	vs := createTestVideoStream()

	// Check if the video stream exists.
	if !services.VideoStreamExists(vs.ID) {
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

	vs := createTestVideoStream()

	videoStream, err := services.GetVideoStreamByID(vs.ID)
	if err != nil {
		t.Errorf("Could not get video stream entity %s", err)
	}
	if !reflect.DeepEqual(vs, *videoStream) {
		t.Errorf("Video stream entity does not match created entity, %v, %v", vs, *videoStream)
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
	vs := createTestVideoStream()

	// Update the url.
	url := "http://youtube.com/watch?v=xyz789"
	err := services.UpdateVideoStream(vs.ID, services.PartialVideoStreamContents{StreamURL: &url})
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ := services.GetVideoStreamByID(vs.ID)
	if videoStream.StreamURL != url {
		t.Errorf("URL did not update, expected %s, actual %s", url, videoStream.StreamURL)
	}

	// Update the capture source.
	cs2 := createTestCaptureSource()
	err = services.UpdateVideoStream(vs.ID, services.PartialVideoStreamContents{CaptureSource: &cs2.ID})
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ = services.GetVideoStreamByID(vs.ID)
	if videoStream.CaptureSource != cs2.ID {
		t.Errorf("Capture source did not update, expected %d, actual %d", cs2.ID, videoStream.CaptureSource)
	}

	// Update stream times.
	start, _ := time.Parse(time.RFC3339, "2025-03-21T13:58:00.000+10:30")
	end, _ := time.Parse(time.RFC3339, "2025-03-21T15:58:00.000+10:30")
	timezone := timezone.UncheckedParse("Australia/Adelaide")
	err = services.UpdateVideoStream(vs.ID, services.PartialVideoStreamContents{StartTime: &start, EndTime: &end, TimeZone: &timezone})
	if err != nil {
		t.Errorf("Could not update video stream entity %s", err)
	}

	videoStream, _ = services.GetVideoStreamByID(vs.ID)
	if !videoStream.StartTime.Equal(start) || !videoStream.EndTime.Equal(end) || videoStream.TimeZone.String() != "Australia/Adelaide" {
		t.Errorf(
			"Start / end times did not update, expected: %s %s %s, got: %s %s %s",
			start.String(),
			end.String(),
			timezone.String(),
			videoStream.StartTime.String(),
			videoStream.EndTime.String(),
			videoStream.TimeZone.String(),
		)
	}
}

func TestUpdateVideoStreamForNonExistentEntity(t *testing.T) {
	setup()

	url := "http://youtube.com/watch?v=xyz789"
	err := services.UpdateVideoStream(int64(123456789), services.PartialVideoStreamContents{StreamURL: &url})
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent video stream")
	}
}

func TestDeleteVideoStream(t *testing.T) {
	setup()

	// Create a new video stream entity.
	vs := createTestVideoStream()

	// Delete the video stream entity.
	err := services.DeleteVideoStream(vs.ID)
	if err != nil {
		t.Errorf("Could not delete video stream entity %d: %s", vs.ID, err)
	}

	// Check if the video stream exists.
	if services.VideoStreamExists(vs.ID) {
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
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})

	// Create a new video stream entity.
	vs := createTestVideoStream()
	services.CreateAnnotation(services.AnnotationContents{
		KeyPoints: []keypoint.KeyPoint{
			{
				BoundingBox: keypoint.BoundingBox{X1: 10, X2: 20, Y1: 70, Y2: 80},
				Time:        videotime.UncheckedParse("00:00:01.000"),
			},
			{
				BoundingBox: keypoint.BoundingBox{X1: 20, X2: 30, Y1: 60, Y2: 70},
				Time:        videotime.UncheckedParse("00:00:02.000"),
			},
		},
		VideostreamID: vs.ID,
		CreatedByID:   uid,
	})

	err := services.DeleteVideoStream(vs.ID)
	if err == nil {
		t.Errorf("Did not receive expected error when deleting video stream with associated annotation")
	}
}
