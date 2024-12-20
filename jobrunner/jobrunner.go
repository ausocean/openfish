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

// Package jobrunner offers running docker containers as jobs in the background with
// multiple implementations.
//
//   - CloudRunJobRunner is a Google Cloud Run implementation.
//   - DockerJobRunner is a local docker based implementation.
package jobrunner

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
)

// JobRunner defines the jobrunner interface. It lets us run a job, with arguments.
type JobRunner interface {
	Run(name string, args ...string) error
}

// CloudRunJobRunner implements JobRunner for containers executed using Google Cloud-Run Jobs.
type CloudRunJobRunner struct {
	project  string
	location string
}

// DockerJobRunner implements JobRunner for containers executed using docker on your local machine.
type DockerJobRunner struct {
}

// NewCloudRunJobRunner creates a new CloudRunJobRunner.
// project is your project name in Google Cloud, and location is the data center
// it is running in, e.g. australia-southeast1
func NewCloudRunJobRunner(project string, location string) JobRunner {
	return &CloudRunJobRunner{project, location}
}

// NewDockerJobRunner creates a new DockerJobRunner.
func NewDockerJobRunner() JobRunner {
	return &DockerJobRunner{}
}

// Run executes the given job in the background. Run is non-blocking,
// so you need to handle getting its results in some other way, such
// as using the datastore and polling.
func (r *CloudRunJobRunner) Run(name string, args ...string) error {
	ctx := context.Background()
	c, err := run.NewJobsClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	container := []*runpb.RunJobRequest_Overrides_ContainerOverride{{
		Args: args,
	}}

	req := &runpb.RunJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", r.project, r.location, name),
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: container,
		},
	}
	_, err = c.RunJob(ctx, req)

	return err
}

// Run executes the given job in the background. Run is non-blocking,
// so you need to handle getting its results in some other way, such
// as using the datastore and polling.
func (r *DockerJobRunner) Run(name string, args ...string) error {

	googleAppCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	dockerArgs := []string{
		"run",
		"-e",
		"GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcloud.json",
		"-v",
		fmt.Sprintf("%s:/tmp/gcloud.json:z", googleAppCredentials),
		name,
	}
	dockerArgs = append(dockerArgs, args...)

	fmt.Println(dockerArgs)
	cmd := exec.Command("docker", dockerArgs...)
	go cmd.Run()

	return nil
}
