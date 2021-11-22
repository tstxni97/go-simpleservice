package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
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

var (
	// global db struct
	db   *gorm.DB
	wait time.Duration
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to read config file: %w ", err))
	}

	MIGRATE := viper.GetBool("MIGRATE")
	ENV := viper.GetString("ENV")
	db_connect := viper.GetString("db.local")
	fmt.Println("migrate", MIGRATE)
	fmt.Println("ENV is ", ENV)

	db, err = gorm.Open(postgres.Open(db_connect), &gorm.Config{})
	if err != nil {
		panic("Fail to connect db")
	}

	if MIGRATE {
		fmt.Println("migrating ORM")
		migrate(db)
		create(db)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Fail to call DB")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	r := mux.NewRouter()
	r.HandleFunc("/", Alive)
	r.HandleFunc("/calculator/sum", Sum).Methods("POST")
	r.HandleFunc("/calculator/mul", Mul).Methods("POST")
	r.HandleFunc("/calculator/sub", Sub).Methods("POST")
	r.HandleFunc("/calculator/div", Div).Methods("POST")
	r.HandleFunc("/users", GetUserByID).Methods("GET", "OPTIONS").Queries("id", "{id}")
	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Addr: "0.0.0.0:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	http.Handle("/", r)
	fmt.Println("Starting up on 3000")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
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
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(err.Error()))
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

	usr := &models.User{
		Name:     "Deeprom",
		Email:    "deeprom@bridgeasiagroup.com",
		Age:      21,
		Birthday: &time.Time{},
		CreditCards: []models.CreditCard{
			{
				Number: "11111111111111",
			},
			{
				Number: "22222222222222",
			}},
	}
	db.Create(usr)

}

func GetUserByID(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}
	ID := req.URL.Query().Get("id")
	// If the primary key is a string (for example, like a uuid), the query will be written as follows:
	// db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
	usr := &models.User{}
	if err := db.First(usr, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(err, "")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code":   "200",
		"result": usr,
	})
}
