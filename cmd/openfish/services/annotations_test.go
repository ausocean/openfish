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
	"reflect"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

func createTestAnnotation() services.Annotation {
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	sp := createTestSpecies()
	vs := createTestVideoStream()
	contents := services.AnnotationContents{
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
		Identifications: map[int64][]int64{
			sp.ID: {uid},
		},
		CreatedByID: uid,
	}
	created, _ := services.CreateAnnotation(contents)

	return *created
}

func TestCreateAnnotation(t *testing.T) {
	setup()

	// Create a new annotation entity.
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "coral.fischer@example.com",
		DisplayName: "Coral Fischer",
		Role:        role.Annotator,
	})
	sp := createTestSpecies()
	vs := createTestVideoStream()
	_, err := services.CreateAnnotation(services.AnnotationContents{
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
		Identifications: map[int64][]int64{
			sp.ID: {uid},
		},
		CreatedByID: uid,
	})

	if err != nil {
		t.Errorf("Could not create annotation entity %s", err)
	}
}

func TestAnnotationExists(t *testing.T) {
	setup()

	// Create a new annotation entity.
	annotation := createTestAnnotation()

	// Check if the annotation exists.
	if !services.AnnotationExists(annotation.ID) {
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

	// Create a new expected entity.
	expected := createTestAnnotation()

	actual, err := services.GetAnnotationByID(expected.ID)
	if err != nil {
		t.Errorf("Could not get annotation entity %s", err)
	}
	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("Annotation does not match created, expected %v, got %v", expected, *actual)
	}
}

func TestGetAnnotationByIDForNonexistentEntity(t *testing.T) {
	setup()

	actual, err := services.GetAnnotationByID(int64(123456789))
	if actual != nil && err == nil {
		t.Errorf("GetAnnotationByID returned non-existing entity %s", err)
	}
}

// TODO: Write tests for GetAnnotations. Test limit, offset and filtering.

func TestAnnotationApplyJoin(t *testing.T) {
	setup()
	annotation := createTestAnnotation()

	_, err := annotation.JoinFields()
	if err != nil {
		t.Errorf("Could not join fields %s", err)
	}
}

func TestAnnotationAddIdentification(t *testing.T) {
	setup()
	original := createTestAnnotation()
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "sandy.whiting@example.com",
		DisplayName: "Sandy Whiting",
		Role:        role.Annotator,
	})
	sp, _ := services.CreateSpecies(services.SpeciesContents{
		ScientificName: "Rhincodon typus",
		CommonName:     "Rhincodon typus",
	})

	services.AddIdentification(original.ID, uid, sp.ID)

	modified, _ := services.GetAnnotationByID(original.ID)
	if len(modified.Identifications) != len(original.Identifications)+1 {
		t.Errorf("Expected an additional identification to be added")
	}
}

func TestAnnotationAddNonExistingSpeciesIdentification(t *testing.T) {
	// TODO: Run test in CI when issue is fixed.
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	setup()
	original := createTestAnnotation()
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "sandy.whiting@example.com",
		DisplayName: "Sandy Whiting",
		Role:        role.Annotator,
	})

	err := services.AddIdentification(original.ID, uid, int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when adding non-existent species as an identification")
	}
}

func TestAnnotationRemoveIdentification(t *testing.T) {
	setup()
	original := createTestAnnotation()
	uid, _ := services.CreateUser(services.UserContents{
		Email:       "sandy.whiting@example.com",
		DisplayName: "Sandy Whiting",
		Role:        role.Annotator,
	})
	sp, _ := services.CreateSpecies(services.SpeciesContents{
		ScientificName: "Rhincodon typus",
		CommonName:     "Rhincodon typus",
	})
	services.AddIdentification(original.ID, uid, sp.ID)
	services.DeleteIdentification(original.ID, uid, sp.ID)

	modified, _ := services.GetAnnotationByID(original.ID)
	if len(modified.Identifications) != len(original.Identifications) {
		t.Errorf("Expected identification to be removed")
	}
}

func TestDeleteAnnotation(t *testing.T) {
	setup()

	// Create a new annotation entity.
	created := createTestAnnotation()

	// Delete the annotation entity.
	err := services.DeleteAnnotation(created.ID)
	if err != nil {
		t.Errorf("Could not delete annotation entity %d: %s", created.ID, err)
	}

	// Check if the annotation exists.
	if services.AnnotationExists(created.ID) {
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
