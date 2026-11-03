<p align="center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg" alt="Go logo" height="90">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white" alt="Go version">
  <img src="https://img.shields.io/badge/stdlib-net%2Fhttp-00ADD8?logo=go&logoColor=white" alt="Go standard library">
  <img src="https://img.shields.io/badge/database-MySQL%208.4-4479A1?logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/docker-compose-2496ED?logo=docker&logoColor=white" alt="Docker Compose">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

<h1 align="center">Products API</h1>

<p align="center">A small REST API I wrote in Go to manage products.</p>

---

## About

I built this API to practice writing a REST service in Go using the standard
library as much as possible, without a web framework. I wanted to really
understand how routing, JSON encoding and decoding, contexts and database
access work under the hood before reaching for a framework.

I used the native router that ships with `net/http` (available since Go 1.22),
so I define routes such as `GET /products` and `PUT /products/{id}` myself,
without a third-party router like Chi or Gorilla.

I kept the whole HTTP layer on the Go standard library. The only thing I had to
bring in from outside is the official MySQL driver
(`github.com/go-sql-driver/mysql`), because the standard library does not ship a
driver for MySQL.

## What this API does

- I can create, read, update and delete products (full CRUD).
- I route the requests with `net/http` (no web framework).
- I send and receive JSON.
- I persist everything in MySQL through `database/sql` and one shared
  connection pool.
- I configure the connection with environment variables.
- I made it run with a single Docker Compose command (API + MySQL).

## Tech stack

Here is what I chose to build this:

| Layer         | Choice                                                       |
| ------------- | ------------------------------------------------------------ |
| Language      | Go 1.27                                                      |
| HTTP          | `net/http` (standard library, native routing)                |
| JSON          | `encoding/json` (standard library)                           |
| Database      | MySQL 8.4                                                    |
| DB access     | `database/sql` + `github.com/go-sql-driver/mysql`            |
| Container     | Docker + Docker Compose                                      |

## How I organized the project

```text
rocketseats/
├── main.go               # I set up the router and start the server here
├── database/
│   └── database.go       # My MySQL connection pool
├── handlers/
│   └── product.go        # My HTTP handlers (Index, Store, Update, Delete)
├── models/
│   └── product.go        # The Product struct I share across the app
├── schema.sql            # The table schema my Docker MySQL container uses
├── Dockerfile            # The multi-stage build I use for the API image
├── docker-compose.yml    # The API + MySQL stack I put together
├── go.mod
└── go.sum
```

## What I need before running it

- Go 1.27 or newer (only if I want to run it locally).
- A MySQL instance, or Docker to use the Compose setup.

## How I configure it

I read the database settings from environment variables. When I do not set a
variable, the default value (between parentheses) is used.

| Variable      | Description                | Default           |
| ------------- | -------------------------- | ----------------- |
| `DB_USER`     | MySQL user                 | `root`            |
| `DB_PASSWORD` | MySQL password             | `Root@123`        |
| `DB_ADDRESS`  | MySQL host and port        | `127.0.0.1:3306`  |
| `DB_NAME`     | Database name              | `products`        |

## Running with Docker Compose

I wrote a Compose file that starts both the API and a MySQL instance, and
creates the `products` table automatically the first time it boots.

```bash
docker compose up --build
```

I can then reach the API at `http://localhost:8080`.

To stop everything:

```bash
docker compose down
```

To stop everything and also delete the database volume:

```bash
docker compose down -v
```

## Running locally

1. I create the database and the table:

```sql
CREATE DATABASE IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products (
    id    INT AUTO_INCREMENT PRIMARY KEY,
    name  VARCHAR(255)   NOT NULL,
    stock INT            NOT NULL DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0
);
```

2. I download the dependencies and start the server:

```bash
go mod download
go run .
```

My server listens on `http://localhost:8080`.

## The endpoints I created

| Method   | Path             | What it does                    | Success |
| -------- | ---------------- | ------------------------------- | ------- |
| `GET`    | `/products`      | Lists all products              | `200`   |
| `POST`   | `/products`      | Creates a new product           | `201`   |
| `PUT`    | `/products/{id}` | Updates a product by id         | `200`   |
| `DELETE` | `/products/{id}` | Deletes a product by id         | `200`   |

### The product payload I use

```json
{
  "id": 1,
  "name": "Gopher Plush",
  "stock": 10,
  "price": 29.9
}
```

## Trying it out

I create a product:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Gopher Plush","stock":10,"price":29.9}'
```

I list all products:

```bash
curl http://localhost:8080/products
```

I update a product:

```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Gopher Plush XL","stock":5,"price":39.9}'
```

I delete a product:

```bash
curl -X DELETE http://localhost:8080/products/1
```

## What I learned

- How `net/http` handles method-based routing and path parameters (`r.PathValue`).
- Why I have to check `rows.Err()` after iterating over a query result.
- How `database/sql` manages a connection pool I share across handlers.
- How to pass the request `context` down to the database layer so my queries are
  cancelled when the client disconnects.

## License

I released this project under the MIT License.
