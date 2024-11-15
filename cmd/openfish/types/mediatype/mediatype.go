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

// mediatype contains a mime-type enum for the types of media that can be downloaded
// to use as training data.
package mediatype

// MediaType is a mime-type enum for the types of media that can be downloaded to
// use as training data.
type MediaType uint8

// Accepted mime-types of media that can be downloaded to use as training data.
const (
	Invalid MediaType = iota
	JPEG
	MP4
)

// AllMimeTypes returns a list of all accepted mime types of media that can be
// downloaded to use as training data.
func AllMimeTypes() []string {
	return []string{JPEG.MimeType(), MP4.MimeType()}
}

// FromMimeType returns the MediaType for a given mime-type string.
func FromMimeType(s string) MediaType {
	switch s {
	case "image/jpeg":
		return JPEG
	case "video/mp4":
		return MP4
	}
	return Invalid
}

// FileExtension returns the file extension for a given MediaType.
func (t MediaType) FileExtension() string {
	switch t {
	case JPEG:
		return "jpeg"
	case MP4:
		return "mp4"
	}
	panic("unreachable")
}

// MimeType returns the mime-type for a given MediaType.
func (t MediaType) MimeType() string {
	switch t {
	case JPEG:
		return "image/jpeg"
	case MP4:
		return "video/mp4"
	}
	panic("unreachable")
}

// IsVideo returns true if the MediaType is a video.
func (t MediaType) IsVideo() bool {
	switch t {
	case MP4:
		return true
	}
	return false
}
