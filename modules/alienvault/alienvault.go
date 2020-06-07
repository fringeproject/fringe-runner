package alienvault

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type AlienVaultPassiveDNSResponse struct {
	PassiveDNS []AlienVaultPassiveDNS `json:"passive_dns"`
}

type AlienVaultURLListResponse struct {
	HasNext bool            `json:"has_next"`
	URLList []AlienVaultURL `json:"url_list"`
}

type AlienVaultPassiveDNS struct {
	Address  string `json:"address"`
	Hostname string `json:"hostname"`
}

type AlienVaultURL struct {
	Domain   string `json:"domain"`
	URL      string `json:"url"`
	Hostname string `json:"hostname"`
}

type AlienVault struct {
}

func NewAlienVault() *AlienVault {
	mod := &AlienVault{}

	return mod
}

func (m *AlienVault) Name() string {
	return "AlienVault"
}

func (m *AlienVault) Slug() string {
	return "alienvault"
}

func (m *AlienVault) Description() string {
	return "Requests AlienVault API. Ref: https://otx.alienvault.com"
}

func (m *AlienVault) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *AlienVault) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	url := "https://otx.alienvault.com/api/v1/indicators/domain/" + hostname + "/passive_dns"
	passiveDNS := AlienVaultPassiveDNSResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &passiveDNS, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request AlienVault Passive DNS API.")
		logrus.Warn(err)
		return err
	}

	url = "https://otx.alienvault.com/api/v1/indicators/domain/" + hostname + "/url_list"
	urlList := AlienVaultURLListResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &urlList, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request AlienVault Passive DNS API.")
		logrus.Warn(err)
		return err
	}

	hostnames := []string{}
	for _, dns := range passiveDNS.PassiveDNS {
		host := dns.Hostname
		if !common.StringInSlice(hostnames, host) {
			hostnames = append(hostnames, host)

			err = ctx.CreateNewAssetAsHostname(host)
			if err != nil {
				logrus.Info("Something went wrong creating the new asset.")
				logrus.Debug(err)
			}
		}
	}
	for _, url := range urlList.URLList {
		host := url.Hostname
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
