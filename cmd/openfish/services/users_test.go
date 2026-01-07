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

	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
)

func TestCreateUser(t *testing.T) {
	setup()

	// Create a new user entity.
	_, err := services.CreateUser(services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	})
	if err != nil {
		t.Errorf("Could not create user entity %s", err)
	}
}

func TestUserExists(t *testing.T) {
	setup()

	// Create a new user entity.
	id, err := services.CreateUser(services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	})
	if err != nil {
		t.Errorf("Could not create user: %s", err)
	}

	// Check if the user exists.
	if !services.UserExists(id) {
		t.Errorf("Expected user to exist")
	}
}

func TestUserExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the user exists.
	// We expect it to return false.
	if services.UserExists(int64(1234567890)) {
		t.Errorf("Did not expect user to exist")
	}
}

func TestGetUserByID(t *testing.T) {
	setup()

	contents := services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	}
	id, err := services.CreateUser(contents)
	if err != nil {
		t.Errorf("Could not create user: %s", err)
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		t.Errorf("Could not get user entity %s", err)
	}

	if user.ID != id {
		t.Errorf("User ID does not match")
	}

	if user.DisplayName != contents.DisplayName {
		t.Errorf("User display name does not match")
	}

	if user.Role != contents.Role {
		t.Errorf("User role does not match")
	}
}

func TestGetUserByEmail(t *testing.T) {
	setup()

	contents := services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	}
	_, err := services.CreateUser(contents)
	if err != nil {
		t.Errorf("Could not create user: %s", err)
	}

	user, err := services.GetUserByEmail("coral.fischer@example.com")
	if err != nil {
		t.Errorf("Could not get user entity: %s", err)
	}

	if user.DisplayName != contents.DisplayName {
		t.Errorf("User display name does not match")
	}

	if user.Role != contents.Role {
		t.Errorf("User role does not match")
	}

	if user.Email != contents.Email {
		t.Errorf("User email does not match")
	}
}

func TestUpdateUser(t *testing.T) {
	setup()

	// Create a new user entity.
	id, err := services.CreateUser(services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	})
	if err != nil {
		t.Errorf("Could not create user: %s", err)
	}

	// Update the role.
	role := role.Admin
	err = services.UpdateUser(id, services.PartialUserContents{
		Role: &role,
	})
	if err != nil {
		t.Errorf("Could not update user entity %s", err)
	}

	user, _ := services.GetUserByID(id)
	if user.Role != role {
		t.Errorf("Role did not update, expected %s, actual %s", role.String(), user.Role.String())
	}

	// Update the display name.
	displayName := "Coral Fischer"
	err = services.UpdateUser(id, services.PartialUserContents{
		DisplayName: &displayName,
	})
	if err != nil {
		t.Errorf("Could not update user entity %s", err)
	}

	user, _ = services.GetUserByID(id)
	if user.DisplayName != displayName {
		t.Errorf("DisplayName did not update, expected %s, actual %s", displayName, user.DisplayName)
	}
}

func TestUpdateUserForNonExistentEntity(t *testing.T) {
	setup()

	role := role.Admin
	err := services.UpdateUser(int64(1234567890), services.PartialUserContents{
		Role: &role,
	})
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent user")
	}
}

func TestDeleteUser(t *testing.T) {
	setup()

	// Create a new user entity.
	id, err := services.CreateUser(services.UserContents{
		DisplayName: "Coral Fischer",
		Email:       "coral.fischer@example.com",
		Role:        role.Default,
	})
	if err != nil {
		t.Errorf("Could not create user: %s", err)
	}

	// Delete the capture source entity.
	err = services.DeleteUser(id)
	if err != nil {
		t.Errorf("Could not delete user entity")
	}

	// Check if the capture source exists.
	if services.UserExists(id) {
		t.Errorf("User entity exists after delete")
	}
}

func TestDeleteUserForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteUser(int64(1234567890))
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent capture source")
	}
}
