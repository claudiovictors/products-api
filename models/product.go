// Package models holds the data structures I share across the application.
package models

// Product represents a single row of the products table. I use the struct tags
// to tell encoding/json which key name to use for each field.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}
