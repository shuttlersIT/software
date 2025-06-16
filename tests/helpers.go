package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func PerformRequest(method, path string, body interface{}) (*httptest.ResponseRecorder, error) {
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w, nil
}
