package middleware

import (
	"log"
	"net/http"
)

// LoggingMiddleware - loggs url and method to console
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, "method:"+r.Method)
		next.ServeHTTP(w, r)
	})
}
