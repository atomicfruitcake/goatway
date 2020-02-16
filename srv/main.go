package main

import (
	"flag"
	"fmt"
	"io/ioutil"
    "log"
    "net/http"
	"time"
	"encoding/json"

	"github.com/atomicfruitcake/goatway/redis"
	"github.com/atomicfruitcake/goatway/client"

	"github.com/gorilla/mux"
)


type Job struct {
    AssetId string
	JobId   string
	Service string
	Status  int
}

const (
	Pending     int = 0
	Processing  int = 1
	Success   	int = 2
	Error 		int = 3
 )

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://0.0.0.0:9090/health", 302)
	
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}


func createJobHandler(w http.ResponseWriter, r *http.Request) {
	var j Job
    err := json.NewDecoder(r.Body).Decode(&j)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	j.Status = Processing
	bytes, err := json.Marshal(j)
    if err != nil {
        msg := fmt.Sprintf("Error reading body due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	err = redis.Set(j.JobId, bytes)
	if err != nil {
		msg := fmt.Sprintf("Error saving job status due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)
	w.Write([]byte("FOOBAR"))
	json.NewEncoder(w).Encode("")
  }

type authenticationMiddleware struct {

	tokenUsers map[string]string
}

var amw *authenticationMiddleware

func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["thisIsAnExampleUserToken"] = "user1"
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

	// amw := authenticationMiddleware{}
	// amw.Populate()
	// router.Use(amw.Middleware)
	

	var wait time.Duration
    flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second * 15,
		"Graceful Shutdown time is 15 seconds",
	)
	flag.Parse()
	
	srv := &http.Server{
        Addr:         "0.0.0.0:9090",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: router, 
	}

	fmt.Println("Goatway HTTP Webserver is on port 9090")

	srv.ListenAndServe()
	// go func() {
	// 	http.ListenAndServe(":9090", nil)
    //     if err := srv.ListenAndServe(); err != nil {
	// 		fmt.Println("Server crashed")
	// 		fmt.Println(err)
	// 		log.Fatal(err)
    //     }
	// }()
}
