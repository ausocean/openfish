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

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/services"
)

func setup() {
	ds_client.Init(true)

	// Create directories if they do not exist.
	_ = os.Mkdir("store/openfish/CaptureSource", os.ModePerm)
	_ = os.Mkdir("store/openfish/VideoStream", os.ModePerm)
	_ = os.Mkdir("store/openfish/Annotation", os.ModePerm)
	_ = os.Mkdir("store/openfish/Species", os.ModePerm)
	_ = os.Mkdir("store/openfish/User", os.ModePerm)
	_ = os.Mkdir("store/openfish/Task", os.ModePerm)
	_ = os.Mkdir("store/openfish/Media", os.ModePerm)
}

func TestCreateCaptureSource(t *testing.T) {
	setup()

	// Create a new capture source entity.
	_, err := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	if err != nil {
		t.Errorf("Could not create capture source entity %s", err)
	}
}

func TestCaptureSourceExists(t *testing.T) {
	setup()

	// Create a new capture source entity.
	id, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)

	// Check if the capture source exists.
	if !services.CaptureSourceExists(int64(id)) {
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
		lat            float64
		long           float64
		cameraHardware string
	}{
		{"Stony Point camera 1", 0.0, 0.0, "RPI camera"},
		{"CuttleCam", -100.0, 100.0, "wide angle camera"},
	}

	for _, tc := range testCases {
		// Create capture source entities for each test case.
		id, _ := services.CreateCaptureSource(tc.name, tc.lat, tc.long, tc.cameraHardware, nil)

		// Check if the capture sources can be fetched and is the same.
		captureSource, err := services.GetCaptureSourceByID(int64(id))
		if err != nil {
			t.Errorf("Could not get capture source entity %s", err)
		}
		if captureSource.Name != tc.name || captureSource.Location.Lat != tc.lat || captureSource.Location.Lng != tc.long || captureSource.CameraHardware != tc.cameraHardware {
			t.Errorf("Capture source entity does not match created entity")
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
	id, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)

	// Update the name.
	name := "new name"
	err := services.UpdateCaptureSource(int64(id), &name, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, _ := services.GetCaptureSourceByID(int64(id))
	if captureSource.Name != name {
		t.Errorf("Name did not update, expected %s, actual %s", name, captureSource.Name)
	}

	// Update latitude and longitude.
	lat, long := -37.0, 145.0
	err = services.UpdateCaptureSource(int64(id), nil, &lat, &long, nil, nil)
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, _ = services.GetCaptureSourceByID(int64(id))
	if captureSource.Location.Lat != lat || captureSource.Location.Lng != long {
		t.Errorf("Location did not update, expected %f %f, actual %f %f", lat, long, captureSource.Location.Lat, captureSource.Location.Lng)
	}

	// Update cameraHardware.
	cameraHardware := "RPI wide angle camera"
	err = services.UpdateCaptureSource(int64(id), nil, nil, nil, &cameraHardware, nil)
	if err != nil {
		t.Errorf("Could not update capture source entity %s", err)
	}

	captureSource, _ = services.GetCaptureSourceByID(int64(id))
	if captureSource.CameraHardware != cameraHardware {
		t.Errorf("CameraHardware did not update, expected %s, actual %s", cameraHardware, captureSource.CameraHardware)
	}

	// TODO: test updating siteID.
}

func TestUpdateCaptureSourceForNonExistentEntity(t *testing.T) {
	setup()

	err := services.UpdateCaptureSource(int64(123456789), nil, nil, nil, nil, nil)
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent capture source")
	}
}

func TestDeleteCaptureSource(t *testing.T) {
	setup()

	// Create a new capture source entity.
	id, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)

	// Delete the capture source entity.
	err := services.DeleteCaptureSource(int64(id))
	if err != nil {
		t.Errorf("Could not delete capture source entity %d: %s", id, err)
	}

	// Check if the capture source exists.
	if services.CaptureSourceExists(int64(id)) {
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
	id, _ := services.CreateCaptureSource("Stony Point camera 1", 0.0, 0.0, "RPI camera", nil)
	services.CreateVideoStream("http://youtube.com/watch?v=abc123", id, _8am, &_4pm, []int64{})

	err := services.DeleteCaptureSource(id)
	if err == nil {
		t.Errorf("Did not receive expected error when deleting capture source with associated video stream")
	}
}
