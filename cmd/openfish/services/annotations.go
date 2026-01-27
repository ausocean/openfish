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

package services

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// Identification is a species suggestion made by users.
type Identification struct {
	Species      SpeciesSummary `json:"species"`
	IdentifiedBy []PublicUser   `json:"identified_by"`
}

// Annotation is a bounding box added to a video with one or many identifications.
// Users can suggest additional identifications for this annotation.
type Annotation struct {
	ID int64
	AnnotationContents
}

// AnnotationContents is the contents of an annotation.
type AnnotationContents struct {
	KeyPoints       []keypoint.KeyPoint
	Identifications map[int64][]int64
	VideostreamID   int64
	CreatedByID     int64
}

// AnnotationWithJoins is an annotation with its foreign key fields joined with
// their respective entities.
type AnnotationWithJoins struct {
	ID              int64               `json:"id" example:"1234567890"`
	KeyPoints       []keypoint.KeyPoint `json:"keypoints"`
	Identifications []Identification    `json:"identifications"`
	Videostream     VideoStreamSummary  `json:"videostream"`
	CreatedBy       PublicUser          `json:"created_by"`
	Start           videotime.VideoTime `json:"start" swaggertype:"string" example:"01:56:05.500"`
	End             videotime.VideoTime `json:"end" swaggertype:"string" example:"01:56:05.500"`
	Duration        int64               `json:"duration" example:"15"`
}

// JoinFields joins the foreign key fields of an annotation with their respective entities.
func (a *Annotation) JoinFields() (*AnnotationWithJoins, error) {

	// Get video stream details.
	videostream, err := GetVideoStreamByID(a.VideostreamID)
	if err != nil {
		return nil, err // TODO: More informative message.
	}

	// Get user details.
	user, err := GetUserByID(a.CreatedByID)
	if err != nil {
		return nil, err // TODO: More informative message.
	}

	// Get identifications.
	identifications := make([]Identification, 0, len(a.Identifications))
	for speciesID, userIDs := range a.Identifications {
		species, err := GetSpeciesByID(speciesID)
		if err != nil {
			return nil, err
		}

		users := make([]PublicUser, 0, len(userIDs))
		for _, userID := range userIDs {
			user, err := GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			users = append(users, user.ToPublicUser())
		}
		identifications = append(identifications, Identification{
			Species:      species.ToSummary(),
			IdentifiedBy: users,
		})
	}

	return &AnnotationWithJoins{
		ID:              a.ID,
		KeyPoints:       a.KeyPoints,
		Videostream:     videostream.ToSummary(),
		Identifications: identifications,
		CreatedBy:       user.ToPublicUser(),
		Start:           a.KeyPoints[0].Time,
		End:             a.KeyPoints[len(a.KeyPoints)-1].Time,
		Duration:        a.KeyPoints[len(a.KeyPoints)-1].Time.Int() - a.KeyPoints[0].Time.Int(),
	}, nil
}

// ToEntity converts an AnnotationContents struct to an entities.Annotation struct.
func (a *AnnotationContents) ToEntity() entities.Annotation {

	// Convert keypoints into storable format.
	kp := make([]struct {
		Time string
		keypoint.BoundingBox
	}, len(a.KeyPoints))
	for i := range a.KeyPoints {
		kp[i] = struct {
			Time string
			keypoint.BoundingBox
		}{
			Time:        a.KeyPoints[i].Time.String(),
			BoundingBox: a.KeyPoints[i].BoundingBox,
		}
	}

	// Convert identifications map into storable format (two arrays).
	users := make([]int64, 0, len(a.Identifications))
	species := make([]int64, 0, len(a.Identifications))
	for speciesID, userIDs := range a.Identifications {
		for _, userID := range userIDs {
			users = append(users, userID)
			species = append(species, speciesID)
		}
	}

	return entities.Annotation{
		VideoStreamID:           a.VideostreamID,
		Start:                   a.KeyPoints[0].Time.Int(),
		CreatedBy:               a.CreatedByID,
		Keypoints:               kp,
		IdentificationUserID:    users,
		IdentificationSpeciesID: species,
	}
}

// AnnotationContentsFromEntity converts an entity to an AnnotationContents struct.
func AnnotationContentsFromEntity(e entities.Annotation) AnnotationContents {

	identifications := make(map[int64][]int64)
	for i := range e.IdentificationUserID {
		userID := e.IdentificationUserID[i]
		speciesID := e.IdentificationSpeciesID[i]
		if _, exists := identifications[speciesID]; exists {
			identifications[speciesID] = append(identifications[speciesID], userID)
		} else {
			identifications[speciesID] = []int64{userID}
		}
	}

	keypoints := make([]keypoint.KeyPoint, len(e.Keypoints))
	for i, k := range e.Keypoints {
		keypoints[i] = keypoint.KeyPoint{
			BoundingBox: k.BoundingBox,
			Time:        videotime.UncheckedParse(k.Time),
		}
	}

	return AnnotationContents{
		KeyPoints:       keypoints,
		Identifications: identifications,
		VideostreamID:   e.VideoStreamID,
		CreatedByID:     e.CreatedBy,
	}
}

// GetAnnotationByID returns an annotation by ID.
func GetAnnotationByID(id int64) (*Annotation, error) {

	store := globals.GetStore()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var e entities.Annotation
	err := store.Get(context.Background(), key, &e)
	if err != nil {
		return nil, err
	}

	annotation := Annotation{
		ID:                 id,
		AnnotationContents: AnnotationContentsFromEntity(e),
	}

	return &annotation, nil
}

// AnnotationExists checks if an annotation exists.
func AnnotationExists(id int64) bool {
	store := globals.GetStore()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var annotation entities.Annotation
	err := store.Get(context.Background(), key, &annotation)
	return err == nil
}

// GetAnnotations gets a list of annotations, filtering by videostream if specified.
func GetAnnotations(limit int, offset int, order *string, videostream *int64) ([]Annotation, error) {
	store := globals.GetStore()
	query := store.NewQuery(entities.ANNOTATION_KIND, false)

	// Apply filters.
	if videostream != nil {
		query.FilterField("VideoStreamID", "=", *videostream)
	}

	// Apply pagination and ordering.
	query.Limit(limit)
	query.Offset(offset)
	if order != nil {
		query.Order(*order)
	}

	// Fetch entities from the datastore.
	var ents []entities.Annotation
	keys, err := store.GetAll(context.Background(), query, &ents)
	if err != nil {
		return []Annotation{}, err
	}

	// Convert entities.
	annotations := make([]Annotation, len(ents))
	for i := range ents {
		annotations[i] = Annotation{
			ID:                 keys[i].ID,
			AnnotationContents: AnnotationContentsFromEntity(ents[i]),
		}
	}

	return annotations, nil
}

// CreateAnnotation creates a new annotation.
func CreateAnnotation(contents AnnotationContents) (*Annotation, error) {

	// Verify VideoStream exists.
	if !VideoStreamExists(contents.VideostreamID) {
		return nil, errors.New("VideoStream does not exist")
	}

	// Get a unique ID for the new annotation.
	store := globals.GetStore()
	key := store.IncompleteKey(entities.ANNOTATION_KIND)
	ent := contents.ToEntity()
	key, err := store.Put(context.Background(), key, &ent)
	if err != nil {
		return nil, err
	}

	// Return newly created annotation.
	created := Annotation{
		ID:                 key.ID,
		AnnotationContents: contents,
	}
	return &created, nil
}

// AddIdentification adds a new species identification to an annotation.
func AddIdentification(id int64, userID int64, speciesID int64) error {
	// Update data in the datastore.
	store := globals.GetStore()

	// Check if speciesID exists.
	if !SpeciesExists(speciesID) {
		return fmt.Errorf("species ID %d does not exist", speciesID)
	}

	// Proceed with adding identification.
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var annotation entities.Annotation

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		ent, ok := e.(*entities.Annotation)
		if ok {
			a := AnnotationContentsFromEntity(*ent)
			ids := a.Identifications[speciesID]
			// Add an identification only if the user hasn't already identified the species.
			if !slices.Contains(ids, userID) {
				a.Identifications[speciesID] = append(ids, userID)
			}
			*ent = a.ToEntity()
		}
	}, &annotation)
}

// DeleteIdentification removes a species identification from an annotation.
func DeleteIdentification(id int64, userID int64, speciesID int64) error {
	// Update data in the datastore.
	store := globals.GetStore()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var annotation entities.Annotation

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		ent, ok := e.(*entities.Annotation)
		if ok {
			a := AnnotationContentsFromEntity(*ent)
			// If there is only one identification and it is by the calling user, remove the species identification from the map.
			if len(a.Identifications[speciesID]) == 1 && a.Identifications[speciesID][0] == userID {
				delete(a.Identifications, speciesID)
			} else {
				for i, id := range a.Identifications[speciesID] {
					if id == userID {
						a.Identifications[speciesID] = append(a.Identifications[speciesID][:i], a.Identifications[speciesID][i+1:]...)
						break
					}
				}
			}
			*ent = a.ToEntity()
		}
	}, &annotation)
}

// DeleteAnnotation deletes an annotation.
func DeleteAnnotation(id int64) error {
	// Delete entity.
	store := globals.GetStore()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	return store.Delete(context.Background(), key)
}
