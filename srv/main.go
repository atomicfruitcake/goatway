package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/atomicfruitcake/goatway/client"
	"github.com/atomicfruitcake/goatway/redis"

	"github.com/gorilla/mux"
)

const port = "9090"

type Job struct {
	AssetId string
	JobId   string
	Service string
	Status  int
}

type JobReq struct {
	JobId string
}

const (
	Pending    int = 0
	Processing int = 1
	Success    int = 2
	Error      int = 3
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("http://0.0.0.0:%s/health", port), 302)

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	var j Job
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body due to %v", err)
		log.Print(msg)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j.Status = Processing
	bytes, err := json.Marshal(j)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling request body due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = client.SendJob(j.Service, bytes)
	if err != nil {
		msg := fmt.Sprintf("Error sending job to service due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = redis.Set(j.JobId, bytes)
	if err != nil {
		msg := fmt.Sprintf("Error saving job status due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	var jr JobReq
	err := json.NewDecoder(r.Body).Decode(&jr)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body due to %v", err)
		log.Print(msg)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	job, err := redis.Get(jr.JobId)
	if err != nil {
		msg := fmt.Sprintf("Error finding job %s due to %v", jr.JobId, err)
		log.Print(msg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

type authenticationMiddleware struct {
	tokenUsers map[string]string
}

var amw *authenticationMiddleware

func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["thisIsAnExampleUserToken"] = "token123"
	amw.tokenUsers["adminToken"] = "admin"
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("adminToken")
		if user, found := amw.tokenUsers[token]; found {
			log.Printf("Successfully Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Token %s is not a valid adminToken", token)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {
	redis.Ping()
	fmt.Println("Starting a new Goatway HTTP Server")
	fmt.Println("Building the Gorilla MUX Router")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/createJob", createJobHandler).Methods("POST")
	router.HandleFunc("/getJob", getJobHandler).Methods("GET", "POST")
	http.Handle("/", router)

	// Not applying auth since it will apply to /health and block Load Balancer health checks

	// amw := authenticationMiddleware{}
	// amw.Populate()
	// router.Use(amw.Middleware)

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"Graceful Shutdown time is 15 seconds",
	)
	flag.Parse()
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	fmt.Printf("Goatway HTTP Webserver is running on port %s\n", port)
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
