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

// services contains the main logic for the OpenFish API.
package services

import (
	"context"
	"net/url"

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
)

// TaskStatus represents the status of a task.
type TaskStatus uint8

const (
	Pending TaskStatus = iota
	Complete
	Cancelled
	Failed
)

// String returns TaskStatus as a string.
func (t TaskStatus) String() string {
	switch t {
	case Pending:
		return "pending"
	case Complete:
		return "complete"
	case Cancelled:
		return "cancelled"
	case Failed:
		return "failed"
	}
	panic("unreachable")
}

// Task is used to track progress of asynchronous tasks and track completion status.
// A task is initially pending but can be marked as complete, cancelled or failed.
// A completed task may have a URL to a resource.
//
// State machine:
//
//	+-----------+      +-----------+
//	|  Pending  | ---> | Complete  |
//	+-----------+      +-----------+
//	      |            +-----------+
//	      |----------> | Cancelled |
//	      |            +-----------+
//	      |            +-----------+
//	      +----------> |  Failed   |
//	                   +-----------+
type Task struct {
	ID       int64
	Status   TaskStatus
	Resource *url.URL
	Error    string
}

// GetTaskByID gets a task when provided with an ID.
func GetTaskById(id int64) (*Task, error) {
	store := globals.GetStore()
	key := store.IDKey(entities.TASK_KIND, id)
	var t entities.Task
	err := store.Get(context.Background(), key, &t)
	if err != nil {
		return nil, err
	}
	var r *url.URL
	if t.Resource != "" {
		r, err = url.ParseRequestURI(t.Resource)
		if err != nil {
			return nil, err
		}
	}

	task := Task{
		ID:       key.ID,
		Status:   TaskStatus(t.Status),
		Resource: r,
		Error:    t.Error,
	}
	return &task, nil
}

// CreateTask creates a new pending task.
func CreateTask() (int64, error) {
	store := globals.GetStore()
	t := entities.Task{
		Status: int(Pending),
	}
	key := store.IncompleteKey(entities.TASK_KIND)
	key, err := store.Put(context.Background(), key, &t)
	if err != nil {
		return 0, err
	}

	return key.ID, nil
}

// CancelTask marks a task as cancelled.
func CancelTask(id int64) error {
	store := globals.GetStore()
	key := store.IDKey(entities.TASK_KIND, id)

	var ta entities.Task
	return store.Update(context.Background(), key, func(e datastore.Entity) {
		t, ok := e.(*entities.Task)
		if ok && t.Status == int(Pending) {
			t.Status = int(Cancelled)
		}
	}, &ta)
}

// FailTask marks a task as failed and optionally attaches an error message to it.
func FailTask(id int64, err error) error {
	store := globals.GetStore()
	key := store.IDKey(entities.TASK_KIND, id)

	var ta entities.Task
	return store.Update(context.Background(), key, func(e datastore.Entity) {
		t, ok := e.(*entities.Task)
		if ok && t.Status == int(Pending) {
			t.Status = int(Failed)
		}
		if err != nil {
			t.Error = err.Error()
		}
	}, &ta)
}

// CompleteTask marks a task as complete and optionally attaches a resource URL to it.
func CompleteTask(id int64, resource *url.URL) error {
	store := globals.GetStore()
	key := store.IDKey(entities.TASK_KIND, id)

	var ta entities.Task
	return store.Update(context.Background(), key, func(e datastore.Entity) {
		t, ok := e.(*entities.Task)
		if ok {
			t.Status = int(Complete)
			if resource != nil {
				t.Resource = resource.RequestURI()
			}
		}
	}, &ta)
}
