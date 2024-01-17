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

	"github.com/ausocean/openfish/api/entities"
	"github.com/ausocean/openfish/api/services"
)

func TestCreateAnnotation(t *testing.T) {
	setup()

	// Create a new annotation entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm)
	_, err := services.CreateAnnotation(vs,
		entities.TimeSpan{Start: _8am, End: _9am},
		nil, "scott@ausocean.org",
		map[string]string{"species": "Sepia Apama"})

	if err != nil {
		t.Errorf("Could not create annotation entity %s", err)
	}
}

func TestAnnotationExists(t *testing.T) {
	setup()

	// Create a new annotation entity.
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm)
	id, _ := services.CreateAnnotation(vs,
		entities.TimeSpan{Start: _8am, End: _9am},
		nil, "scott@ausocean.org",
		map[string]string{"species": "Sepia Apama"})

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
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm)
	id, _ := services.CreateAnnotation(vs,
		entities.TimeSpan{Start: _8am, End: _9am},
		nil, "scott@ausocean.org",
		map[string]string{"species": "Sepia Apama"})

	annotation, err := services.GetAnnotationByID(int64(id))
	if err != nil {
		t.Errorf("Could not get annotation entity %s", err)
	}
	if annotation.VideoStreamID != vs || !annotation.TimeSpan.Start.Equal(_8am) || !annotation.TimeSpan.End.Equal(_9am) || annotation.BoundingBox != nil || annotation.Observer != "scott@ausocean.org" || annotation.ObservationKeys[0] != "species" || annotation.ObservationPairs[0] != "species:Sepia Apama" {
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
	cs, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	vs, _ := services.CreateVideoStream("http://youtube.com/watch?v=abc123", int64(cs), _8am, &_4pm)
	id, _ := services.CreateAnnotation(vs,
		entities.TimeSpan{Start: _8am, End: _9am},
		nil, "scott@ausocean.org",
		map[string]string{"species": "Sepia Apama"})

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
