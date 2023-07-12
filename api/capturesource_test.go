package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func TestCreateCaptureSource(t *testing.T) {

	testcases := []struct {
		reqBody   string
		resStatus int
	}{
		// TODO: Test cases with missing fields cause the server to crash.
		//       Add test case for this and fix issue.
		{
			resStatus: fiber.StatusOK,
			reqBody: `{
				"name": "Ochre Point",
				"location": "-35.221389,138.46333333",
				"camera_hardware": "rpi cam v2 - wide angle lens",
				"site_id": 10192840284
			}`,
		},
		{
			resStatus: fiber.StatusOK,
			reqBody: `{
				"name": "Ochre Point",
				"location": "-35.221389,138.46333333",
				"camera_hardware": "rpi cam v2 - wide angle lens"
			}`,
		},
		{
			resStatus: fiber.StatusOK,
			reqBody: `{
				"name": "Ochre Point",
				"location": "-35.221389,138.46333333"
			}`,
		},
		{
			resStatus: fiber.StatusBadRequest,
			reqBody: `{
				"name": true,
				"location": [1,2,3],
				"camera_hardware": 123456,
				"site_id": false,
				"extra": "fields"
			}`,
		},
		{
			resStatus: fiber.StatusBadRequest,
			reqBody:   "something that is not json",
		},
		{
			resStatus: fiber.StatusBadRequest,
			reqBody:   "",
		},
	}

	// Initialize app.
	ds_client.Init(false)
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	registerAPIRoutes(app)

	for i, testcase := range testcases {
		// Send test request.
		req := httptest.NewRequest("POST", "/api/v1/capturesources", strings.NewReader(testcase.reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		// Check request, response.
		if err != nil {
			t.Errorf("[test %d]  %s", i, err.Error())
		}

		if resp.StatusCode != testcase.resStatus {
			t.Errorf("[test %d]  expected HTTP code %d, got %d", i, testcase.resStatus, resp.StatusCode)
		}
	}
}

func TestGetCaptureSourcesWithLimitAndOffset(t *testing.T) {

	// Initialize app.
	ds_client.Init(false)
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	registerAPIRoutes(app)

	// Populate datastore with test data.
	for _, body := range []string{
		`{
			"name": "Ochre Point",
			"location": "-35.221389,138.46333333",
			"camera_hardware": "rpi cam v2 - wide angle lens",
			"site_id": 10192840284
		}`,
		`{
			"name": "Carrickalinga",
			"location": "-35.4270000,138.3148330",
			"camera_hardware": "rpi cam v2"
		}`,
		`{
			"name": "Windara",
			"location": "-35.221389,138.46333333",
			"camera_hardware": "rpi cam v2 - wide angle lens",
			"site_id": 12345678910
		}`,
	} {
		// Send test request.
		req := httptest.NewRequest("POST", "/api/v1/capturesources", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req, -1)
	}

	testcases := []struct {
		limit  int
		offset int
	}{
		{
			limit: 0, offset: 0,
		},
		{
			limit: 1, offset: 0,
		},
		{
			limit: 2, offset: 0,
		},
		{
			limit: 1, offset: 1,
		},
		{
			limit: 2, offset: 1,
		},
	}

	for i, testcase := range testcases {
		// Send request
		url := fmt.Sprintf("/api/v1/capturesources?limit=%d&offset=%d", testcase.limit, testcase.offset)
		req := httptest.NewRequest("GET", url, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("[test %d]  %s", i, err.Error())
		}

		// Check response.
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("[test %d]  expected HTTP code 200, got %d", i, resp.StatusCode)
		}

		// Unmarshal response.
		respBody, _ := ioutil.ReadAll(resp.Body)
		var result api.Result[handlers.CaptureSourceResult]
		if json.Unmarshal(respBody, &result) != nil {
			t.Errorf("[test %d]  invalid JSON format %s", i, string(respBody))
		}

		// Check response has correct limit and offset.
		if len(result.Results) != testcase.limit {
			t.Errorf("[test %d]  expected %d result, got %d results", i, testcase.limit, len(result.Results))
		}

		if result.Limit != testcase.limit {
			t.Errorf("[test %d]  expected limit to be equal to %d, got limit = %d", i, testcase.limit, result.Limit)
		}

		if result.Offset != testcase.offset {
			t.Errorf("[test %d]  expected offset to be equal to %d, got offset = %d", i, testcase.offset, result.Offset)
		}
	}
}
