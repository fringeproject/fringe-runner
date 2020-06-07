package common

type ModuleResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ModuleInterface interface {
	Name() string
	Slug() string
	Description() string
	ResourceURLs() []ModuleResource

	Run(*ModuleContext) error
}

type Module struct {
	Name        string `json:"displayName"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
