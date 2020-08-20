package nessus

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type NessusAPI struct {
	Endpoint  string
	APIToken  string
	UserToken string
	Context   *common.ModuleContext
}

type NessusAPIError struct {
	Error string `json:"error"`
}

type NessusSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NessusSessionResponse struct {
	Token string `json:"token"`
}

type NessusScanSettings struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TextTargets string `json:"text_targets"`
	LaunchNow   bool   `json:"launch_now"`
	Enabled     bool   `json:"enabled"`
}

type NessusScanRequest struct {
	UUID     string             `json:"uuid"`
	Settings NessusScanSettings `json:"settings"`
}

type NessusScan struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type NessusScanResponse struct {
	Scan NessusScan `json:"scan"`
	NessusAPIError
}

type NessusScansResponse struct {
	Scans []NessusScan `json:"scans"`
	NessusAPIError
}

type NessusTemplate struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type NessusTemplateResponse struct {
	Templates []NessusTemplate `json:"templates"`
	NessusAPIError
}

func checkResponseError(res interface{}) error {
	v := reflect.ValueOf(res)
	r := reflect.Indirect(v)
	f := r.FieldByName("Error")

	if f.IsValid() {
		errorMessage := f.String()

		if errorMessage != "" {
			return fmt.Errorf("The API returns the following error: \"%s\"", errorMessage)
		}
	}

	return nil
}

func NewNessusAPI(ctx *common.ModuleContext, endpoint string, username string, password string) (*NessusAPI, error) {
	api := &NessusAPI{
		Endpoint:  endpoint,
		APIToken:  "",
		UserToken: "",
		Context:   ctx,
	}

	err := api.getAPIToken()
	if err != nil {
		return nil, fmt.Errorf("Cannot fetch and parse Nessus API Token.")
	}

	err = api.Authenticate(username, password)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func (api *NessusAPI) generateURL(nessusPath string) string {
	url, err := url.Parse(api.Endpoint)
	if err != nil {
		fmt.Println("Cannot parse Nessus endpoit")
	}

	url.Path = path.Join(url.Path, nessusPath)

	return url.String()
}

// Send a JSON requests to the API with the authentication headers
func (api *NessusAPI) doJSON(method string, nessusPath string, request interface{}, response interface{}) (*http.Response, error) {
	url := api.generateURL(nessusPath)
	opts := api.Context.GetDefaultHTTPOptions()
	opts.Headers = append(opts.Headers, common.HTTPHeader{Name: "Content-Type", Value: "application/json"})

	if api.APIToken != "" {
		opts.Headers = append(opts.Headers, common.HTTPHeader{Name: "X-API-Token", Value: api.APIToken})
	}

	if api.UserToken != "" {
		opts.Headers = append(opts.Headers, common.HTTPHeader{Name: "X-Cookie", Value: fmt.Sprintf("token=%s", api.UserToken)})
	}

	_, _, _, err := api.Context.HTTPRequestJson(method, url, request, &response, opts)

	if err != nil {
		return nil, err
	}

	// Check if the API returns an error
	err = checkResponseError(response)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Fetch the API token within the large Javascript file
func (api *NessusAPI) getAPIToken() error {
	url := api.generateURL("/nessus6.js")
	_, body, _, err := api.Context.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		fmt.Println("Cannot fetch Nessus Javascript file.")
		return err
	}

	tokenRegexp, err := regexp.Compile(`key:"getApiToken",value:function\(\){return"(.*?)"}`)
	if err != nil {
		fmt.Println("Cannot compile token regexp.")
		return err
	}

	tokenMatches := tokenRegexp.FindStringSubmatch(string(*body))
	if len(tokenMatches) != 2 {
		fmt.Println("Cannot find token in Javascript file.")
		return err
	}

	token := tokenMatches[1]
	api.APIToken = token

	return nil
}

// Authenticate the user to the Nessus API
func (api *NessusAPI) Authenticate(login string, password string) error {
	session := NessusSessionRequest{Username: login, Password: password}
	var sessionResponse NessusSessionResponse
	_, err := api.doJSON(http.MethodPost, "/session", &session, &sessionResponse)

	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("The Nessus credentials are invalid.")
		return err
	}

	api.UserToken = sessionResponse.Token

	return nil
}

// List available templates
func (api *NessusAPI) ListScanTemplates() (map[string]string, error) {
	res := &NessusTemplateResponse{}

	_, err := api.doJSON(http.MethodGet, "/editor/scan/templates", nil, res)
	if err != nil {
		return nil, err
	}

	templates := map[string]string{}
	for _, template := range res.Templates {
		templates[template.Name] = template.UUID
	}

	return templates, nil
}

// Create a new scan using the template UID
func (api *NessusAPI) CreateScan(template string, target string, scanName string) (string, error) {
	description := "Scan created by fringe-runner."

	req := &NessusScanRequest{
		UUID: template,
		Settings: NessusScanSettings{
			Name:        scanName,
			Description: description,
			TextTargets: target,
			LaunchNow:   true,
			Enabled:     false,
		},
	}
	res := &NessusScanResponse{}

	_, err := api.doJSON(http.MethodPost, "/scans", req, res)

	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("An error occured while creating the new Nessus scan.")
		return "", err
	}

	return res.Scan.UUID, nil
}

func (api *NessusAPI) ListScans() ([]NessusScan, error) {
	res := &NessusScansResponse{}
	_, err := api.doJSON(http.MethodGet, "/scans", nil, res)

	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("An error occured while fetching the Nessus scans.")
		return nil, err
	}

	return res.Scans, nil
}
