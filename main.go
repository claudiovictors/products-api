// Package main is the entry point of my API. Here I build the router, open the
// database connection pool and start the HTTP server.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/claudiovictors/crud/database"
	"github.com/claudiovictors/crud/handlers"
)

func main() {
	// I create the router that maps request patterns to handlers.
	app := http.NewServeMux()

	// I open the MySQL connection pool once and share it with every handler.
	db, err := database.Connection()
	if err != nil {
		log.Fatal(err)
	}

	// Here I register the routes. The "METHOD /path" pattern is supported
	// natively by net/http since Go 1.22, so I did not need a third-party
	// router. I inject the shared database pool and the request context into
	// each handler.
	app.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		handlers.Index(w, r, db, r.Context())
	})

	app.HandleFunc("POST /products", func(w http.ResponseWriter, r *http.Request) {
		handlers.Store(w, r, db, r.Context())
	})

	app.HandleFunc("PUT /products/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Update(w, r, db, r.Context())
	})

	app.HandleFunc("DELETE /products/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Delete(w, r, db, r.Context())
	})

	// ListenAndServe blocks until the server is shut down or fails to start.
	fmt.Println("Server running on port :8080")
	if err := http.ListenAndServe(":8080", app); err != nil {
		panic(err)
	}
}
