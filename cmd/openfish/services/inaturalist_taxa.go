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

// services contains the main logic for the OpenFish API.
package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PaginatedAPIResponse struct {
	TotalResults int    `json:"total_results"`
	Page         int    `json:"page"`
	PerPage      int    `json:"per_page"`
	Results      []Taxa `json:"results"`
}

type Taxa struct {
	ID                        int        `json:"id"`
	Rank                      string     `json:"rank"`
	RankLevel                 int        `json:"rank_level"`
	IconicTaxonID             int        `json:"iconic_taxon_id"`
	AncestorIDS               []int      `json:"ancestor_ids"`
	IsActive                  bool       `json:"is_active"`
	Name                      string     `json:"name"`
	ParentID                  int        `json:"parent_id"`
	Ancestry                  string     `json:"ancestry"`
	Names                     []Name     `json:"names"`
	Extinct                   bool       `json:"extinct"`
	DefaultPhoto              *Photo     `json:"default_photo"`
	TaxonChangesCount         int        `json:"taxon_changes_count"`
	TaxonSchemesCount         int        `json:"taxon_schemes_count"`
	ObservationsCount         int        `json:"observations_count"`
	FlagCounts                FlagCounts `json:"flag_counts"`
	CurrentSynonymousTaxonIDS any        `json:"current_synonymous_taxon_ids"`
	AtlasID                   any        `json:"atlas_id"`
	CompleteSpeciesCount      any        `json:"complete_species_count"`
	WikipediaURL              string     `json:"wikipedia_url"`
	CompleteRank              string     `json:"complete_rank"`
	IconicTaxonName           string     `json:"iconic_taxon_name"`
	PreferredCommonName       string     `json:"preferred_common_name"`
}

type Photo struct {
	ID                 int             `json:"id"`
	LicenseCode        string          `json:"license_code"`
	Attribution        string          `json:"attribution"`
	URL                string          `json:"url"`
	OriginalDimensions ImageDimensions `json:"original_dimensions"`
	Flags              []any           `json:"flags"`
	SquareURL          string          `json:"square_url"`
	MediumURL          string          `json:"medium_url"`
}

type ImageDimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type FlagCounts struct {
	Resolved   int `json:"resolved"`
	Unresolved int `json:"unresolved"`
}

type Name struct {
	Name     string `json:"name"`
	Locale   string `json:"locale"`
	Lexicon  string `json:"lexicon"`
	Position int    `json:"position"`
	IsValid  bool   `json:"is_valid"`
}

const inatUserAgent = "openfish" // As recommended by: https://www.inaturalist.org/pages/api+recommended+practices

func GetSpeciesByDescendant(parentID int) ([]Taxa, error) {
	cursor := 0
	var species []Taxa

	for {
		// Get a page of results.
		url := fmt.Sprintf("https://api.inaturalist.org/v1/taxa?is_active=true&rank=species&taxon_id=%d&order_by=id&order=asc&id_above=%d", parentID, cursor)
		agent := fiber.Get(url)
		agent.UserAgent(inatUserAgent)

		var res PaginatedAPIResponse
		code, _, errs := agent.Struct(&res)

		if errs != nil {
			return nil, errors.Join(errs...)
		}

		if code != 200 {
			return nil, fmt.Errorf("iNaturalist API returned status code %d", code)
		}

		// If no more results available, return.
		if len(res.Results) == 0 {
			return species, nil
		}

		// Set slice capacity to the total results as reported by API.
		if species == nil {
			species = make([]Taxa, 0, res.TotalResults)
		}

		// Set cursor position to be id of last taxa in response.
		cursor = res.Results[len(res.Results)-1].ID

		// Append results to species.
		species = append(species, res.Results...)
	}
}

func GetTaxonByName(taxonName string) (*Taxa, error) {
	name := strings.Split(taxonName, " ")
	url := fmt.Sprintf("https://api.inaturalist.org/v1/taxa?rank=%s&q=%s", name[0], name[1])
	agent := fiber.Get(url)
	agent.UserAgent(inatUserAgent)

	var res PaginatedAPIResponse
	code, _, errs := agent.Struct(&res)

	if code != 200 {
		return nil, fmt.Errorf("iNaturalist API returned status code %d", code)
	}

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	if len(res.Results) == 0 {
		return nil, nil
	}

	return &res.Results[0], nil
}
