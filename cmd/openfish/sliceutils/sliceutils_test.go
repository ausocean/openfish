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

package sliceutils_test

import (
	"reflect"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/sliceutils"
)

func TestWindow(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	expected := [][]int{[]int{1, 2}, []int{2, 3}, []int{3, 4}, []int{4, 5}}

	var output [][]int

	for v := range sliceutils.Window(slice, 2) {
		output = append(output, v)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v but got %v", expected, output)
	}
}

func TestWindowPermutations(t *testing.T) {
	slice := []int{1, 2, 3}

	expected := [][]int{
		[]int{1},
		[]int{2},
		[]int{3},
		[]int{1, 2},
		[]int{2, 3},
		[]int{1, 2, 3},
	}

	var output [][]int

	for v := range sliceutils.WindowPermutations(slice) {
		output = append(output, v)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v but got %v", expected, output)
	}
}
