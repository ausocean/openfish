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

// role defines the user roles that control access to the API.
package role

import "fmt"

type Role int8

const (
	Readonly  Role = iota // Can only use GET APIs.
	Annotator             // Can create annotations.
	Curator               // Can create annotations and videostreams.
	Admin                 // Can do everything.
	Default   = Annotator
)

// String returns the string representation of a Role.
func (r Role) String() string {
	switch r {
	case Readonly:
		return "readonly"
	case Annotator:
		return "annotator"
	case Curator:
		return "curator"
	case Admin:
		return "admin"
	}
	return "unknown"
}

// Parse parses a string into a Role.
func Parse(s string) (Role, error) {
	switch s {
	case "readonly":
		return Readonly, nil
	case "annotator":
		return Annotator, nil
	case "curator":
		return Curator, nil
	case "admin":
		return Admin, nil
	}
	return Default, fmt.Errorf("invalid role provided: %s", s)
}

// UnmarshalText is used for decoding query params or JSON into a Role.
func (r *Role) UnmarshalText(text []byte) error {
	var err error
	*r, err = Parse(string(text))
	return err
}

// MarshalText is used for encoding a Role into JSON or query params.
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}
