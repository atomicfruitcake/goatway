package createjob

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"

	"github.com/atomicfruitcake/goatway/client"
	"github.com/atomicfruitcake/goatway/redis"
	"github.com/atomicfruitcake/goatway/constants"
)


// Handler HTTP Request handler for create new jobs for requests
func Handler(w http.ResponseWriter, r *http.Request) {
	var j constants.Job
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body due to %v", err)
		log.Print(msg)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j.Status = constants.Processing
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

	err = redis.Set(j.JobID, bytes)
	if err != nil {
		msg := fmt.Sprintf("Error saving job status due to %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}