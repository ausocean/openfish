package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/handlers"
	"github.com/ausocean/openfish/api/utils"
	"github.com/gofiber/fiber/v2"
)

func TestCreateCaptureSource(t *testing.T) {
	// Test Setup.
	ds_client.Init(false)
	app := fiber.New()
	RegisterAPIRoutes(app)

	// Send test request.
	body := strings.NewReader(`{ "name": "Camera 1", "location": "-37.12345678,140.12345678" }`)
	req := httptest.NewRequest("POST", "/api/v1/capturesources", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// Check response.
	if resp.StatusCode != fiber.StatusOK {
		t.Log("expected HTTP code 200, got", resp.StatusCode)
		t.Fail()
	}
}

func TestGetCaptureSources(t *testing.T) {
	// Test Setup.
	ds_client.Init(false)
	app := fiber.New()
	RegisterAPIRoutes(app)

	// Populate datastore.
	body := strings.NewReader(`{ "name": "Camera 1", "location": "-37.12345678,140.12345678" }`)
	req := httptest.NewRequest("POST", "/api/v1/capturesources", body)
	req.Header.Set("Content-Type", "application/json")

	app.Test(req)

	// Send request
	req = httptest.NewRequest("GET", "/api/v1/capturesources?limit=1&offset=0", nil)
	resp, _ := app.Test(req)

	// Check response.
	valid := true

	if resp.StatusCode != fiber.StatusOK {
		t.Log("expected HTTP code 200, got", resp.StatusCode)
		valid = false
	}

	// Unmarshal response.
	respBody, _ := ioutil.ReadAll(resp.Body)
	var result utils.Result[handlers.CaptureSourceResult]
	json.Unmarshal(respBody, &result)
	t.Log("json response: ", string(respBody))

	if len(result.Results) != 1 {
		t.Log("expected 1 result, got", len(result.Results), "results")
		valid = false
	}

	if result.Total != 1 {
		t.Log("expected total to be equal to 1, got total =", result.Total)
		valid = false
	}

	if result.Limit != 1 {
		t.Log("expected limit to be equal to 1, got limit =", result.Limit)
		valid = false
	}

	if result.Offset != 0 {
		t.Log("expected offset to be equal to 0, got offset =", result.Offset)
		valid = false
	}

	if !valid {
		t.Fail()
	}

}
