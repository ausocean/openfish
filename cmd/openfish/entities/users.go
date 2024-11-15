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

package entities

import (
	"fmt"

	"github.com/ausocean/openfish/datastore"
)

// Kind of entity to store / fetch from the datastore.
const USER_KIND = "User"

// User contains the user role and email address.
type User struct {
	Email string
	Role  Role
}

// Role enum.
type Role int8

const (
	ReadonlyRole  Role = iota // Can only use GET APIs.
	AnnotatorRole             // Can create annotations.
	CuratorRole               // Can create annotations and videostreams.
	AdminRole                 // Can do everything.
	DefaultRole   = AnnotatorRole
)

func (r Role) String() string {
	switch r {
	case ReadonlyRole:
		return "readonly"
	case AnnotatorRole:
		return "annotator"
	case CuratorRole:
		return "curator"
	case AdminRole:
		return "admin"
	}
	return "unknown"
}

func ParseRole(s string) (Role, error) {
	switch s {
	case "readonly":
		return ReadonlyRole, nil
	case "annotator":
		return AnnotatorRole, nil
	case "curator":
		return CuratorRole, nil
	case "admin":
		return AdminRole, nil
	}
	return DefaultRole, fmt.Errorf("invalid role provided: %s", s)
}

// Implements Copy from the Entity interface.
func (vs *User) Copy(dst datastore.Entity) (datastore.Entity, error) {
	var v *User
	if dst == nil {
		v = new(User)
	} else {
		var ok bool
		v, ok = dst.(*User)
		if !ok {
			return nil, datastore.ErrWrongType
		}
	}
	*v = *vs
	return v, nil
}

// GetCache returns nil, because no caching is used.
func (vs *User) GetCache() datastore.Cache {
	return nil
}

// NewUser returns a new User entity.
func NewUser() datastore.Entity {
	return &User{}
}
