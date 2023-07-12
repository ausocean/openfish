package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ausocean/openfish/api/ds_client"
	"github.com/gofiber/fiber/v2"
)

func setup(t *testing.T) (*fiber.App, int64) {
	// Initialize app.
	ds_client.Init(false)
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	registerAPIRoutes(app)

	// Populate datastore with capturesource for testing.
	// Send test request.
	body := `{
		"name": "Ochre Point",
		"location": "-35.221389,138.46333333",
		"camera_hardware": "rpi cam v2 - wide angle lens",
		"site_id": 10192840284
	}`
	req := httptest.NewRequest("POST", "/api/v1/capturesources", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test setup failed: %s", err.Error())
	}
	if res.StatusCode != 200 {
		t.Fatalf("test setup failed: status code = %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("test setup failed: %s", err.Error())
	}

	csID, err := extractID(resBody)
	if err != nil {
		t.Fatalf("test setup failed: %s", err.Error())
	}

	return app, csID
}

func teardown() {
	// TODO: remove testing capturesource.
}

// extractID extracts the ID from the response.
func extractID(body []byte) (int64, error) {
	var resData struct {
		ID int64
	}

	err := json.Unmarshal(body, &resData)
	if err != nil {
		return 0, err
	}

	return resData.ID, nil
}

// Checks that JSON has the expected format.
func bodyMatchesType[T any](body []byte) bool {
	var format T
	return json.Unmarshal(body, &format) == nil
}

// TestCreateVideoStream checks that the API for creating a video stream returns a 200 (ok) response
// when a valid request is made, and a 400 (bad request) response otherwise.
func TestCreateVideoStream(t *testing.T) {
	app, csID := setup(t)
	defer teardown()

	testcases := []struct {
		reqBody       string // Test case data.
		resStatusCode int    // Expected response: HTTP status code.
		resIsSuccess  bool   // Expected response: Has body with format {"id": <number>} if successful, else { "message": <string> } if not.
	}{
		{
			// Valid requests.
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
				"capturesource": %d,
				"startTime": "2023-06-07T08:00:00.00Z",
				"endTime": "2023-06-07T16:30:00.00Z"
			}`,
			resStatusCode: fiber.StatusOK,
			resIsSuccess:  true,
		},
		{
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=1234567890A",
				"capturesource": %d,
				"startTime": "2023-05-01T20:00:00.00Z",
				"endTime": "2023-05-02T06:15:00.00Z"
			}`,
			resStatusCode: fiber.StatusOK,
			resIsSuccess:  true,
		},
		// Should be a bad request because end time is before start time.
		{
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
				"capturesource": %d,
				"startTime": "2023-05-01T20:00:00.00Z",
				"endTime": "2022-05-01T20:00:00.00Z"
			}`,
			resStatusCode: fiber.StatusOK,
			resIsSuccess:  true,
		},
		// Should be bad requests because missing required fields.
		{
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
				"startTime": "2023-06-07T08:00:00.00Z",
				"endTime": "2023-06-07T16:30:00.00Z"
			}`,
			resStatusCode: fiber.StatusBadRequest,
			resIsSuccess:  false,
		},
		{
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
				"capturesource": %d,
				"endTime": "2023-06-07T16:30:00.00Z"
			}`,
			resStatusCode: fiber.StatusBadRequest,
			resIsSuccess:  false,
		},
		{
			reqBody: `{
				"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
				"capturesource": %d,
				"startTime": "2023-06-07T08:00:00.00Z"
			}`,
			resStatusCode: fiber.StatusBadRequest,
			resIsSuccess:  false,
		},
		// Should be bad requests because they are invalid JSON.
		{
			reqBody:       "something that is not json",
			resStatusCode: fiber.StatusBadRequest,
			resIsSuccess:  false,
		},
		{
			reqBody:       "",
			resStatusCode: fiber.StatusBadRequest,
			resIsSuccess:  false,
		},
	}

	for i, testcase := range testcases {
		// Insert data into test cases.
		reqBody := fmt.Sprintf(testcase.reqBody, csID)

		// Send test request.
		req := httptest.NewRequest("POST", "/api/v1/videostreams", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		if err != nil {
			t.Errorf("[test %d] app.Test failed with error: %s", i, err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("[test %d]  io.ReadAll failed with error: %s", i, err)
		}

		// Check response status code matches expected.
		if res.StatusCode != testcase.resStatusCode {
			t.Errorf("[test %d]  expected HTTP code %d, got %d", i, testcase.resStatusCode, res.StatusCode)
		}

		// Check response body matches expected.
		if testcase.resIsSuccess && !bodyMatchesType[struct{ ID uint64 }](body) {
			t.Errorf(`[test %d]  expected response body to have format { "id": <int> }, got %s`, i, string(body))
		}
		if !testcase.resIsSuccess && !bodyMatchesType[struct{ Message string }](body) {
			t.Errorf(`[test %d]  expected response body to have format { "message": <string> }, got %s`, i, string(body))
		}

		// Clean up after ourselves by removing the video stream.
		if testcase.resIsSuccess {
			id, _ := extractID(body)
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/videostreams/%d", id), nil)
			app.Test(req, -1)
		}
	}
}

func TestDeleteVideoStream(t *testing.T) {
	app, csID := setup(t)
	defer teardown()

	t.Run("it should return 200 for valid ID", func(t *testing.T) {

		// Create a video stream to delete later in test
		createReqBody := `{
			"stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
			"capturesource": %d,
			"startTime": "2023-05-01T20:00:00.00Z",
			"endTime": "2022-05-01T20:00:00.00Z"
		}`
		createReqBody = fmt.Sprintf(createReqBody, csID)

		createReq := httptest.NewRequest("POST", "/api/v1/videostreams", strings.NewReader(createReqBody))
		createReq.Header.Set("Content-Type", "application/json")
		createRes, err := app.Test(createReq, -1)
		if err != nil {
			t.Logf("app.Test failed with error: %s", err)
			t.Fail()
		}
		createResBody, _ := io.ReadAll(createRes.Body)
		id, _ := extractID(createResBody)

		// Send delete request.
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/videostreams/%d", id), nil)
		res, err := app.Test(req, -1)

		if err != nil {
			t.Errorf("app.Test failed with error: %s", err)
		}

		// Check response status code matches expected.
		if res.StatusCode != 200 {
			t.Errorf("expected HTTP code 200, got %d", res.StatusCode)
		}

	})

	t.Run("it should return an error for nonexistent IDs", func(t *testing.T) {
		// Send delete request.
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/videostreams/%d", 1234), nil)
		res, err := app.Test(req, -1)

		if err != nil {
			t.Errorf("app.Test failed with error: %s", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("io.ReadAll failed with error: %s", err)
		}

		// Check response status code matches expected.
		if res.StatusCode != 400 {
			t.Errorf("expected HTTP code 400, got %d", res.StatusCode)
		}

		// Check response body matches expected.
		if !bodyMatchesType[struct{ Message string }](body) {
			t.Errorf(`expected response body to have format { "message": <string> }, got %s`, string(body))
		}
	})
}
