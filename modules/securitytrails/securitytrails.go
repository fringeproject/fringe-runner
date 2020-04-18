package securitytrails

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type securityTrailsSubdomainsResponse struct {
	Subdomains []string `json:"subdomains"`
}

type SecurityTrails struct {
}

func NewSecurityTrails() *SecurityTrails {
	mod := &SecurityTrails{}

	return mod
}

func (m *SecurityTrails) Name() string {
	return "SecurityTrails"
}

func (m *SecurityTrails) Slug() string {
	return "securitytrails"
}

func (m *SecurityTrails) Description() string {
	return "Requests SecurityTrails API. Ref: https://securitytrails.com/corp/api"
}

func (m *SecurityTrails) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	stAPIKey, err := ctx.GetConfigurationValue("SECURITYTRAILS_API_KEY")
	if err != nil {
		err = fmt.Errorf("You must provide a SECURITYTRAILS_API_KEY value to fetch the API.")
		return err
	}

	url := "https://api.securitytrails.com/v1/domain/" + hostname + "/subdomains?apikey=" + stAPIKey

	var stSubdomainsResponse securityTrailsSubdomainsResponse
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &stSubdomainsResponse, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request SecurityTrails API.")
		logrus.Warn(err)
		return err
	}

	for _, subdomain := range stSubdomainsResponse.Subdomains {
		s := subdomain + "." + hostname
		err = ctx.CreateNewAssetAsHostname(s)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
