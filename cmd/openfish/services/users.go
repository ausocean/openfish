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
	"context"

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
)

// PublicUser is a user without private information, for public display, e.g. on annotations.
type PublicUser struct {
	ID          int64     `json:"id" example:"1234567890"`
	DisplayName string    `json:"display_name" example:"Coral Fischer"`
	Role        role.Role `json:"role" swaggertype:"string" example:"annotator"`
}

// User is a complete user.
type User struct {
	ID int64 `json:"id" example:"1234567890"`
	UserContents
}

// UserContents is the contents of a user in private contexts.
type UserContents struct {
	Email       string    `json:"email" example:"coral.fischer@example.com"`
	DisplayName string    `json:"display_name" example:"Coral Fischer"`
	Role        role.Role `json:"role" swaggertype:"string" example:"annotator"`
}

// PartialUserContents is for updating a user with a partial update (such as a PATCH request).
type PartialUserContents struct {
	Email       *string    `json:"email" validate:"optional" example:"coral.fischer@example.com"`
	DisplayName *string    `json:"display_name" validate:"optional" example:"Coral Fischer"`
	Role        *role.Role `json:"role" validate:"optional" example:"annotator" swaggertype:"string" enums:"readonly,annotator,curator,admin"`
}

// UserContentsFromEntity converts an entities.User to a UserContents.
func UserContentsFromEntity(u entities.User) UserContents {
	return UserContents{
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Role:        u.Role,
	}
}

// ToEntity converts a UserContents to an entities.User for storage in the datastore.
func (u *UserContents) ToEntity() entities.User {
	e := entities.User{
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Role:        u.Role,
	}

	return e
}

// ToPublicUser converts a services.User into a services.PublicUser. This is used for hiding private information
// such as email addresses from other users when they make use of the API.
func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Role:        u.Role,
	}
}

// GetUserByID gets a user when provided when provided with an ID.
func GetUserByID(id int64) (*User, error) {
	store := globals.GetStore()
	key := store.IDKey(entities.USER_KIND, id)
	var user entities.User
	err := store.Get(context.Background(), key, &user)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		UserContents: UserContentsFromEntity(user),
	}, nil
}

// GetUserByEmail gets a user when provided with an email address.
func GetUserByEmail(email string) (*User, error) {
	store := globals.GetStore()
	query := store.NewQuery(entities.USER_KIND, false)

	query.FilterField("Email", "=", email)
	query.Limit(1)

	var users []entities.User
	keys, err := store.GetAll(context.Background(), query, &users)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 || len(users) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}

	return &User{
		ID:           keys[0].ID,
		UserContents: UserContentsFromEntity(users[0]),
	}, nil
}

// UserExists checks if a user exists for a given ID.
func UserExists(id int64) bool {
	store := globals.GetStore()
	key := store.IDKey(entities.USER_KIND, id)
	var user entities.User
	err := store.Get(context.Background(), key, &user)
	return err == nil
}

// GetUsers gets a list of users.
func GetUsers(limit int, offset int) ([]User, error) {

	// Fetch data from the datastore.
	store := globals.GetStore()
	query := store.NewQuery(entities.USER_KIND, false)

	query.Limit(limit)
	query.Offset(offset)

	var users []entities.User
	ids, err := store.GetAll(context.Background(), query, &users)
	if err != nil {
		return []User{}, err
	}

	results := make([]User, len(users))
	for i := range users {
		results[i] = User{
			ID:           ids[i].ID,
			UserContents: UserContentsFromEntity(users[i]),
		}
	}

	return results, nil
}

// CreateUser creates a new user.
func CreateUser(user UserContents) (int64, error) {

	store := globals.GetStore()
	key := store.IncompleteKey(entities.USER_KIND)

	u := user.ToEntity()

	// Add to datastore.
	id, err := store.Put(context.Background(), key, &u)
	if err != nil {
		return 0, err
	}

	return id.ID, nil
}

// UpdateUser updates a user's role.
func UpdateUser(id int64, updates PartialUserContents) error {

	// Update data in the datastore.
	store := globals.GetStore()
	key := store.IDKey(entities.USER_KIND, id)
	var user entities.User

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*entities.User)
		if ok {
			if updates.DisplayName != nil {
				v.DisplayName = *updates.DisplayName
			}
			if updates.Role != nil {
				v.Role = *updates.Role
			}
			if updates.Email != nil {
				v.Email = *updates.Email
			}
		}
	}, &user)
}

// DeleteUser deletes a user.
func DeleteUser(id int64) error {
	// Delete entity.
	store := globals.GetStore()
	key := store.IDKey(entities.USER_KIND, id)
	return store.Delete(context.Background(), key)
}
