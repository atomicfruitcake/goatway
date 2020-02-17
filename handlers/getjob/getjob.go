package getjob

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"

	"github.com/atomicfruitcake/goatway/constants"
	"github.com/atomicfruitcake/goatway/redis"
)

// Handler HTTP Request handler for getting the status of processing jobs
func Handler(w http.ResponseWriter, r *http.Request) {
	var jr constants.JobReq
	err := json.NewDecoder(r.Body).Decode(&jr)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body due to %v", err)
		log.Print(msg)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	job, err := redis.Get(jr.JobID)
	if err != nil {
		msg := fmt.Sprintf("Error finding job %s due to %v", jr.JobID, err)
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