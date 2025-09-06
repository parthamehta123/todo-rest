package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	h "todo-rest/internal/http"
	"todo-rest/internal/todo"
)

func main() {
	store := todo.NewMemoryStore()
	srv := h.New(store)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      srv.Mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("REST server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Println("server stopped")
}
