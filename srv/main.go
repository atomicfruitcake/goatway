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
	"github.com/atomicfruitcake/goatway/middleware/auth"
	"github.com/atomicfruitcake/goatway/middleware/logging"
	"github.com/atomicfruitcake/goatway/redis"

	"github.com/gorilla/mux"
)

func main() {
	err := redis.Ping()
	if err != nil {
		log.Fatal("Could not connect to Redis, cannot boot Goatway")
		
	}
	log.Println("Starting a new Goatway HTTP Server")
	log.Println("Building the Gorilla MUX Router")
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", root.Handler).Methods("GET")
	r.HandleFunc("/health", health.Handler).Methods("GET")
	r.Use(logging.Middleware)

	jr := r.PathPrefix("/job").Subrouter()
	jr.HandleFunc("/createJob", createjob.Handler).Methods("POST")
	jr.HandleFunc("/getJob", getjob.Handler).Methods("GET", "POST")

	am := auth.Middleware{}
	am.Populate()
	jr.Use(am.Middleware)
	http.Handle("/", r)

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second * 15,
		"Graceful Shutdown time is 15 seconds",
	)
	flag.Parse()
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", constants.AppPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 120,
		Handler:      r,
	}
	log.Printf("Goatway HTTP Webserver is running on port %s\n", constants.AppPort)
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
