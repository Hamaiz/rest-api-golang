package serve

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Server - struct for server
type Server struct {
	*http.Server
}

// init - run before the program starts
func init() {
	log.Println("configuring env files...")

	// load godotenv
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found: %v", err)
	}

	// switch between env files depending on enviornment
	switch os.Getenv("APP_ENV") {
	case "development":
		godotenv.Load(".env.development")
	case "production":
		godotenv.Load(".env.production")
	}
}

// NewServer - Starts the http server
func NewServer() (*Server, error) {
	log.Println("configuring server...")

	// getting port from env variables
	port := os.Getenv("PORT")

	var addr string

	// get the handler - ./api.go
	api, err := New()
	if err != nil {
		return nil, err
	}

	// checking if string contains ":" for address
	if strings.Contains(port, ":") {
		addr = port
	} else {
		addr = ":" + port
	}

	// http.Server configuration
	srv := http.Server{
		Handler:      api,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{&srv}, nil
}

// Start - runs server (and hold graceful shutdown)
func (srv *Server) Start() {
	log.Println("starting server...")

	// running server in a go routine
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Printf("Listening on %v...", srv.Addr)

	// logic for graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	log.Println("Shutting down server... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	log.Println("server gracefully stopped")
}
