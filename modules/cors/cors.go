package cors

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type CORS struct {
}

func NewCORS() *CORS {
	mod := &CORS{}

	return mod
}

func (m *CORS) Name() string {
	return "HTTP CORS"
}

func (m *CORS) Slug() string {
	return "http-cors"
}

func (m *CORS) Description() string {
	return "Check if the server reflect the Origin header in the CORS response."
}

func (m *CORS) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *CORS) Run(ctx *common.ModuleContext) error {
	url, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	beacon := ctx.GetBeaconAsURL()
	otps := ctx.GetDefaultHTTPOptions()
	otps.Headers = append(otps.Headers, common.HTTPHeader{Name: "Origin", Value: beacon})

	_, _, responseHeaders, err := ctx.HttpRequest(http.MethodGet, url, nil, otps)
	if err != nil {
		logrus.Debug(err)
		logrus.Warn("Error fetching URL with custom Orign header.")
	}

	for _, headerValues := range *responseHeaders {
		for _, headerValue := range headerValues {
			if strings.Contains(headerValue, beacon) {
				err = ctx.CreateNewAssetAsRaw("URL may be vulneranle to CORS.")
				if err != nil {
					logrus.Debug(err)
					logrus.Warn("Could not create vulnerability.")
				}
			}
		}
	}

	return nil
}
