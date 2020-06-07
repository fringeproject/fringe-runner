package threatminer

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type ThreatMiner struct {
}

type ThreatMinerAPIResponse struct {
	Results []string `json:"results"`
}

func NewThreatMiner() *ThreatMiner {
	mod := &ThreatMiner{}

	return mod
}

func (m *ThreatMiner) Name() string {
	return "ThreatMiner Hostname"
}

func (m *ThreatMiner) Slug() string {
	return "threatminer-hostname"
}

func (m *ThreatMiner) Description() string {
	return "Requests ThreatMiner Hostname API. Ref: https://www.threatminer.org/api.php"
}

func (m *ThreatMiner) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *ThreatMiner) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	url := "https://api.threatminer.org/v2/domain.php?q=" + hostname + "&rt=5"
	response := ThreatMinerAPIResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &response, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request ThreatMiner API.")
		logrus.Warn(err)
		return err
	}

	hostnames := []string{}
	for _, host := range response.Results {
		if !common.StringInSlice(hostnames, host) {
			hostnames = append(hostnames, host)

			err = ctx.CreateNewAssetAsHostname(host)
			if err != nil {
				logrus.Info("Something went wrong creating the new asset.")
				logrus.Debug(err)
			}
		}
	}

	return nil
}
