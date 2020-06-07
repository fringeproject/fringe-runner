package takeover

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type ServiceProvider struct {
	Name        string   `json:"name"`
	CName       []string `json:"cname"`
	Fingerprint []string `json:"fingerprint"`
}

type TakeOver struct {
}

func NewTakeOver() *TakeOver {
	mod := &TakeOver{}

	return mod
}

func (m *TakeOver) Name() string {
	return "Subdomain takeover"
}

func (m *TakeOver) Slug() string {
	return "takeover"
}

func (m *TakeOver) Description() string {
	return "Checks if the hostname has a dangling CNAME record pointing to a service."
}

func (m *TakeOver) ResourceURLs() []common.ModuleResource {
	return []common.ModuleResource{
		{Name: "takeover_providers.json", URL: "https://static.fringeproject.com/fringe-runner/resources/takeover_providers.json"},
	}
}

func (m *TakeOver) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	providerFile, err := ctx.GetRessourceFile("takeover_providers.json")
	if err != nil {
		return err
	}

	providers := []ServiceProvider{}
	err = common.ReadJSONFile(providerFile, &providers)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot read providers JSON file. Please, check the ressource \"takeover_providers.json\" file.")
		logrus.Warn(err)
		return err
	}

	dnsServer, err := ctx.GetConfigurationValue("dns_server")
	if err != nil {
		logrus.Warn("No DNS server set, use 8.8.8.8 as a default value.")
		dnsServer = "8.8.8.8"
	}

	// TODO: Add IP resolution, MX (mail) and nameserver (CDN, ...)
	cnames, err := common.LookupCName(hostname, dnsServer)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot get CNAME record.")
		logrus.Warn(err)
		return err
	}

	// Request the hostname to fringerprint it
	responseString := ""
	_, body, _, err := ctx.HttpRequest(http.MethodGet, "https://"+hostname, nil, nil)
	if err != nil {
		logrus.Info("Error fetching hostname.")
	} else {
		responseString = string(*body)
	}

	for _, provider := range providers {
		for _, providerCNAME := range provider.CName {
			cnameRegExp, err := regexp.Compile(providerCNAME)
			if err != nil {
				logrus.Infof("Cannot compile Regexp: %s", providerCNAME)
				continue
			}

			for _, cname := range cnames {
				cnameMatches := cnameRegExp.FindStringSubmatch(cname)

				// Check if we have a match
				if len(cnameMatches) > 0 {
					for _, fringerprint := range provider.Fingerprint {
						if strings.Contains(responseString, fringerprint) {
							err := ctx.CreateNewAssetAsRaw("The hostname is vulnerable to DNS TakeOver on " + provider.Name)
							if err != nil {
								logrus.Info("Could not create vulnerability.")
							}

							err = ctx.AddTag("take-over")
							if err != nil {
								logrus.Info("Could not add tag.")
							}
						}
					}
				}
			}
		}
	}

	return nil
}
