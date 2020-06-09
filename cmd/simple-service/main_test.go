package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Alive)

	handler.ServeHTTP(rr, req)

	expected := `{"alive": true}`

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}
