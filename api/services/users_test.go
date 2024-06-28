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
	"github.com/ausocean/openfish/api/types/user"
)

func TestCreateUser(t *testing.T) {
	setup()

	// Create a new user entity.
	err := services.CreateUser("test@test.com", user.DefaultRole)
	if err != nil {
		t.Errorf("Could not create user entity %s", err)
	}
}

func TestUserExists(t *testing.T) {
	setup()

	// Create a new user entity.
	services.CreateUser("test@test.com", user.DefaultRole)

	// Check if the user exists.
	if !services.UserExists("test@test.com") {
		t.Errorf("Expected user to exist")
	}
}

func TestUserExistsForNonexistentEntity(t *testing.T) {
	setup()

	// Check if the user exists.
	// We expect it to return false.
	if services.UserExists("nonexistent@test.com") {
		t.Errorf("Did not expect user to exist")
	}
}

func TestGetUserByEmail(t *testing.T) {
	setup()

	// Define test cases.
	testCases := []struct {
		email string
		role  user.Role
	}{
		{"admin@test.com", user.AdminRole},
		{"curator@test.com", user.CuratorRole},
		{"annotator@test.com", user.AnnotatorRole},
		{"readonly@test.com", user.ReadonlyRole},
	}

	for _, tc := range testCases {
		// Create user entities for each test case.
		services.CreateUser(tc.email, tc.role)

		// Check if the user can be fetched and is the same.
		user, err := services.GetUserByEmail(tc.email)
		if err != nil {
			t.Errorf("Could not get user entity %s", err)
		}
		if user.Email != tc.email || user.Role != tc.role {
			t.Errorf("User entity does not match created entity")
		}
	}
}

func TestGetUserByEmailForNonexistentEntity(t *testing.T) {
	setup()

	user, err := services.GetUserByEmail("nonexistent@test.com")
	if user != nil && err == nil {
		t.Errorf("GetUserByEmail returned non-existing entity %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	setup()

	// Create a new user entity.
	services.CreateUser("test@test.com", user.DefaultRole)

	// Update the role.
	role := user.AdminRole
	err := services.UpdateUser("test@test.com", role)
	if err != nil {
		t.Errorf("Could not update user entity %s", err)
	}

	user, _ := services.GetUserByEmail("test@test.com")
	if user.Role != role {
		t.Errorf("Role did not update, expected %s, actual %s", role.String(), user.Role.String())
	}
}

func TestUpdateUserForNonExistentEntity(t *testing.T) {
	setup()

	err := services.UpdateUser("nonexistent@test.com", user.AdminRole)
	if err == nil {
		t.Errorf("Did not receive expected error when updating non-existent user")
	}
}

func TestDeleteUser(t *testing.T) {
	setup()

	// Create a new user entity.
	services.CreateUser("test@test.com", user.DefaultRole)

	// Delete the capture source entity.
	err := services.DeleteUser("test@test.com")
	if err != nil {
		t.Errorf("Could not delete user entity")
	}

	// Check if the capture source exists.
	if services.UserExists("test@test.com") {
		t.Errorf("User entity exists after delete")
	}
}

func TestDeleteUserForNonexistentEntity(t *testing.T) {
	setup()

	err := services.DeleteUser("nonexistent@test.com")
	if err == nil {
		t.Errorf("Did not receive expected error when deleting non-existent capture source")
	}
}
