package threatcrowd

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type ThreatCrowd struct {
}

type ThreatcrowdResponse struct {
	// We just need this field for now
	Subdomains []string `json:"subdomains,omitempty"`
}

func NewThreatCrowd() *ThreatCrowd {
	mod := &ThreatCrowd{}

	return mod
}

func (m *ThreatCrowd) Name() string {
	return "ThreatCrowd"
}

func (m *ThreatCrowd) Slug() string {
	return "threatcrowd"
}

func (m *ThreatCrowd) Description() string {
	return "Request Threatcrowd website to get informations about hostname. Ref: https://www.threatcrowd.org/"
}

func (m *ThreatCrowd) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *ThreatCrowd) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}
	hostname = common.CleanHostname(hostname)

	url := "https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=" + hostname

	var threatcrowdResponse ThreatcrowdResponse
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &threatcrowdResponse, nil)

	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request ThreatCrowd")
		logrus.Warn(err)
		return err
	}

	var hostnames []string
	for _, host := range threatcrowdResponse.Subdomains {
		host = common.CleanHostname(host)

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
