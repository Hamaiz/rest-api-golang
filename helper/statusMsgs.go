package helper

import (
	"net/http"
	"strconv"
)

// Not Found - 404
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"status": "404", "message": "not found"}`))
}

// status struct holds status msg
type statusStruct struct {
	status int
	msg    string
}

// ASM - send any desired message
func ASM(w http.ResponseWriter, s int, msg string) {
	w.Header().Set("Content-Type", "application/json")

	// status - map holds all status msgs
	status := map[int]statusStruct{
		200: {http.StatusOK, "ok"},
		201: {http.StatusCreated, "created"},
		204: {http.StatusNoContent, "no content"},
		401: {http.StatusUnauthorized, "unauthorized"},
		403: {http.StatusForbidden, "forbidden"},
		404: {http.StatusNotFound, "not found"},
		405: {http.StatusMethodNotAllowed, "method not allowed"},
		409: {http.StatusConflict, "conflict"},
		422: {http.StatusUnprocessableEntity, "unprocessable entity"},
		500: {http.StatusInternalServerError, "internal server error"},
	}

	sm := status[s].msg

	// status header
	w.WriteHeader(status[s].status)

	if msg != "" {
		w.Write([]byte(`{"status": "` + strconv.Itoa(s) + `", "message": "` + msg + `"}`))
	} else {
		w.Write([]byte(`{"status": "` + strconv.Itoa(s) + `", "message": "` + sm + `"}`))
	}
}
