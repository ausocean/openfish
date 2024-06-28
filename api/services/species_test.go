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

	"github.com/ausocean/openfish/api/services"
	"github.com/ausocean/openfish/api/types/species"
)

func TestCreateSpecies(t *testing.T) {
	setup()

	// Create a new species entity.
	_, err := services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]species.Image, 0), nil)
	if err != nil {
		t.Errorf("Could not create species entity %s", err)
	}
}

func TestSpeciesExists(t *testing.T) {
	setup()

	// Create a new species entity.
	id, _ := services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]species.Image, 0), nil)

	// Check if the species exists.
	if !services.SpeciesExists(int64(id)) {
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
	id, _ := services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]species.Image, 0), nil)

	species, err := services.GetSpeciesByID(int64(id))
	if err != nil {
		t.Errorf("Could not get species entity %s", err)
	}
	if species.CommonName != "Southern Reef Squid" && species.Species != "Sepioteuthis australis" && len(species.Images) != 0 {
		t.Errorf("Video stream entity does not match created entity")
	}
}

func TestGetSpeciesByIDForNonexistentEntity(t *testing.T) {
	setup()

	species, err := services.GetSpeciesByID(int64(123456789))
	if species != nil && err == nil {
		t.Errorf("GetSpeciesByID returned non-existing entity %s", err)
	}
}

// TODO: Write tests for GetRecommendedSpecies. Test limit, offset, and sorting.

func TestDeleteSpecies(t *testing.T) {
	setup()

	// Create a new species entity.
	id, _ := services.CreateSpecies("Sepioteuthis australis", "Southern Reef Squid", make([]species.Image, 0), nil)

	// Delete the species entity.
	err := services.DeleteSpecies(int64(id))
	if err != nil {
		t.Errorf("Could not delete species entity %d: %s", id, err)
	}

	// Check if the species exists.
	if services.SpeciesExists(int64(id)) {
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
