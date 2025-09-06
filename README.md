# Todo REST API (Go + chi)

A simple RESTful API for managing todos, built in Go using the [chi](https://github.com/go-chi/chi) router.  
This project demonstrates clean project structure, an in-memory store, and basic CRUD operations with JSON.

---

## ✨ Features

- Health check endpoint (`/healthz`)
- CRUD endpoints under `/v1/todos/`:
  - `POST /v1/todos/` → Create a todo
  - `GET /v1/todos/` → List all todos
  - `GET /v1/todos/{id}/` → Get a specific todo
  - `PATCH /v1/todos/{id}/` → Update a todo (title and/or mark as done)
  - `DELETE /v1/todos/{id}/` → Delete a todo
- Graceful server shutdown
- In-memory store (thread-safe with RWMutex)

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Make (optional, for shortcuts)
- curl (for testing)

### Run the server

```bash
go run ./cmd/server
```

or using Make:

```bash
make run
```

The server listens on `http://localhost:8080`.

---

## 📚 Example Usage

### Health check
```bash
curl -s localhost:8080/healthz
# {"status":"ok"}
```

### Create a todo
```bash
curl -s -XPOST localhost:8080/v1/todos/   -d '{"title":"learn go"}'   -H 'Content-Type: application/json'
```

### List todos
```bash
curl -s localhost:8080/v1/todos/
```

### Get a todo by ID
```bash
curl -s localhost:8080/v1/todos/1/
```

### Update (mark as done)
```bash
curl -s -XPATCH localhost:8080/v1/todos/1/   -d '{"done":true}'   -H 'Content-Type: application/json'
```

### Delete
```bash
curl -s -XDELETE localhost:8080/v1/todos/1/
```

---

## 🧪 Development

Run tests:
```bash
make test
```

Build binary:
```bash
make build
```

Run with Docker:
```bash
docker build -t todo-rest .
docker run -p 8080:8080 todo-rest
```

---

## 📖 Notes

- Data is **not persisted** (in-memory store only). Restarting the server clears all todos.
- Good starting point for learning Go REST APIs with chi.
- Can be extended with a real database, authentication, or OpenAPI/Swagger.
