package constants

// AppPort HTTP port where the application is run on
const AppPort = "9090"

// RedisURL URL where the redis-server can be found
const RedisURL = "localhost:6379"

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