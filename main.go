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
	"github.com/kelseyhightower/envconfig"

	"sa-course-app/pkg/infrastructure/database"
	"sa-course-app/pkg/user/infrastructure/mysql"
)

const appName = "saapp"

type Config struct {
	Port       int    `envconfig:"port"`
	DbHost     string `envconfig:"db_host"`
	DbPort     int    `envconfig:"db_port"`
	DbUser     string `envconfig:"db_user"`
	DbPassword string `envconfig:"db_password"`
	DbName     string `envconfig:"db_name"`
}

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

func startHttpServer(port int, router *mux.Router) *http.Server {
	s := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	log.Printf("Try to start server on %v port...\n", port)

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	return s
}

func main() {
	var c Config
	err := envconfig.Process(appName, &c)
	if err != nil {
		log.Fatal(err)
	}

	dsn := database.NewDSN(c.DbHost, strconv.Itoa(c.DbPort), c.DbUser, c.DbPassword, c.DbName)
	conn, err := database.NewMySQLConnector(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userRepository := mysql.NewUserRepository(conn.Client())
	appService := app.NewService(userRepository)

	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := buildRouter(appService)
	srv := startHttpServer(c.Port, router)

	sig := <-osSignalChan
	log.Printf("OS signal received: %+v", sig)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Print("Server shutdown successfully!")
}
