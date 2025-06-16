package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"software_management/routes"

	"github.com/stretchr/testify/assert"
)

func TestGetAllLogs(t *testing.T) {
	router := routes.RegisterRoutes()

	req, _ := http.NewRequest("GET", "/api/software-assignment-logs/details", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestSoftwareAssignmentLogs(t *testing.T) {
	router := routes.RegisterRoutes()

	log := map[string]any{
		"software_id": 1,
		"staff_id":    1,
		"action":      "assigned",
		"note":        "Manual assignment",
	}
	payload, _ := json.Marshal(log)
	req, _ := http.NewRequest("POST", "/api/software-assignment-logs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Get all logs
	req, _ = http.NewRequest("GET", "/api/software-assignment-logs", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Invalid log entry
	invalid := map[string]any{"software_id": 0, "staff_id": 0, "action": ""}
	payload, _ = json.Marshal(invalid)
	req, _ = http.NewRequest("POST", "/api/software-assignment-logs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
