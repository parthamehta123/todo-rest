package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"todo-rest/internal/todo"
)

type Server struct {
	Store todo.Store
	Mux   *chi.Mux
}

func New(store todo.Store) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	s := &Server{Store: store, Mux: r}
	r.Route("/v1/todos", func(r chi.Router) {
		r.Get("/", s.list)
		r.Post("/", s.create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.get)
			r.Patch("/", s.update)
			r.Delete("/", s.delete)
		})
	})
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	return s
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Store.List())
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || len(in.Title) == 0 {
		http.Error(w, "invalid body: title required", http.StatusBadRequest)
		return
	}
	it := s.Store.Create(in.Title)
	writeJSON(w, http.StatusCreated, it)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	it, err := s.Store.Get(id)
	if errors.Is(err, todo.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id, err := atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var in struct {
		Title *string `json:"title"`
		Done  *bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	it, err := s.Store.Update(id, in.Title, in.Done)
	if errors.Is(err, todo.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	id, err := atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := s.Store.Delete(id); errors.Is(err, todo.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func atoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
