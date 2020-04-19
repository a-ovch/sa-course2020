package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
)

const portEnvVar = "SERVICE_PORT"

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func httpLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v\n", r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func startHttpServer() *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	port := os.Getenv(portEnvVar)
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Can not use \"%v\" as HTTP port.", port)
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: httpLogMiddleware(r),
	}

	log.Printf("Server successfully started on %v port\n", port)

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	return s
}

func main() {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := startHttpServer()

	sig := <-osSignalChan
	log.Printf("OS signal received: %+v", sig)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Print("Server shutdown successfully!")
}
