package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type product struct {
	ID       int     "json:id"
	Name     string  "json:name"
	Quantity int     "json:quantity"
	Price    float64 "json:price"
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	products := []product{}
	for rows.Next() {
		var p product
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := "INSERT INTO products(name, quantity, price) VALUES(?, ?, ?)"

	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	log.Printf("Executing query: %s with values: %v, %v, %v", query, p.Name, p.Quantity, p.Price)
	log.Println(result, err)
	if err != nil {
		//log.Fatal(err)
		log.Println(result)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM products WHERE id=%d", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name=?, quantity=?, price=? where id=?")
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if n == 0 {
		return errors.New("Product does not exists")
	}
	return nil
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=?")
	result, err := db.Exec(query, p.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return errors.New("Product not exist")
	}
	return nil
}
