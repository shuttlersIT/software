package tests

import (
	//"encoding/json"
	"bytes"
	"net/http"
	"net/http/httptest"
	"software_management/routes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDepartment(t *testing.T) {
	router := routes.RegisterRoutes()
	if router == nil {
		t.Fatal("router is nil – check routes.RegisterRoutes()")
	}

	body := []byte(`{"name": "IT"}`)

	req, _ := http.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetDepartments(t *testing.T) {
	w, err := PerformRequest("GET", "/api/departments", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}
