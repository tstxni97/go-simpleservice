package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ResBody struct {
	Result float64 `json:"result"`
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

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
	handler := http.HandlerFunc(Sum)

	handler.ServeHTTP(rr, req)

	expected := 300.2

	var resBody ResBody

	if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, resBody.Result)
}

func TestCalculatorMul(t *testing.T) {
	t.Parallel()

	var jsonReq = []byte(`{"a": 100, "b": 200.2}`)
	req, err := http.NewRequest("POST", "/calculator.mul", bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Mul)

	handler.ServeHTTP(rr, req)

	expected := float64(20020)

	var resBody ResBody

	if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, resBody.Result)
}

func TestCalculatorSub(t *testing.T) {
	t.Parallel()

	var jsonReq = []byte(`{"a": 100, "b": 200.2}`)
	req, err := http.NewRequest("POST", "/calculator.mul", bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Sub)

	handler.ServeHTTP(rr, req)

	expected := float64(-100.2)

	var resBody ResBody

	if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
		panic(err)
	}

	assert.Equal(t, true, almostEqual(expected, resBody.Result))
}

func TestCalculatorDiv(t *testing.T) {
	t.Parallel()

	var jsonReq = []byte(`{"a": 100, "b": 200.2}`)
	req, err := http.NewRequest("POST", "/calculator.div", bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Div)

	handler.ServeHTTP(rr, req)

	expected := float64(100 / 200.2)

	var resBody ResBody

	if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
		panic(err)
	}

	assert.Equal(t, true, almostEqual(expected, resBody.Result))

}
