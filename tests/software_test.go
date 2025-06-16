package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"software_management/models"
	"software_management/routes"

	"github.com/stretchr/testify/assert"
)

func TestCreateSoftware(t *testing.T) {
	router := routes.RegisterRoutes()

	newSoftware := models.Software{
		Name: "Zoom",
		Type: "Communication",
	}

	body, _ := json.Marshal(newSoftware)
	req, _ := http.NewRequest("POST", "/api/software", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestSoftwareLifecycle(t *testing.T) {
	router := routes.RegisterRoutes()

	// 1. Create software
	software := models.Software{
		Name:        "Google Docs",
		Description: "Collaborative word processor",
		Type:        "SaaS",
	}
	payload, _ := json.Marshal(software)
	req, _ := http.NewRequest("POST", "/api/software", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var created models.Software
	json.Unmarshal(w.Body.Bytes(), &created)

	// 2. Get software
	req, _ = http.NewRequest("GET", "/api/software", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 3. Update software
	updated := map[string]string{"description": "Real-time collaboration"}
	payload, _ = json.Marshal(updated)
	req, _ = http.NewRequest("PUT", "/api/software/"+ToStr(created.ID), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 4. Delete software
	req, _ = http.NewRequest("DELETE", "/api/software/"+ToStr(created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 5. Invalid creation
	invalid := map[string]string{"name": ""}
	payload, _ = json.Marshal(invalid)
	req, _ = http.NewRequest("POST", "/api/software", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func ToStr(id uint) string {
	return strconv.Itoa(int(id))
}
