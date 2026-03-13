package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/guilycst/GoATTH-penguinui/internal/server"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8090", "Server port")
	flag.Parse()

	// Get the project root (one level up from cmd/server)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	srv := server.New(projectRoot)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("Button Component Demo: http://localhost%s/components/button", addr)

	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
