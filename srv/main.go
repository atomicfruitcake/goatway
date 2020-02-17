package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/atomicfruitcake/goatway/constants"
	"github.com/atomicfruitcake/goatway/handlers/createjob"
	"github.com/atomicfruitcake/goatway/handlers/getjob"
	"github.com/atomicfruitcake/goatway/handlers/health"
	"github.com/atomicfruitcake/goatway/handlers/root"
	"github.com/atomicfruitcake/goatway/auth"
	"github.com/atomicfruitcake/goatway/redis"

	"github.com/gorilla/mux"
)

func main() {
	redis.Ping()
	fmt.Println("Starting a new Goatway HTTP Server")
	fmt.Println("Building the Gorilla MUX Router")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", root.Handler).Methods("GET")
	router.HandleFunc("/health", health.Handler).Methods("GET")

	jobrouter := router.PathPrefix("/job").Subrouter()
	jobrouter.HandleFunc("/createJob", createjob.Handler).Methods("POST")
	jobrouter.HandleFunc("/getJob", getjob.Handler).Methods("GET", "POST")

	amw := auth.AuthenticationMiddleware{}
	amw.Populate()
	jobrouter.Use(amw.Middleware)
	http.Handle("/", router)

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second * 15,
		"Graceful Shutdown time is 15 seconds",
	)
	flag.Parse()
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", constants.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	fmt.Printf("Goatway HTTP Webserver is running on port %s\n", constants.Port)
	go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    srv.Shutdown(ctx)

    log.Println("Shutting Down Goatway Server")
    os.Exit(0)
}
