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

package sliceutils

import "iter"

// WindowPermutations creates an iterator of all the permutations of overlapping slices of a slice.
// For example: [1, 2, 3] will yield [1], [2], [3], [1, 2], [2, 3], [1, 2, 3].
func WindowPermutations[Slice ~[]E, E any](s Slice) iter.Seq[Slice] {
	return func(yield func(Slice) bool) {
		for i := 1; i <= len(s); i++ {
			for subslice := range Window(s, i) {
				yield(subslice)
			}
		}
	}
}

// Window creates an iterator of over overlapping slices of length n.
// Very similar to https://pkg.go.dev/slices#Chunk and std::slice::Windows in rust.
func Window[Slice ~[]E, E any](s Slice, n int) iter.Seq[Slice] {
	if n < 1 || n > len(s) {
		panic("cannot be less than 1 or greater than len(s)")
	}

	return func(yield func(Slice) bool) {
		for i := 0; i < len(s)-n+1; i += 1 {
			yield(s[i : i+n])
		}
	}
}
