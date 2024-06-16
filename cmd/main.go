package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/vctaragao/book-crud/internal/book"
	"github.com/vctaragao/book-crud/internal/infra/database"
	"github.com/vctaragao/book-crud/internal/infra/http/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load .env file: %v", err)
	}

	dbConn := database.NewConn()

	if err := database.Migrate(dbConn); err != nil {
		log.Fatalf("unable to run migrations: %v\n", err)
	}

	log.Println("migrations applied succefully")

	postgresAdapter := database.NewPostgresAdapter(dbConn)

	bookService := book.NewBookService(&postgresAdapter)

	getBookHandler := handler.NewGetBookHandler(bookService)
	createBookHandler := handler.NewCreateBookHandler(bookService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"message": "ok"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("POST /book", createBookHandler.Handler)
	mux.HandleFunc("GET /book/{id}", getBookHandler.Handler)

	server := http.Server{
		Addr:    ":7777",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server! Listening on :7777...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("unable to start server: %v", err)
		}

		log.Println("Stopped serving new connections")
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	ctx, shoutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shoutdownCancel()

	if err := dbConn.Close(); err != nil {
		log.Printf("unable to close db connection: %v\n", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("unable to close server: %v", err)
	}

	log.Println("server closed!")
}
