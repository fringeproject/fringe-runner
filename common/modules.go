package common

type ModuleInterface interface {
	Name() string
	Slug() string
	Description() string

	Run() error
}

type Module struct {
	Name        string `json:"displayName"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
