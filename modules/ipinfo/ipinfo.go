package ipinfo

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type IPInfoResponse struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Loc      string `json:"loc"`
}

type IPInfo struct {
}

func NewIPInfo() *IPInfo {
	mod := &IPInfo{}

	return mod
}

func (m *IPInfo) Name() string {
	return "IP Info"
}

func (m *IPInfo) Slug() string {
	return "ipinfo"
}

func (m *IPInfo) Description() string {
	return "Requests IPInfo API to retreive information about an IP. Ref: https://ipinfo.io/developers"
}

func (m *IPInfo) Run(ctx *common.ModuleContext) error {
	ip, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	IPInfoAPIKey, err := ctx.GetConfigurationValue("ipinfo_api_key")
	if err != nil {
		logrus.Info("IPInfo API key is empty, you will be throttled after few requests.")
		IPInfoAPIKey = ""
	}

	url := "https://ipinfo.io/" + ip + "/json?token=" + IPInfoAPIKey
	ipinfoResponse := &IPInfoResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, ipinfoResponse, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error requesting IPInfo API and parsing results.")
		logrus.Warn(err)
		return err
	}

	if len(ipinfoResponse.Hostname) > 0 {
		err = ctx.CreateNewAssetAsHostname(ipinfoResponse.Hostname)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	if len(ipinfoResponse.Loc) > 0 {
		err = ctx.AddTag("loc:" + ipinfoResponse.Loc)
		if err != nil {
			logrus.Info("Could not add tag.")
			logrus.Debug(err)
		}
	}

	return nil
}
