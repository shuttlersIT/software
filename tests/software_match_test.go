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

func TestCreateDepartmentMatch(t *testing.T) {
	router := routes.RegisterRoutes()

	match := models.SoftwareDepartmentMatch{
		SoftwareID:   1,
		DepartmentID: 1,
	}

	body, _ := json.Marshal(match)
	req, _ := http.NewRequest("POST", "/api/software_department_matches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestSoftwareMatches(t *testing.T) {
	router := routes.RegisterRoutes()

	// Create org-wide match
	orgMatch := map[string]any{"software_id": 1}
	payload, _ := json.Marshal(orgMatch)
	req, _ := http.NewRequest("POST", "/api/software_organization_matches", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Create department match
	deptMatch := map[string]any{"software_id": 1, "department_id": 1}
	payload, _ = json.Marshal(deptMatch)
	req, _ = http.NewRequest("POST", "/api/software_department_matches", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Update department match (assume ID 1)
	update := map[string]any{"software_id": 1, "department_id": 2}
	payload, _ = json.Marshal(update)
	req, _ = http.NewRequest("PUT", "/api/software_department_matches/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Delete department match
	req, _ = http.NewRequest("DELETE", "/api/software_department_matches/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
