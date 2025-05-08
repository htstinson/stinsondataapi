package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			qp := r.URL.Query()
			dir := qp.Get("dir")
			fld := qp.Get("field")
			fmt.Printf(
				"[%v] [middleware Logger] %s %s %s %s %s %s %s\n",
				time.Now().Format(time.RFC3339),
				r.Method,
				r.URL.Path,
				dir,
				fld,
				r.Header.Get("X-Forwarded-For"), r.URL.User,
				time.Since(start),
			)
		})
	}
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		key := "requestID"
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), key, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	fmt.Println("h CORS middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
