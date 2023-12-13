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

// Package ds_client initializes the datastore and makes it available to other packages through the use of Get().
package ds_client

import (
	"context"

	"github.com/ausocean/openfish/api/model"
	"github.com/ausocean/openfish/datastore"
)

var store datastore.Store

// Get returns the datastore global variable.
func Get() datastore.Store {
	return store
}

// Init initializes the datastore global variable and datastore client.
func Init(local bool) {
	ctx := context.Background()
	var err error
	if local {
		store, err = datastore.NewStore(ctx, "file", "openfish", "./store")
	} else {
		store, err = datastore.NewStore(ctx, "cloud", "openfish-dev", "")
	}
	if err != nil {
		panic(err)
	}

	datastore.RegisterEntity(model.CAPTURESOURCE_KIND, model.NewCaptureSource)
	datastore.RegisterEntity(model.VIDEOSTREAM_KIND, model.NewVideoStream)
	datastore.RegisterEntity(model.ANNOTATION_KIND, model.NewAnnotation)
}
