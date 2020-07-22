package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	fmt.Println("exec: ensureTableExists")
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
	fmt.Println("end: ensureTableExists")
}

func clearTable() {
	fmt.Println("exec: clearTable")
	fmt.Println("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	fmt.Println("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	fmt.Println("end: clearTable")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)
}`
