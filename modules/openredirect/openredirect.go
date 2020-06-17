package openredirect

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type OpenRedirect struct {
}

func NewOpenRedirect() *OpenRedirect {
	mod := &OpenRedirect{}

	return mod
}

func (m *OpenRedirect) Name() string {
	return "Open Redirect"
}

func (m *OpenRedirect) Slug() string {
	return "http-open-redirect"
}

func (m *OpenRedirect) Description() string {
	return "Check if the server is vulnerable to open-redirect."
}

func (m *OpenRedirect) ResourceURLs() []common.ModuleResource {
	return []common.ModuleResource{
		{Name: "open_redirect.json", URL: "https://static.fringeproject.com/fringe-runner/resources/open_redirect.json"},
	}
}

func AddArgToURL(rawurl, key, value string) (string, error) {
	parseurl, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	queries := parseurl.Query()
	queries.Set(key, value)

	parseurl.RawQuery = queries.Encode()

	return parseurl.String(), nil
}

func (m *OpenRedirect) Run(ctx *common.ModuleContext) error {
	rawurl, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	ressourceFile, err := ctx.GetRessourceFile("open_redirect.json")
	if err != nil {
		return err
	}

	parameters := []string{}
	err = common.ReadJSONFile(ressourceFile, &parameters)
	if err != nil {
		return err
	}

	beacon := ctx.GetBeaconAsURL()

	for _, parameter := range parameters {
		newurl, err := AddArgToURL(rawurl, parameter, beacon)
		if err != nil {
			logrus.Info("Cannot add parameter to the url.")
			logrus.Debug(err)
			continue
		}

		_, _, responseHeaders, err := ctx.HttpRequest(http.MethodGet, newurl, nil, nil)
		if err != nil {
			logrus.Info("Cannot request url")
			logrus.Debug(err)
			continue
		}

		for _, headerValues := range *responseHeaders {
			for _, headerValue := range headerValues {
				if strings.Contains(headerValue, beacon) {
					rawAsset := fmt.Sprintf("URL may be vulneranle to Open-Redirect on paramter %s.", parameter)
					err = ctx.CreateNewAssetAsRaw(rawAsset)
					if err != nil {
						logrus.Debug(err)
						logrus.Warn("Could not create vulnerability.")
					}
				}
			}
		}
	}

	return nil
}
