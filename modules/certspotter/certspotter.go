package certspotter

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type CertSpotter struct {
}

type CertSpotterResponse struct {
	DNSNames []string `json:"dns_names"`
}

func NewCertSpotter() *CertSpotter {
	mod := &CertSpotter{}

	return mod
}

func (m *CertSpotter) Name() string {
	return "Cert Spotter"
}

func (m *CertSpotter) Slug() string {
	return "certspotter"
}

func (m *CertSpotter) Description() string {
	return "Request CertSpotter website to get informations about hostname. Ref: https://sslmate.com/certspotter/"
}

func (m *CertSpotter) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}
	hostname = common.CleanHostname(hostname)

	// url := "https://certspotter.com/api/v0/certs?domain=" + hostname
	url := "https://api.certspotter.com/v1/issuances?domain=" + hostname + "&include_subdomains=true&expand=dns_names"

	var certSpotterResponses []CertSpotterResponse
	statusCode, _, _, err := ctx.HTTPRequestJson(http.MethodGet, url, nil, &certSpotterResponses, nil)

	if err != nil {
		logrus.Debug(err)

		if statusCode != nil && *statusCode == 429 {
			// TODO: use API key, there's nothing in the documentation.
			err = fmt.Errorf("You have exceeded the domain search rate limit for the Cert Spotter API. Please try again later, or authenticate with an API key.")
			logrus.Warn(err)
			return err
		}

		err = fmt.Errorf("Cannot request CertSpotter")
		logrus.Warn(err)
		return err
	}

	var hostnames []string
	for _, certSpotterResponse := range certSpotterResponses {
		for _, host := range certSpotterResponse.DNSNames {
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
	}

	return nil
}
