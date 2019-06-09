package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body struct {
		Status string `json:"status"`
	}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "OK", body.Status)
}
