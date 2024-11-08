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

package services

import (
	"context"
	"fmt"

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/types/mediatype"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
	"github.com/ausocean/openfish/datastore"
)

// Media is a saved video or image.
type Media struct {
	ID int64
	MediaContents
}

// MediaContents represents Media without an ID.
type MediaContents struct {
	Type              mediatype.MediaType
	VideoStreamSource int64
	StartTime         videotime.VideoTime
	EndTime           *videotime.VideoTime
	Bytes             []byte
}

func (m *MediaContents) ToEntity() entities.Media {
	e := entities.Media{
		Type:              int(m.Type),
		VideoStreamSource: m.VideoStreamSource,
		StartTime:         m.StartTime.Int(),
		EndTime:           nil,
		Bytes:             m.Bytes,
	}

	if m.EndTime != nil {
		end := m.EndTime.Int()
		e.EndTime = &end
	}

	return e
}

func MediaContentsFromEntity(e entities.Media) MediaContents {
	m := MediaContents{
		Type:              mediatype.MediaType(e.Type),
		VideoStreamSource: e.VideoStreamSource,
		StartTime:         videotime.FromInt(e.StartTime),
		Bytes:             e.Bytes,
	}

	if e.EndTime != nil {
		endTime := videotime.FromInt(*e.EndTime)
		m.EndTime = &endTime
	}

	return m
}

// GetMediaByID gets an image or video when provided with an ID.
func GetMediaByID(id int64) (*Media, error) {
	store := ds_client.Get()
	key := store.IDKey(entities.MEDIA_KIND, id)
	var entity entities.Media
	err := store.Get(context.Background(), key, &entity)
	if err != nil {
		return nil, err
	}

	media := Media{
		ID:            key.ID,
		MediaContents: MediaContentsFromEntity(entity),
	}
	return &media, nil
}

// GetMediaByTypeStreamAndTime gets the media with the specified type, source video stream and time.
// It can do so because media is uniquely identified by the combination of these three parameters.
func GetMediaByTypeStreamAndTime(mtype mediatype.MediaType, source int64, start videotime.VideoTime, end *videotime.VideoTime) (*Media, error) {

	if mtype.IsVideo() == (end == nil) {
		return nil, fmt.Errorf("video media types must be provided with an end time and image media types must not")
	}

	store := ds_client.Get()
	query := store.NewQuery(entities.MEDIA_KIND, false)
	query.FilterField("Type", "=", int(mtype))
	query.FilterField("VideoStreamSource", "=", source)
	query.FilterField("StartTime", "=", start.Int())
	if end != nil {
		query.FilterField("EndTime", "=", end.Int())
	}
	query.Limit(1)

	var entities []entities.Media
	keys, err := store.GetAll(context.Background(), query, &entities)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}

	media := Media{
		ID:            keys[0].ID,
		MediaContents: MediaContentsFromEntity(entities[0]),
	}
	return &media, nil
}

// MediaExists checks if the media exists in the datastore.
func MediaExists(id int64) bool {
	store := ds_client.Get()
	key := store.IDKey(entities.MEDIA_KIND, id)
	var media entities.Media
	err := store.Get(context.Background(), key, &media)
	return err == nil
}

// CreateMedia puts an image or video in the datastore.
func CreateMedia(media MediaContents) (int64, error) {

	// Verify VideoStream exists.
	if !VideoStreamExists(media.VideoStreamSource) {
		return 0, fmt.Errorf("video stream does not exist")
	}

	// TODO: Check media does not already exist.

	// Insert entity.
	store := ds_client.Get()
	entity := media.ToEntity()
	key := store.IncompleteKey(entities.MEDIA_KIND)
	key, err := store.Put(context.Background(), key, &entity)
	if err != nil {
		return 0, err
	}

	// Return ID of created media.
	return key.ID, nil
}

// DeleteMedia deletes an image/video.
func DeleteMedia(id int64) error {
	store := ds_client.Get()
	key := store.IDKey(entities.MEDIA_KIND, id)
	return store.Delete(context.Background(), key)
}
