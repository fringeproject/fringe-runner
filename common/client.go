package common

// Representeq a client that request next job and send update to a backend
type RunnerClient interface {
	SendModuleList([]Module) error
	RequestJob() (*Job, error)
	UpdateJob(*Job, []Asset) error
}
