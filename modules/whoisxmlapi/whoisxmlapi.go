package whoisxmlapi

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type WhoisXMLAPI struct {
}

type whoisXMLAPIViewResult struct {
	Name string `json:"name"`
	// FirstSeen string `json:"first_seen"`
	// LastSeen  string `json:"last_visit"`
}

type whoisXMLAPIViewResponse struct {
	Results     []whoisXMLAPIViewResult `json:"result"`
	CurrentPage string                  `json:"current_page"`
	Size        int                     `json:"size"`
}

func NewWhoisXMLAPI() *WhoisXMLAPI {
	mod := &WhoisXMLAPI{}

	return mod
}

func (m *WhoisXMLAPI) Name() string {
	return "WhoisXMLAPI"
}

func (m *WhoisXMLAPI) Slug() string {
	return "whoisxmlapi"
}

func (m *WhoisXMLAPI) Description() string {
	return "Requests WhoisXMLAPI API. Ref: https://reverse-ip.whoisxmlapi.com/api/documentation/making-requests"
}

func (m *WhoisXMLAPI) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *WhoisXMLAPI) Run(ctx *common.ModuleContext) error {
	ip, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	whoisxmlapiKey, err := ctx.GetConfigurationValue("whoisxmlapi_key")
	if err != nil {
		err := fmt.Errorf("You must provide a whoisxmlapi_key value to fetch the API.")
		return err
	}

	// The API provide a pager based on the last domain.
	// We don't want to waste all the API call so we limit to 3 calls per run
	lastDomain := "1"
	for i := 0; i < 3; i++ {
		url := "https://reverse-ip.whoisxmlapi.com/api/v1?apiKey=" + whoisxmlapiKey + "&ip=" + ip + "&outputFormat=json&from=" + lastDomain
		res := &whoisXMLAPIViewResponse{}
		_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, res, nil)
		if err != nil {
			logrus.Debug(err)
			err = fmt.Errorf("Cannot request WhoisXMLAPI API.")
			logrus.Warn(err)
			return err
		}

		if res.Size == 0 {
			break
		}

		for _, result := range res.Results {
			err = ctx.CreateNewAssetAsHostname(result.Name)
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create hostname.")
			}
		}

		lastDomain = res.Results[len(res.Results)-1].Name

	}

	return nil
}
