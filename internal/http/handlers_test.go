package http

import (
	"net/http/httptest"
	"testing"

	"todo-rest/internal/todo"
)

func TestListEmpty(t *testing.T) {
	s := New(todo.NewMemoryStore())
	req := httptest.NewRequest("GET", "/v1/todos/", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}
}
