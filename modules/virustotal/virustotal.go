package virustotal

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type VirusTotal struct {
}

type VirusTotalSubdomains struct {
	Data []VirusTotalSubdomainData `json:"data"`
}

type VirusTotalSubdomainData struct {
	Id string `json:"id"`
}

func NewVirusTotal() *VirusTotal {
	mod := &VirusTotal{}

	return mod
}

func (m *VirusTotal) Name() string {
	return "VirusTotal"
}

func (m *VirusTotal) Slug() string {
	return "virustotal"
}

func (m *VirusTotal) Description() string {
	return "Requests VirusTotal API (UI). Ref: https://www.virustotal.com/"
}

func (m *VirusTotal) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	url := "https://www.virustotal.com/ui/domains/" + hostname + "/subdomains"
	resp := &VirusTotalSubdomains{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &resp, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request VirusTotal API (UI).")
		logrus.Warn(err)
		return err
	}

	for _, data := range resp.Data {
		host := common.CleanHostname(data.Id)
		err = ctx.CreateNewAssetAsHostname(host)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
