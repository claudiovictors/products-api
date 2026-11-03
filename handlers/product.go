// Package handlers contains the HTTP handlers of my API. Every handler receives
// the shared database pool and the request context, and writes the response
// directly to the ResponseWriter.
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/claudiovictors/crud/models"
)

// Index returns every product stored in the database as a JSON array.
func Index(w http.ResponseWriter, r *http.Request, db *sql.DB, ctx context.Context) {
	query := "SELECT id, name, stock, price FROM products"

	// I run the SELECT and get a cursor over the result rows.
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("failed to query products: %v", err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}
	// I always close the cursor to release the connection back to the pool.
	defer rows.Close()

	// I start with an empty slice so an empty table encodes as [] and not null.
	products := []models.Product{}

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Stock, &product.Price); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}

		products = append(products, product)
	}

	// rows.Next() returns false both on the normal end and on error, so I check
	// rows.Err() to catch failures that happened during iteration.
	if err := rows.Err(); err != nil {
		log.Printf("failed to iterate over rows: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// Store decodes a product from the request body, inserts it and returns the
// created product (with its generated id) as JSON.
func Store(w http.ResponseWriter, r *http.Request, db *sql.DB, ctx context.Context) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO products(name, stock, price) VALUES (?, ?, ?)"

	// I use ExecContext for statements that do not return rows.
	res, err := db.ExecContext(ctx, query, product.Name, product.Stock, product.Price)
	if err != nil {
		log.Printf("failed to insert product: %v", err)
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	// I read the auto-increment id that MySQL generated for the new row.
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("failed to retrieve inserted product id: %v", err)
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	product.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// Update replaces the product identified by the {id} path value.
func Update(w http.ResponseWriter, r *http.Request, db *sql.DB, ctx context.Context) {
	// I read the {id} segment that the router captured from the URL.
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET name = ?, stock = ?, price = ? WHERE id = ?"
	res, err := db.ExecContext(ctx, query, product.Name, product.Stock, product.Price, id)
	if err != nil {
		log.Printf("failed to update product: %v", err)
		http.Error(w, "failed to update product", http.StatusInternalServerError)
		return
	}

	// RowsAffected tells me how many rows matched the WHERE clause.
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("failed to retrieve affected rows: %v", err)
		http.Error(w, "failed to update product", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// Delete removes the product identified by the {id} path value.
func Delete(w http.ResponseWriter, r *http.Request, db *sql.DB, ctx context.Context) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM products WHERE id = ?"
	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to delete product: %v", err)
		http.Error(w, "failed to delete product", http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("failed to retrieve affected rows: %v", err)
		http.Error(w, "failed to delete product", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	// Fprintf writes the body and sends a 200 OK status by default.
	fmt.Fprintf(w, "Product id %s deleted successfully", id)
}
