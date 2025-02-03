package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(main.DbUser, main.DbPass, main.DbHost, main.DbPort, "test")
	if err != nil {
		log.Fatal(err)
	}
	CreateTable()
	code := m.Run()
	os.Exit(code)
}

func CreateTable() {
	createTableQuery := `create table if not exists products (
		id int not null AUTO_INCREMENT,
		name varchar(255) not null,
		quantity int,
		price float(10,7),
		primary key (id)
	);`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("delete from products;")
	a.DB.Exec("alter table products AUTO_INCREMENT=1;")
	log.Println("Cleared table")
	if err != nil {
		log.Fatal(err)
	}
}

func TestAddProduct(t *testing.T) {
	clearTable()
	addProduct("John Doe", 10, 10.0)
	addProduct("Mark Brien", 20, 11.0)
	log.Println("Added product")
}

func TestGetProducts(t *testing.T) {
	log.Println("got the product")
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendNewResponse(req)
	CheckStatusCode(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{"name":"John Doe", "quantity":10, "price":10}`)
	req, _ := http.NewRequest("POST", "/product/", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendNewResponse(req)
	CheckStatusCode(t, http.StatusCreated, response.Code)
	log.Printf("Response Status Code: %d", response.Code)
	log.Printf("Response Body: %s", response.Body.String())
	var m map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		t.Fatalf("JSON Unmarshal error: %v. Raw body: %s", err, response.Body.String())
	}
	log.Printf("Parsed JSON: %+v", m)
	if m["Name"] != "John Doe" {
		t.Errorf("Expected name to be 'John Doe'. Got '%v'", m["name"])
	}

	if quantity, ok := m["Quantity"].(float64); !ok || int(quantity) != 10 {
		t.Errorf("Expected quantity to be 10. Got '%v'", m["quantity"])
	}

	if price, ok := m["Price"].(float64); !ok || int(price) != 10 {
		t.Errorf("Expected price to be 10. Got '%v'", m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("John Doe", 10, 10.0)
	addProduct("Mark Brien", 20, 11.0)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendNewResponse(req)
	CheckStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendNewResponse(req)
	CheckStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/2", nil)
	response = sendNewResponse(req)
	CheckStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/2", nil)
	response = sendNewResponse(req)
	CheckStatusCode(t, http.StatusNotFound, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendNewResponse(req)
	CheckStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("John Doe", 10, 10.0)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendNewResponse(req)
	CheckStatusCode(t, http.StatusOK, response.Code)
	var oldValue map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &oldValue)
	if err != nil {
		t.Fatalf("JSON Unmarshal error: %v. Raw body: %s", err, response.Body.String())
	}
	product := []byte(`{"name":"John Doe", "quantity":10, "price":100}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response = sendNewResponse(req)
	log.Printf("Response Status Code: %d", response.Code)
	log.Printf("Response Body: %s", response.Body.String())
	var newValue map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &newValue)
	if err != nil {
		t.Fatalf("JSON Unmarshal error: %v. Raw body: %s", err, response.Body.String())
	}
	if oldValue["Name"] != newValue["Name"] {
		log.Printf("Name Got updated. Got '%v' '%v'", oldValue["Name"], newValue["Name"])
	}
	if oldValue["Quantity"] != newValue["Quantity"] {
		log.Printf("Quantity Got updated. Got '%v'  '%v'", oldValue["Quantity"], newValue["Quantity"])
	}
	if oldValue["Price"] != newValue["Price"] {
		log.Printf("Price Got updated. Got '%v' '%v'", oldValue["Price"], newValue["Price"])
	}
}

func addProduct(name string, quantity int, price float64) {
	query := `insert into products(name, quantity, price) values(?, ?, ?)`
	_, err := a.DB.Exec(query, name, quantity, price)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code to be %d, got %d", expectedStatusCode, actualStatusCode)
	}
}

func sendNewResponse(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}
