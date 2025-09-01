package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/HemahWeb/chirpy/internal/database"
	"github.com/HemahWeb/chirpy/internal/handlers"
	"github.com/HemahWeb/chirpy/internal/types"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := types.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
	}

	handler := handlers.New(&apiCfg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handler.Healthz)

	// Chirps
	mux.HandleFunc("POST /api/chirps", handler.PostChirps)
	mux.HandleFunc("GET /api/chirps", handler.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", handler.GetChirpsByID)

	// Users
	mux.HandleFunc("POST /api/users", handler.UsersCreate)
	mux.HandleFunc("POST /api/login", handler.Login)

	// Admin
	mux.HandleFunc("POST /admin/reset", handler.UsersReset) // resets users and metrics
	mux.HandleFunc("GET /admin/metrics", handler.MetricsView)

	mux.Handle("/app/", http.StripPrefix("/app/", handler.MiddlewareMetricsInc(http.FileServer(http.Dir("app")))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}
