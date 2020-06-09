package main

import (
	"bytes"
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

func TestCalculatorSum(t *testing.T) {
	t.Parallel()

	var jsonReq = []byte(`{"a": 100, "b": 200.2}`)
	req, err := http.NewRequest("POST", "/calculator.sum", bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Alive)

	handler.ServeHTTP(rr, req)

	expected := `{"result": 300.2}`

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}
