package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.config.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) MetricsView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		h.config.FileserverHits.Load())
}
