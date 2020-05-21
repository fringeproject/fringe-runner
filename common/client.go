package common

type FringeClientModuleListRequest struct {
	Modules []Module `json:"modules"`
}

type FringeClientUpdateJobRequest struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	Assets      []Asset  `json:"datas"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	StartedAt   int64    `json:"startedAt"`
	EndedAt     int64    `json:"endedAt"`
}

const (
	JOB_STATUS_WAITING  string = "WA"
	JOB_STATUS_ON_GOING string = "OG"
	JOB_STATUS_SUCCESS  string = "SU"
	JOB_STATUS_ERROR    string = "ER"
)

// Representeq a client that request next job and send update to a backend
type RunnerClient interface {
	SendModuleList([]Module) error
	RequestJob() (*Job, error)
	UpdateJob(*Job, []Asset) error
}
