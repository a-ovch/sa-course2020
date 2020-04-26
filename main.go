package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"

	"sa-course-app/pkg/infrastructure/database"
	"sa-course-app/pkg/user/domain"
	"sa-course-app/pkg/user/infrastructure/mysql"
)

const portEnvVar = "SERVICE_PORT"

var us *domain.UserService

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func addResponseContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, "{\"host\": \"%v\"}", r.Host)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func createUserHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Create: OK")
}

func updateUserHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Update: OK")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Incorrect 'id' value: %v", vars["id"])
		return
	}

	if u, err := us.FindUser(domain.UserID(id)); u != nil && err == nil {
		bytes, _ := json.Marshal(u)
		_, _ = w.Write(bytes)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteUserHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Delete: OK")
}

func buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ready", readyHandler).Methods(http.MethodGet)
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	r.HandleFunc("/user", createUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", getUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", updateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/user/{id}", deleteUserHandler).Methods(http.MethodDelete)

	r.Use(addResponseContentTypeJSONMiddleware, logMiddleware)

	return r
}

func startHttpServer() *http.Server {
	port := os.Getenv(portEnvVar)
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Can not use \"%v\" as HTTP port.", port)
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: buildRouter(),
	}

	log.Printf("Server successfully started on %v port\n", port)

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	return s
}

func main() {
	dsn := database.NewDSN("localhost", "3306", "root", "12345Q", "mydb")
	c, err := database.NewMySQLConnector(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ur := mysql.NewUserRepository(c.Client())
	us = domain.NewUserService(ur)

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
