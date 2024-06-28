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

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/datastore"
)

// GetUserByEmail gets a user when provided with an email.
func GetUserByEmail(email string) (*entities.User, error) {
	store := ds_client.Get()
	key := store.NameKey(entities.USER_KIND, email)
	var user entities.User
	err := store.Get(context.Background(), key, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserExists(email string) bool {
	store := ds_client.Get()
	key := store.NameKey(entities.USER_KIND, email)
	var user entities.User
	err := store.Get(context.Background(), key, &user)
	return err == nil
}

// GetUsers gets a list of users.
func GetUsers(limit int, offset int) ([]entities.User, error) {

	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery(entities.USER_KIND, false)

	query.Limit(limit)
	query.Offset(offset)

	var users []entities.User
	_, err := store.GetAll(context.Background(), query, &users)
	if err != nil {
		return []entities.User{}, err
	}

	return users, nil
}

// CreateUser creates a new user.
func CreateUser(email string, role entities.Role) error {

	// Use the user's email as a unique ID.
	store := ds_client.Get()
	key := store.NameKey(entities.USER_KIND, email)

	user := entities.User{
		Email: email,
		Role:  role,
	}

	// Add to datastore.
	_, err := store.Put(context.Background(), key, &user)
	return err
}

// UpdateUser updates a user's role.
func UpdateUser(email string, role entities.Role) error {

	// Update data in the datastore.
	store := ds_client.Get()
	key := store.NameKey(entities.USER_KIND, email)
	var user entities.User

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*entities.User)
		if ok {
			v.Role = role
		}
	}, &user)
}

// DeleteUser deletes a user.
func DeleteUser(email string) error {
	// Delete entity.
	store := ds_client.Get()
	key := store.NameKey(entities.USER_KIND, email)
	return store.Delete(context.Background(), key)
}
