package httpprobe

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type HttpProbe struct {
}

func NewHttpProbe() *HttpProbe {
	mod := &HttpProbe{}

	return mod
}

func (m *HttpProbe) Name() string {
	return "HttpProbe"
}

func (m *HttpProbe) Slug() string {
	return "http-probe"
}

func (m *HttpProbe) Description() string {
	return "From a hostname, probe for a HTTP or HTTPS endpoint."
}

func (m *HttpProbe) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *HttpProbe) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	hostname = common.CleanHostname(hostname)

	urls := []string{
		"http://" + hostname,
		"https://" + hostname,
	}
	opts := ctx.GetDefaultHTTPOptions()
	opts.FollowRedirect = false

	for _, url := range urls {
		logrus.Debugf("Requesting: \"%s\"", url)
		_, _, _, err := ctx.HttpRequest(http.MethodGet, url, nil, opts)
		if err != nil {
			logrus.Debug(err)
			continue
		}

		err = ctx.CreateNewAssetAsURL(url)
		if err != nil {
			logrus.Debug(err)
			logrus.Warn("Could not create URL asset.")
		}
	}

	return nil
}
