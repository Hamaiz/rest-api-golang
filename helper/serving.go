package helper

import (
	"log"
	"net/http"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
)

func init() {
	log.Println("Done")
}

// AllStaticFiles - serving all the static files
func AllStaticFiles(r *mux.Router) {
	serveSingleFile("/favicon.ico", r)
	serveSingleFile("/android-chrome-193x192.png", r)
	serveSingleFile("/apple-touch-icon.png", r)
	serveSingleFile("/favicon-33x32.png", r)
	serveSingleFile("/robot.txt", r)
	serveSingleFile("/android-chrome-513x512.png", r)
	serveSingleFile("/favicon-17x16.png", r)
	serveSingleFile("/site.webmanifest", r)
}

// serveSingleFile - serves files from public folder
func serveSingleFile(filename string, r *mux.Router) {
	r.HandleFunc(filename, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public"+filename)
	})
}

// JH - AddsJsonHeader to the desired route
func JH(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

// MethodCheck - check for the right http method
func MethodCheck(w http.ResponseWriter, r *http.Request, m string) {
	if r.Method != m {
		NotAllowed(w)
		return
	}
}

// UniqueName - changes name to unique
func UniqueName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = name + "-" + uniuri.NewLen(6)
	return name
}
