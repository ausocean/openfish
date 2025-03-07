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
	"reflect"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/services"
)

func createTestSpecies() services.Species {
	species, _ := services.CreateSpecies(services.SpeciesContents{
		ScientificName: "Sepioteuthis australis",
		CommonName:     "Southern Reef Squid",
		Images:         []services.SpeciesImage{},
	})

	return *species
}

func TestCreateSpecies(t *testing.T) {
	setup()

	// Create a new species entity.
	_, err := services.CreateSpecies(services.SpeciesContents{
		ScientificName: "Sepioteuthis australis",
		CommonName:     "Southern Reef Squid",
		Images:         []services.SpeciesImage{},
	})
	if err != nil {
		t.Errorf("Could not create species entity %s", err)
	}
}

func TestSpeciesExists(t *testing.T) {
	setup()

	// Create a new species entity.
	species := createTestSpecies()

	// Check if the species exists.
	if !services.SpeciesExists(species.ID) {
		t.Errorf("Expected species to exist")
	}
}

func TestSpeciesExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the species exists.
	// We expect it to return false.
	if services.SpeciesExists(int64(123456789)) {
		t.Errorf("Did not expect species to exist")
	}
}

func TestGetSpeciesByID(t *testing.T) {
	setup()

	// Create a new species entity.
	expected := createTestSpecies()

	found, err := services.GetSpeciesByID(expected.ID)
	if err != nil {
		t.Errorf("Could not get species entity %s", err)
	}
	if !reflect.DeepEqual(expected, *found) {
		t.Errorf("Found species does not match expected, expected: %+v, found: %+v", expected, *found)
	}
}

func TestGetSpeciesByIDForNonexistentEntity(t *testing.T) {
	setup()

	species, err := services.GetSpeciesByID(int64(123456789))
	if species != nil && err == nil {
		t.Errorf("GetSpeciesByID returned non-existing entity %s", err)
	}
}

func TestGetSpeciesByINaturalistTaxonID(t *testing.T) {
	// TODO: Run test in CI when issue is fixed.
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	setup()

	// Create a new species entity.
	taxonID := 314239870
	expected, _ := services.CreateSpecies(services.SpeciesContents{
		ScientificName:     "Sepioteuthis australis",
		CommonName:         "Southern Reef Squid",
		Images:             []services.SpeciesImage{},
		INaturalistTaxonID: &taxonID,
	})

	found, err := services.GetSpeciesByINaturalistID(taxonID)
	if err != nil {
		t.Errorf("Could not get species entity %s", err)
	}
	if found == nil {
		t.Errorf("Species entity was not found")
	}
	if !reflect.DeepEqual(expected, *found) {
		t.Errorf("Found species does not match expected, expected: %+v, found: %+v", expected, *found)
	}
}

// TODO: Write tests for GetSpecies. Test limit, offset, and sorting.

func TestDeleteSpecies(t *testing.T) {
	setup()

	// Create a new species entity.
	species := createTestSpecies()

	// Delete the species entity.
	err := services.DeleteSpecies(species.ID)
	if err != nil {
		t.Errorf("Could not delete species entity %d: %s", species.ID, err)
	}

	// Check if the species exists.
	if services.SpeciesExists(species.ID) {
		t.Errorf("Video stream entity exists after delete")
	}
}

func TestDeleteSpeciesForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteSpecies(int64(123456789))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent species")
	}
}
