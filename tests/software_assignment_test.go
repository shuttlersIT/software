package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"software_management/models"
	"software_management/routes"

	"github.com/stretchr/testify/assert"
)

func TestAssignSoftwareToStaff(t *testing.T) {
	router := routes.RegisterRoutes()

	assignment := models.AssignedSoftware{
		SoftwareID: 1,
		StaffID:    1,
	}

	body, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("POST", "/api/software-assignments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestSoftwareAssignment(t *testing.T) {
	router := routes.RegisterRoutes()

	// Assume software_id = 1, staff_id = 1 already exist
	assignment := map[string]interface{}{
		"software_id": 1,
		"staff_id":    1,
		"note":        "Initial access",
	}
	payload, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("POST", "/api/software-assignments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Edge Case: Duplicate assignment
	req, _ = http.NewRequest("POST", "/api/software-assignments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "already assigned")

	// List assignments
	req, _ = http.NewRequest("GET", "/api/software-assignments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
