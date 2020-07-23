package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
)

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

func TestAbs(t *testing.T) {
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	//	fmt.Println(req)
	response := executeRequest(req)
	fmt.Println(response.Code)

	checkResponseCode(t, http.StatusOK, response.Code)
	body := response.Body.String()
	fmt.Printf("body: %s\n", body)

	if body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistingProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	fmt.Println(response.Code)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product",
		bytes.NewBuffer(jsonStr))
	//	fmt.Println(req)
	req.Header.Set("Content-Type", "application/json")
	//	fmt.Println(req)
	response := executeRequest(req)
	fmt.Println(response.Code)
	checkResponseCode(t, http.StatusCreated, response.Code)

	body := response.Body.String()
	fmt.Printf("body: %s\n", body)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)

	if m["name"] != "test product" {
		t.Errorf("Expected prod name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected prod price to be '11.22'. Got '%v'", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected prod ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(13)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	fmt.Printf(" response code: %d \n responce body: %s\n", response.Code, response.Body.String())
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/products", nil)

	q := req.URL.Query()
	q.Add("start", "3")
	q.Add("count", "7")
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL)
	fmt.Println(req.URL.RawQuery)

	response = executeRequest(req)
	fmt.Printf(" response code: %d \n responce body: %s\n", response.Code, response.Body.String())

	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)
	fmt.Println(response.Code)
	fmt.Println(originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.23}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println(originalProduct)
	fmt.Println(m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
