package shodan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Shodan struct {
}

type shodanAPIResponse struct {
	Ports []int `json:"ports"`
}

func NewShodan() *Shodan {
	mod := &Shodan{}

	return mod
}

func (m *Shodan) Name() string {
	return "Shodan"
}

func (m *Shodan) Slug() string {
	return "shodan"
}

func (m *Shodan) Description() string {
	return "Requests Shodan API."
}

func (m *Shodan) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *Shodan) Run(ctx *common.ModuleContext) error {
	ip, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	shodanAPIKey, err := ctx.GetConfigurationValue("shodan_api_key")
	if err != nil {
		err = fmt.Errorf("You must provide a shodan_api_key value to fetch the API.")
		return err
	}

	url := "https://api.shodan.io/shodan/host/" + ip + "?key=" + shodanAPIKey
	_, body, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request Shodan API.")
		logrus.Warn(err)
		return err
	}

	shodanResponse := &shodanAPIResponse{}
	decoder := json.NewDecoder(bytes.NewReader(*body))
	err = decoder.Decode(&shodanResponse)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot decode Shodan response.")
		logrus.Warn(err)
		return err
	}

	for _, port := range shodanResponse.Ports {
		portMsg := fmt.Sprintf("Port %d seems to be open with service.", port)

		err = ctx.CreateNewAssetAsRaw(portMsg)
		if err != nil {
			logrus.Debug(err)
			logrus.Warn("Could not create vulnerability.")
		}
	}

	return nil
}
