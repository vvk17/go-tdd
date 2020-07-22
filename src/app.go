package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	fmt.Println("exec: App.Initialize")
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	fmt.Printf("connectionString: %s \n", connectionString)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	fmt.Println("end: App.Initialize")
}
func (a *App) Run(addr string) {}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
