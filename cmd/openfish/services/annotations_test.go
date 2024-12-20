/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023, The OpenFish Contributors.

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

	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// Constants
var startTime, _ = videotime.Parse("00:00:01")
var endTime, _ = videotime.Parse("00:01:00")

var startKp = keypoint.KeyPoint{
	BoundingBox: keypoint.BoundingBox{X1: 10, X2: 20, Y1: 70, Y2: 80},
	Time:        startTime,
}

var endKp = keypoint.KeyPoint{
	BoundingBox: keypoint.BoundingBox{X1: 10, X2: 20, Y1: 70, Y2: 80},
	Time:        endTime,
}

func TestCreateAnnotation(t *testing.T) {
	setup()

	// Create a new annotation entity.
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]entities.Image, 0), nil)
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	_, err := services.CreateAnnotation(vs,
		[]keypoint.KeyPoint{startKp, endKp},
		uid,
		map[string]string{"species": "Sepioteuthis australis", "common_name": "Southern Reef Squid"})

	if err != nil {
		t.Errorf("Could not create annotation entity %s", err)
	}
}

func TestAnnotationExists(t *testing.T) {
	setup()

	// Create a new annotation entity.
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]entities.Image, 0), nil)
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	id, _ := services.CreateAnnotation(vs,
		[]keypoint.KeyPoint{startKp, endKp},
		uid,
		map[string]string{"species": "Sepioteuthis australis", "common_name": "Southern Reef Squid"})

	// Check if the annotation exists.
	if !services.AnnotationExists(int64(id)) {
		t.Errorf("Expected annotation to exist")
	}
}

func TestAnnotationExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the annotation exists.
	// We expect it to return false.
	if services.AnnotationExists(int64(123456789)) {
		t.Errorf("Did not expect annotation to exist")
	}
}

func TestGetAnnotationByID(t *testing.T) {
	setup()

	// Create a new annotation entity.
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]entities.Image, 0), nil)
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	id, _ := services.CreateAnnotation(vs,
		[]keypoint.KeyPoint{startKp, endKp},
		uid,
		map[string]string{"species": "Sepioteuthis australis", "common_name": "Southern Reef Squid"})

	annotation, err := services.GetAnnotationByID(int64(id))
	if err != nil {
		t.Errorf("Could not get annotation entity %s", err)
	}
	if annotation.VideoStreamID != vs || annotation.Observer != uid {
		// TODO: verify keypoints and observations are the same
		t.Errorf("Annotation entity does not match created entity")
	}
}

func TestGetAnnotationByIDForNonexistentEntity(t *testing.T) {
	setup()

	annotation, err := services.GetAnnotationByID(int64(123456789))
	if annotation != nil && err == nil {
		t.Errorf("GetAnnotationByID returned non-existing entity %s", err)
	}
}

// TODO: Write tests for GetAnnotations. Test limit, offset and filtering.

func TestDeleteAnnotation(t *testing.T) {
	setup()

	// Create a new annotation entity.
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]entities.Image, 0), nil)
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm, []int64{})
	id, _ := services.CreateAnnotation(vs,
		[]keypoint.KeyPoint{startKp, endKp},
		uid,
		map[string]string{"species": "Sepioteuthis australis", "common_name": "Southern Reef Squid"})

	// Delete the annotation entity.
	err := services.DeleteAnnotation(int64(id))
	if err != nil {
		t.Errorf("Could not delete annotation entity %d: %s", id, err)
	}

	// Check if the annotation exists.
	if services.AnnotationExists(int64(id)) {
		t.Errorf("Video stream entity exists after delete")
	}
}

func TestDeleteAnnotationForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteAnnotation(int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent annotation")
	}
}
