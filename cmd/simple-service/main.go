package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type ReqBody struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// WriteJSON writes JSON to the response with the given http status
func WriteJSON(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Alive)
	r.HandleFunc("/calculator.sum", Sum).Methods("POST")
	r.HandleFunc("/calculator.mul", Mul).Methods("POST")
	r.HandleFunc("/calculator.sub", Sub).Methods("POST")
	r.HandleFunc("/calculator.div", Div).Methods("POST")
	r.Use(mux.CORSMethodMiddleware(r))

	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Alive(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func Sum(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	// decode request
	var reqBody ReqBody
	err := decoder.Decode(&reqBody)
	if err != nil {
		err = errors.Wrap(err, "Malformed request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := reqBody.A + reqBody.B
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}

func Mul(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	// decode request
	var reqBody ReqBody
	err := decoder.Decode(&reqBody)
	if err != nil {
		err = errors.Wrap(err, "Malformed request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := reqBody.A * reqBody.B
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}

func Sub(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	// decode request
	var reqBody ReqBody
	err := decoder.Decode(&reqBody)
	if err != nil {
		err = errors.Wrap(err, "Malformed request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := reqBody.A - reqBody.B

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}

func Div(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	// decode request
	var reqBody ReqBody
	err := decoder.Decode(&reqBody)
	if err != nil {
		err = errors.Wrap(err, "Malformed request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := reqBody.A / reqBody.B
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}
