package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

const (
	RunnerTokenHeaderName  = "X-Runner-Token"
	ContentTypeHeaderName  = "Content-Type"
	AcceptHeaderName       = "Accept"
	ContentTypeHeaderValue = "application/json"
)

type FringeClient struct {
	httpClient  *common.HTTPClient
	coordinator string
	id          string
	token       string
	perimeter   string
}

type FringeClientrModuleListRequest struct {
	Modules []common.Module `json:"modules"`
}

type FringeClientrUpdateJobRequest struct {
	ID          string         `json:"id"`
	Status      string         `json:"status"`
	Assets      []common.Asset `json:"datas"`
	Tags        []string       `json:"tags"`
	Description string         `json:"description"`
	StartedAt   int64          `json:"startedAt"`
	EndedAt     int64          `json:"endedAt"`
}

func NewFringeClient(coordinator string, id string, token string, perimeter string) (common.RunnerClient, error) {

	// Check if the coordinator is a valid URL and add it's IP to the HTTP whitelist
	coordinatorURL, err := url.Parse(coordinator)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("The coordinator URL is not valid.")
		return nil, err
	}

	coordinatorIP, err := net.LookupHost(coordinatorURL.Hostname())
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot resolve coordinator hostname: %s", coordinatorURL.Hostname())
		return nil, err
	}

	proxy, verifyCert := common.GetProxyFromEnv()

	opt := common.HTTPOptions{
		Headers: []common.HTTPHeader{
			{Name: ContentTypeHeaderName, Value: ContentTypeHeaderValue},
			{Name: AcceptHeaderName, Value: ContentTypeHeaderValue},
			{Name: RunnerTokenHeaderName, Value: token},
		},
		Timeout:        time.Second * 20,
		FollowRedirect: true,
		Proxy:          proxy,
		VerifyCert:     verifyCert,
		WhiteListIP:    coordinatorIP,
	}

	httpClient, err := common.NewHTTPClient(context.Background(), &opt)
	if err != nil {
		return nil, err
	}

	client := &FringeClient{
		httpClient:  httpClient,
		coordinator: coordinator,
		id:          id,
		token:       token,
		perimeter:   perimeter,
	}

	return client, nil
}

func (c *FringeClient) String() string {
	return fmt.Sprintf("FringeClient <%s>", c.id)
}

func (c *FringeClient) SendModuleList(modules []common.Module) error {
	url := fmt.Sprintf("%s/perimeters/%s/runners/%s/modules", c.coordinator, c.perimeter, c.id)
	data := &FringeClientrModuleListRequest{
		Modules: modules,
	}

	statusCode, _, _, err := c.httpClient.DoJson(http.MethodPost, url, "", "", data, nil)
	if err != nil {
		return err
	}

	if *statusCode == 404 {
		err = fmt.Errorf("The coordinator does not accept the runner calls.")
		return err
	}

	return err
}

func (c *FringeClient) RequestJob() (*common.Job, error) {
	url := fmt.Sprintf("%s/perimeters/%s/runners/%s/job", c.coordinator, c.perimeter, c.id)
	job := &common.Job{}

	_, _, _, err := c.httpClient.DoJson(http.MethodGet, url, "", "", nil, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *FringeClient) UpdateJob(job *common.Job, newAssets []common.Asset) error {
	url := fmt.Sprintf("%s/perimeters/%s/runners/%s/job", c.coordinator, c.perimeter, c.id)
	data := &FringeClientrUpdateJobRequest{
		ID:          job.ID,
		Status:      "SU",
		Assets:      newAssets,
		Tags:        []string{},
		Description: "",
		StartedAt:   0,
		EndedAt:     0,
	}

	_, _, _, err := c.httpClient.DoJson(http.MethodPut, url, "", "", data, nil)

	return err
}
