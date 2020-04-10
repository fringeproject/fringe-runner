package censys

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
)

type Censys struct {
}

type censysViewResponse struct {
	Tags             []string `json:"tags"`
	Ports            []int    `json:"ports"`
	AutonomousSystem struct {
		Name string `json:"name"`
	} `json:"autonomous_system"`
}

func NewCensys() *Censys {
	mod := &Censys{}

	return mod
}

func (m *Censys) Name() string {
	return "Censys"
}

func (m *Censys) Slug() string {
	return "censys"
}

func (m *Censys) Description() string {
	return "Requests CensysIO API."
}

func (m *Censys) Run(ctx *common.ModuleContext) error {
	ip, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	censysAPIID, err := ctx.GetConfigurationValue("CENSYS_API_ID")
	if err != nil {
		err := fmt.Errorf("You must provide a CENSYS_API_ID value to fetch the API.")
		return err
	}

	CensysAPISecret, err := ctx.GetConfigurationValue("CENSYS_API_SECRET")
	if err != nil {
		err := fmt.Errorf("You must provide a CENSYS_API_SECRET value to fetch the API.")
		return err
	}

	headers := []common.HTTPHeader{
		common.GetBasicAuthHeader(censysAPIID, CensysAPISecret),
	}

	httpOpts := ctx.GetDefaultHTTPOptions()
	httpOpts.Timeout = time.Second * 10
	httpOpts.Headers = headers

	url := "https://censys.io/api/v1/view/ipv4/" + ip
	censysResponse := &censysViewResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, censysResponse, httpOpts)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request Censys API.")
		logrus.Warn(err)
		return err
	}

	// Add the Autonomous System name as a tag for the IP
	err = ctx.CreateNewAsset("tag:"+censysResponse.AutonomousSystem.Name, assets.AssetTypes["raw"])
	if err != nil {
		logrus.Debug(err)
		logrus.Warn("Could not create tag.")
	}

	// Add open ports to the IP
	for _, port := range censysResponse.Ports {
		portMsg := fmt.Sprintf("Port %d seems to be open with service.", port)

		err = ctx.CreateNewAsset(portMsg, assets.AssetTypes["raw"])
		if err != nil {
			logrus.Debug(err)
			logrus.Warn("Could not create vulnerability.")
		}
	}

	return nil
}
