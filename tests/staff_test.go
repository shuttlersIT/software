package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStaff(t *testing.T) {
	dept := map[string]string{"name": "IT"}
	resp, _ := PerformRequest("POST", "/api/departments", dept)
	assert.Equal(t, http.StatusOK, resp.Code)

	body := map[string]interface{}{
		"full_name":     "Ada Lovelace",
		"email":         "ada@shuttlers.africa",
		"department_id": 1,
	}
	w, err := PerformRequest("POST", "/api/staff", body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}
