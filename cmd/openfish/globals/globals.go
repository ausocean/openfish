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

// Package globals makes various variables available to other packages through the use of Get*().
package globals

import (
	"context"

	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/datastore"
	"github.com/ausocean/openfish/jobrunner"
)

var store datastore.Store
var runner jobrunner.JobRunner

// GetStore returns the datastore global variable.
func GetStore() datastore.Store {
	return store
}

// InitStore initializes the datastore global variable and datastore client.
func InitStore(local bool) {
	ctx := context.Background()
	var err error
	if local {
		store, err = datastore.NewStore(ctx, "file", "openfish", "./store")
	} else {
		store, err = datastore.NewStore(ctx, "cloud", "openfish", "")
	}
	if err != nil {
		panic(err)
	}

	datastore.RegisterEntity(entities.CAPTURESOURCE_KIND, entities.NewCaptureSource)
	datastore.RegisterEntity(entities.VIDEOSTREAM_KIND, entities.NewVideoStream)
	datastore.RegisterEntity(entities.ANNOTATION_KIND, entities.NewAnnotation)
	datastore.RegisterEntity(entities.SPECIES_KIND, entities.NewSpecies)
	datastore.RegisterEntity(entities.USER_KIND, entities.NewUser)
}

// GetRunner returns the job runner global variable.
func GetRunner() jobrunner.JobRunner {
	return runner
}

// InitRunner initializes the job runner global variable.
func InitRunner(local bool) {
	if local {
		runner = jobrunner.NewDockerJobRunner()
	} else {
		runner = jobrunner.NewCloudRunJobRunner("openfish", "southeast-1")
	}
}

// GetBucket returns the bucket global variable and bucket API client.
func GetBucket() {
	// TODO: implement.
}

// InitBucket initializes the bucket API client.
func InitBucket(local bool) {
	// TODO: implement.
}
