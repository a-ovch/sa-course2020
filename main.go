package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const serverAddr = ":8000"

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
