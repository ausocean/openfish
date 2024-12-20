/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2024, The OpenFish Contributors.

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

// tasks/youtube-download is a thin wrapper around yt-dlp that writes the downloaded
// video to the datastore and updates its task status when it completes or fails.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/mediatype"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// downloadVideo downloads a video segment from the stream at the specified start and end video times.
func downloadVideo(url url.URL, mtype mediatype.MediaType, start videotime.VideoTime, end videotime.VideoTime) ([]byte, error) {
	// Download video using yt-dlp, writing to standard out.
	cmd := exec.Command("yt-dlp",
		"--download-sections",
		fmt.Sprintf("*%s-%s", start.String(), end.String()),
		"--force-keyframes-at-cuts",
		"-S",
		fmt.Sprintf("ext:%s", mtype.FileExtension()),
		"-o",
		"-",
		"-q",
		url.String())

	// Make pipe from std out.
	r, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	// Run command.
	cmd.Start()

	// Read all bytes to buffer.
	bytes, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return bytes, nil
}

// downloadFrame downloads a frame from the stream at the specified video time.
func downloadFrame(url url.URL, mtype mediatype.MediaType, start videotime.VideoTime) ([]byte, error) {
	// TODO: implement downloading a frame from a youtube stream.
	return []byte{}, nil
}

// fail ends the program and marks the task as failed.
func fail(taskID int64, err error) {
	print(err.Error())
	services.FailTask(taskID, err)
	os.Exit(1)
}

func main() {

	// Parse command line arguments into primitive types.
	taskID := flag.Int64("task-id", 0, "task ID used to track task status (optional)")
	videostreamID := flag.Int64("videostream-id", 0, "source video stream ID of stream (required)")
	urlStr := flag.String("url", "", "URL of video to download (required)")
	typeStr := flag.String("type", "", "Mime type of video to download (required)")
	startStr := flag.String("start", "", "Start time of video snippet (required)")
	endStr := flag.String("end", "", "End time of video snippet (optional)")

	flag.Parse()

	// Datastore initialisation.
	err := globals.InitStore(false)
	if err != nil {
		os.Exit(1)
	}
	if *taskID == 0 {
		os.Exit(1)
	}

	// Parse/validate flags.
	if *videostreamID == 0 {
		fail(*taskID, fmt.Errorf("Video stream ID must be specified"))
	}

	streamURL, err := url.Parse(*urlStr)
	if err != nil {
		fail(*taskID, err)
	}

	mtype, err := mediatype.ParseMimeType(*typeStr)
	if err != nil {
		fail(*taskID, err)
	}

	start, err := videotime.Parse(*startStr)
	if err != nil {
		fail(*taskID, err)
	}

	var end *videotime.VideoTime
	if mtype.IsVideo() {
		endVal, err := videotime.Parse(*endStr)
		if err != nil {
			fail(*taskID, err)
		}
		end = &endVal
	}

	// Storage initialisation.
	err = globals.InitStorage(false)
	if err != nil {
		fail(*taskID, err)
	}

	// Download media.
	var bytes []byte
	if mtype.IsVideo() {
		bytes, err = downloadVideo(*streamURL, mtype, start, *end)
	} else {
		fail(*taskID, fmt.Errorf("downloading images is unimplemented"))
		bytes, err = downloadFrame(*streamURL, mtype, start)
	}

	if err != nil {
		fail(*taskID, err)
	}

	// Write media to the datastore.
	mk := services.MediaKey{
		Type:          mtype,
		VideoStreamID: *videostreamID,
		StartTime:     start,
		EndTime:       end,
	}
	_, err = services.CreateMedia(mk, bytes)
	if err != nil {
		fail(*taskID, err)
	}

	// Mark task as completed.
	var time string
	if mtype.IsVideo() {
		time = fmt.Sprintf("%s-%s", start.String(), end.String())
	} else {
		time = fmt.Sprintf("%s", start.String())
	}
	createdResource, err := url.ParseRequestURI(fmt.Sprintf("/api/v1/videostreams/%d/media/%s?time=%s", *videostreamID, mtype.MimeType(), time))
	if err != nil {
		fail(*taskID, err)
	}

	err = services.CompleteTask(*taskID, createdResource)
	if err != nil {
		fail(*taskID, err)
	}
}
