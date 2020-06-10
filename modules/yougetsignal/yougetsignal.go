package yougetsignal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type YouGetSignal struct {
}

type YouGetSignalDomain []string

type YouGetSignalResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Domains []YouGetSignalDomain `json:"domainArray"`
}

func NewYouGetSignal() *YouGetSignal {
	mod := &YouGetSignal{}

	return mod
}

func (m *YouGetSignal) Name() string {
	return "YouGetSignal"
}

func (m *YouGetSignal) Slug() string {
	return "yougetsignal"
}

func (m *YouGetSignal) Description() string {
	return "Requests yougetsignal.com reverse DNS. Ref: https://www.yougetsignal.com/tools/web-sites-on-web-server/"
}

func (m *YouGetSignal) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *YouGetSignal) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	url := "https://domains.yougetsignal.com/domains.php?remoteAddress=" + hostname
	youGetSignalResponse := &YouGetSignalResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodPost, url, nil, youGetSignalResponse, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request YouGetSignal website.")
		logrus.Warn(err)
		return err
	}

	// There is no results
	if youGetSignalResponse.Status == "Fail" {
		logrus.Debug(youGetSignalResponse.Message)
		return nil
	}

	if youGetSignalResponse.Status != "Success" {
		logrus.Debug(youGetSignalResponse.Message)
		err = fmt.Errorf("YouGetSignal returns an error while fetching the assets.")
		logrus.Warn(err)
		return err
	}

	var hostnames []string
	for _, domain := range youGetSignalResponse.Domains {
		host := domain[0]

		if strings.HasSuffix(host, hostname) && !common.StringInSlice(hostnames, host) {
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
