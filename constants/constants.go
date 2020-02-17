package constants

// Port HTTP port where the application is run on
const Port = "9090"

// Job data used to identify and transfer jobs handles by goatway
type Job struct {
	AssetID string
	JobID   string
	Service string
	Status  int
}

// JobReq data request format used when querying job statuses
type JobReq struct {
	JobID string
}

const (
	Pending    int = 0
	Processing int = 1
	Success    int = 2
	Error      int = 3
)