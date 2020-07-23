package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func ensureTableExists() {
	//	fmt.Println("exec: ensureTableExists")
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
	//	fmt.Println("end: ensureTableExists")
}

func clearTable() {
	//	fmt.Println("exec: clearTable")
	//	fmt.Println("DELETE FROM products")
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	//	fmt.Println("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	//	fmt.Println("end: clearTable")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
 id SERIAL,
 name TEXT NOT NULL,
 price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
 CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	//	fmt.Println("exec: executeRequest")
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	//	fmt.Println("end: executeRequest")
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	//	fmt.Println("exec: checkResponseCode")
	//	fmt.Printf("code expected: %d, actual %d\n", expected, actual)
	if expected != actual {
		t.Errorf("Expected responce code %d. Got %d\n", expected, actual)
	}
	//	fmt.Println("end: checkResponseCode")
}
