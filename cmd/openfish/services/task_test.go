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
	"net/url"
	"testing"

	"github.com/ausocean/openfish/cmd/openfish/services"
)

func TestCreateTask(t *testing.T) {
	setup()

	// Create a new task entity.

	_, err := services.CreateTask()

	if err != nil {
		t.Errorf("Could not create task entity %s", err)
	}
}

func TestGetTaskById(t *testing.T) {
	setup()

	id, _ := services.CreateTask()

	task, err := services.GetTaskById(id)
	if err != nil {
		t.Errorf("Could not get task entity %s", err)
	}

	if task.Resource != nil || task.Status != services.Pending {
		t.Errorf("Task entity does not match expected")
	}
}

func TestCancelTask(t *testing.T) {
	setup()

	id, _ := services.CreateTask()

	err := services.CancelTask(id)
	if err != nil {
		t.Errorf("Could not cancel task entity %s", err)
	}

	task, _ := services.GetTaskById(id)
	if task.Resource != nil || task.Status != services.Cancelled {
		t.Errorf("Task entity does not match expected")
	}
}

func TestFailTask(t *testing.T) {
	setup()

	id, _ := services.CreateTask()

	err := services.FailTask(id)
	if err != nil {
		t.Errorf("Could not fail task entity %s", err)
	}

	task, _ := services.GetTaskById(id)
	if task.Resource != nil || task.Status != services.Failed {
		t.Errorf("Task entity does not match expected")
	}
}

func TestCompleteTask(t *testing.T) {
	setup()

	id, _ := services.CreateTask()

	url, _ := url.Parse("http://openfish.appspot.com/api/v1/media/123")

	err := services.CompleteTask(id, url)
	if err != nil {
		t.Errorf("Could not complete task %s", err)
	}

	task, err := services.GetTaskById(id)
	if err != nil {
		t.Errorf("Could not get task %s", err)
	}
	if task.Resource.String() != url.String() || task.Status != services.Complete {
		t.Errorf("Task entity does not match expected")
	}
}
