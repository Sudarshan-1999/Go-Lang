package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(DbUser string, DbPass string, DbHost string, DbPort string, DbName string) error {
	var err error
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", DbUser, DbPass, DbHost, DbPort, DbName)
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(addr string) error {
	log.Fatal(http.ListenAndServe(addr, app.Router))
	return nil
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func sendResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)

}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product/", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusOK, products)
}

func (app *App) getProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	p := product{
		ID: key,
	}
	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(writer, http.StatusNotFound, "Product not found")
		default:
			sendError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(writer, http.StatusOK, p)
}

func (app *App) createProduct(writer http.ResponseWriter, request *http.Request) {
	var p product
	err := json.NewDecoder(request.Body).Decode(&p)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = p.createProduct(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusCreated, p)
}

// update product
func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	var p product
	err = json.NewDecoder(request.Body).Decode(&p)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	p.ID = key
	err = p.updateProduct(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusOK, p)
}

// Delete products
func (app *App) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	p := product{ID: key}
	err = p.deleteProduct(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusOK, map[string]string{"status": "Product deleted successfully"})
}
