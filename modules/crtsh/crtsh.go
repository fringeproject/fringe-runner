package crtsh

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Crtsh struct {
}

type CrtShObject struct {
	// We just need this field
	NameValue string `json:"name_value"`
}

func NewCrtsh() *Crtsh {
	mod := &Crtsh{}

	return mod
}

func (m *Crtsh) Name() string {
	return "crt.sh"
}

func (m *Crtsh) Slug() string {
	return "crtsh"
}

func (m *Crtsh) Description() string {
	return "Request crt.sh website to get informations about a hostname."
}

func (m *Crtsh) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *Crtsh) Run(ctx *common.ModuleContext) error {
	// Get the hostname
	baseHostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	// Forge the crt.sh URL to fetch certificat informations
	crtshURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", baseHostname)
	var crtshResponse []CrtShObject
	_, _, _, err = ctx.HTTPRequestJson("GET", crtshURL, nil, &crtshResponse, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error requesting crt.sh and parsing results.")
		logrus.Warn(err)
		return err
	}

	// Parse the results to filter the results
	var hostnames []string
	for _, obj := range crtshResponse {
		// The `name_value` may contain multiple hostnames separated by a new line
		hosts := strings.Split(obj.NameValue, "\n")

		for _, host := range hosts {
			host = strings.TrimPrefix(host, "*.")

			// A certificate can contains other domains so we limit to the asset
			// pass as argument
			if strings.Contains(host, "@") {
				parts := strings.Split(host, "@")
				host = parts[len(parts)-1]
			}

			if strings.HasSuffix(host, baseHostname) && !common.StringInSlice(hostnames, host) {
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
