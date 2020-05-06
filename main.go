package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sa-course-app/pkg/user/app"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"

	"sa-course-app/pkg/infrastructure/database"
	"sa-course-app/pkg/user/infrastructure/mysql"
)

const portEnvVar = "SERVICE_PORT"

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

func buildRouter(appService *app.Service) *mux.Router {
	h := handlers{s: appService}

	r := mux.NewRouter()
	r.HandleFunc("/ready", h.readyHandler).Methods(http.MethodGet)
	r.HandleFunc("/health", h.healthHandler).Methods(http.MethodGet)

	r.HandleFunc("/user", h.createUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", h.getUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", h.updateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/user/{id}", h.deleteUserHandler).Methods(http.MethodDelete)

	r.Use(addResponseContentTypeJSONMiddleware, logMiddleware)

	return r
}

func startHttpServer(router *mux.Router) *http.Server {
	port := os.Getenv(portEnvVar)
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Can not use \"%v\" as HTTP port.", port)
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: router,
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

	userRepository := mysql.NewUserRepository(c.Client())
	appService := app.NewService(userRepository)

	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := buildRouter(appService)
	srv := startHttpServer(router)

	sig := <-osSignalChan
	log.Printf("OS signal received: %+v", sig)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Print("Server shutdown successfully!")
}
