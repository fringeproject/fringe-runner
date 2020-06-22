package cloudflare

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type CloudflareDOHResponseAnswer struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type CloudflareDOHResponse struct {
	Answer []CloudflareDOHResponseAnswer `json:"Answer"`
}

type CloudflareDOH struct {
}

func NewCloudflareDOH() *CloudflareDOH {
	mod := &CloudflareDOH{}

	return mod
}

func (m *CloudflareDOH) Name() string {
	return "Cloudflare DOH"
}

func (m *CloudflareDOH) Slug() string {
	return "cloudflare-doh"
}

func (m *CloudflareDOH) Description() string {
	return "Requests Cloudflare DOH. Ref: https://www.cloudflare.com/dns/"
}

func (m *CloudflareDOH) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *CloudflareDOH) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	httpOpts := ctx.GetDefaultHTTPOptions()
	httpOpts.Headers = []common.HTTPHeader{
		{Name: "Accept", Value: "application/dns-json"},
	}

	url := "https://cloudflare-dns.com/dns-query?name=" + hostname + "&type=A"
	res := &CloudflareDOHResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, res, httpOpts)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request Cloudflare DOH DNS.")
		logrus.Warn(err)
		return err
	}

	for _, answer := range res.Answer {
		err = ctx.CreateNewAssetAsIP(answer.Data)

		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
