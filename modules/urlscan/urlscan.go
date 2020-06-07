package urlscan

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type URLScan struct {
}

// Simples structures to get domains from results
type URLScanSearch struct {
	Results []URLScanSearchResult `json:"results"`
}

type URLScanSearchResult struct {
	Page URLScanSearchPage `json:"page"`
}

type URLScanSearchPage struct {
	Domain string `json:"domain"`
}

func NewURLScan() *URLScan {
	mod := &URLScan{}

	return mod
}

func (m *URLScan) Name() string {
	return "urlscan.io"
}

func (m *URLScan) Slug() string {
	return "urlscan"
}

func (m *URLScan) Description() string {
	return "Use the urlscan.io search API to get list of subdomains from URL. Ref: https://urlscan.io/about-api/."
}

func (m *URLScan) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *URLScan) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}
	hostname = common.CleanHostname(hostname)

	url := "https://urlscan.io/api/v1/search/?q=domain:" + hostname
	urlscanSearchResults := &URLScanSearch{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, urlscanSearchResults, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error requesting urlscan Search API and parsing results.")
		logrus.Warn(err)
		return err
	}

	var hostnames []string
	for _, result := range urlscanSearchResults.Results {
		host := common.CleanHostname(result.Page.Domain)

		// The API returns various results
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
