package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	models "github.com/tstxni97/go-simpleservice/cmd/models"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	MIGRATE := viper.GetBool("MIGRATE")

	dsn := "host=localhost user=dp dbname=gorm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Fail to connect db")
	}

	if MIGRATE == true {
		migrate(db)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Fail to call DB")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	create(db)

	r := mux.NewRouter()
	r.HandleFunc("/", Alive)
	r.HandleFunc("/calculator/sum", Sum).Methods("POST")
	r.HandleFunc("/calculator/mul", Mul).Methods("POST")
	r.HandleFunc("/calculator/sub", Sub).Methods("POST")
	r.HandleFunc("/calculator/div", Div).Methods("POST")
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

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.CreditCard{})
	db.AutoMigrate(&models.Dog{})
	db.AutoMigrate(&models.Toy{})
}

func create(db *gorm.DB) {
	card1 := 

	card2 := &models.CreditCard{
		Number: "22222222222222",
	}

	usr := &models.User{
		Name:        "Deeprom",
		Email:       "deeprom@bridgeasiagroup.com",
		Age:         21,
		Birthday:    &time.Time{},
		CreditCards: []models.CreditCard{*card1, *card2},
	}
	db.Create(usr)

}
