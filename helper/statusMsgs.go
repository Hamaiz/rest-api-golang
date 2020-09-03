package helper

import (
	"encoding/json"
	"net/http"
)

// Unauthorized = 405
func NotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"status": "405", "message": "method not allowed"}`))
}

// Server Error = 500
func ServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"status": "500", "message": "internal server error"}`))
}

// Forbidden = 403
func ForbiddenError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	if err == "" {
		w.Write([]byte(`{"status": "403", "message": "forbidden"}`))
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "403", "message": err})
	}
}

// Unauthorized - 401
func UnauthorizedError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"status": "401", "message": "unauthorized"}`))
}

// Ok - 200
func Ok(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"200", "message": "` + msg + `"}`))
}

// Not Found - 404
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"status": "404", "message": "not found"}`))
}
