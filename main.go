package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.HandleFunc("/healthz", apiCfg.healthz)
	mux.HandleFunc("/metrics", apiCfg.metricsView)
	mux.HandleFunc("/reset", apiCfg.metricsReset)

	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("app")))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}
