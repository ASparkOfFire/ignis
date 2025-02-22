package main

import (
	"github.com/ASparkOfFire/ignis/sdk"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("X-Test-Header", "KabirKalsi")
	if _, err := w.Write([]byte("Hello from WASI")); err != nil {
		log.Printf("Error writing response: %v\n", err)
		return
	}

	return
}

func main() {
	router := chi.NewMux()
	router.Get("/", handleRoot)

	sdk.Handle(router)
}
