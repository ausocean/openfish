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
	"os"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/latlong"
	"github.com/ausocean/openfish/cmd/openfish/types/timezone"
)

func setup() {
	globals.InitStore(true)
	globals.InitStorage(true)

	// Create directories if they do not exist.
	os.MkdirAll("store/openfish/CaptureSource", os.ModePerm)
	os.MkdirAll("store/openfish/VideoStream", os.ModePerm)
	os.MkdirAll("store/openfish/Annotation", os.ModePerm)
	os.MkdirAll("store/openfish/Species", os.ModePerm)
	os.MkdirAll("store/openfish/User", os.ModePerm)
	os.MkdirAll("store/openfish/Task", os.ModePerm)
	os.MkdirAll("openfish-media/images", os.ModePerm)
	os.MkdirAll("openfish-media/videos", os.ModePerm)
}

// createTestCaptureSource creates a capture source in the datastore for use in tests.
func createTestCaptureSource() services.CaptureSource {
	cs, _ := services.CreateCaptureSource(services.CaptureSourceContents{
		Name:           "Stony Point camera 1",
		Location:       latlong.UncheckedParse("-37.000,145.000"),
		CameraHardware: "RPI camera",
	})
	return *cs
}

func TestCreateCaptureSource(t *testing.T) {
	setup()

	// Create a new capture source entity.
	_, err := services.CreateCaptureSource(services.CaptureSourceContents{
		Name:           "Stony Point camera 1",
		Location:       latlong.UncheckedParse("-37.000,145.000"),
		CameraHardware: "RPI camera",
	})
	if err != nil {
		t.Errorf("Could not create capture source entity %s", err)
	}
}

func TestCaptureSourceExists(t *testing.T) {
	setup()

	// Create a new capture source entity.
	cs := createTestCaptureSource()

	// Check if the capture source exists.
	if !services.CaptureSourceExists(cs.ID) {
		t.Errorf("Expected capture source to exist")
	}
}

func TestCaptureSourceExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the capture source exists.
	// We expect it to return false.
	if services.CaptureSourceExists(int64(123456789)) {
		t.Errorf("Did not expect capture source to exist")
	}
}

func TestGetCaptureSourceByID(t *testing.T) {
	setup()

	// Define test cases.
	testCases := []struct {
		name           string
		location       string
		cameraHardware string
	}{
		{"Stony Point camera 1", "-37.00000000,145.00000000", "RPI camera"},
		{"CuttleCam", "45.00000000,-20.00000000", "wide angle camera"},
	}

	for _, tc := range testCases {
		// Create capture source entities for each test case.
		cs, err := services.CreateCaptureSource(services.CaptureSourceContents{
			Name:           tc.name,
			Location:       latlong.UncheckedParse(tc.location),
			CameraHardware: tc.cameraHardware,
		})
		if err != nil {
			t.Errorf("Could not create capture source entity %s", err)
		}

		// Check if the capture sources can be fetched and is the same.
		captureSource, err := services.GetCaptureSourceByID(cs.ID)
		if err != nil {
			t.Errorf("Could not get capture source entity %s", err)
		}
		if captureSource.CameraHardware != tc.cameraHardware || captureSource.Name != tc.name || captureSource.Location.String() != tc.location {
			t.Errorf("Capture source entity does not match created entity: expected: %v, actual %v", tc, captureSource)
		}
	}
}

func TestGetCaptureSourceByIDForNonexistentEntity(t *testing.T) {
	setup()

	captureSource, err := services.GetCaptureSourceByID(int64(123456789))
	if captureSource != nil && err == nil {
		t.Errorf("GetCaptureSourceByID returned non-existing entity %s", err)
	}
}

// TODO: Write tests for GetCaptureSources. Test limit, offset and filtering.

func TestUpdateCaptureSource(t *testing.T) {
	setup()

	// Create a new capture source entity.
	cs, err := services.CreateCaptureSource(services.CaptureSourceContents{
		Name:           "Stony Point camera 1",
		Location:       latlong.UncheckedParse("-37.000,145.000"),
		CameraHardware: "RPI camera",
	})
	if err != nil {
		t.Errorf("Could not create capture source entity %s", err)
	}

	// Update the name.
	name := "new name"
	err = services.UpdateCaptureSource(cs.ID, services.PartialCaptureSourceContents{
		Name: &name,
	})
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, err := services.GetCaptureSourceByID(cs.ID)
	if err != nil {
		t.Errorf("Could not get capture source entity %s", err)
	}
	if captureSource.Name != name {
		t.Errorf("Name did not update, expected %s, actual %s", name, captureSource.Name)
	}

	// Update latitude and longitude.
	location := latlong.UncheckedParse("-37.0,145.0")
	err = services.UpdateCaptureSource(cs.ID, services.PartialCaptureSourceContents{
		Location: &location,
	})
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, err = services.GetCaptureSourceByID(cs.ID)
	if err != nil {
		t.Errorf("Could not get capture source entity %s", err)
	}
	if captureSource.Location.String() != location.String() {
		t.Errorf("Location did not update, expected %s, actual %s", location, captureSource.Location)
	}

	// Update cameraHardware.
	cameraHardware := "RPI wide angle camera"
	err = services.UpdateCaptureSource(cs.ID, services.PartialCaptureSourceContents{
		CameraHardware: &cameraHardware,
	})
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, err = services.GetCaptureSourceByID(cs.ID)
	if err != nil {
		t.Errorf("Could not get capture source entity %s", err)
	}
	if captureSource.CameraHardware != cameraHardware {
		t.Errorf("CameraHardware did not update, expected %s, actual %s", cameraHardware, captureSource.CameraHardware)
	}

	// TODO: test updating siteID.
}

func TestUpdateCaptureSourceForNonExistentEntity(t *testing.T) {
	setup()

	err := services.UpdateCaptureSource(int64(123456789), services.PartialCaptureSourceContents{})
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent capture source")
	}
}

func TestDeleteCaptureSource(t *testing.T) {
	setup()

	// Create a new capture source entity.
	cs, err := services.CreateCaptureSource(services.CaptureSourceContents{
		Name:           "Stony Point camera 1",
		Location:       latlong.UncheckedParse("-37.0,145.0"),
		CameraHardware: "RPI camera",
	})
	if err != nil {
		t.Errorf("Could not create capture source entity %s", err)
	}

	// Delete the capture source entity.
	err = services.DeleteCaptureSource(cs.ID)
	if err != nil {
		t.Errorf("Could not delete capture source entity %d: %s", cs.ID, err)
	}

	// Check if the capture source exists.
	if services.CaptureSourceExists(cs.ID) {
		t.Errorf("Capture source entity exists after delete")
	}
}

func TestDeleteCaptureSourceForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteCaptureSource(int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent capture source")
	}
}

func TestDeleteCaptureSourceWithAssociatedVideoStreams(t *testing.T) {
	// TODO: Run test in CI when issue is fixed.
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	setup()

	// Create a new capture source entity and a video stream that references it.
	cs := createTestCaptureSource()
	services.CreateVideoStream(services.VideoStreamContents{
		StartTime:     _8am,
		EndTime:       &_4pm,
		AnnotatorList: []int64{},
		BaseVideoStreamFields: services.BaseVideoStreamFields{
			TimeZone:      timezone.UncheckedParse("Australia/Adelaide"),
			StreamURL:     "http://youtube.com/watch?v=abc123",
			CaptureSource: cs.ID,
		},
	})

	err := services.DeleteCaptureSource(cs.ID)
	if err == nil {
		t.Errorf("Did not receive expected error when deleting capture source with associated video stream")
	}
}
