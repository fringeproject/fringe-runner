package tor

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Tor struct {
}

type TorResponse struct {
	// We just need this field for now
	Subdomains []string `json:"subdomains,omitempty"`
}

func NewTor() *Tor {
	mod := &Tor{}

	return mod
}

func (m *Tor) Name() string {
	return "Tor"
}

func (m *Tor) Slug() string {
	return "tor"
}

func (m *Tor) Description() string {
	return "Request Exonerator to check if the IP is an exit node. Ref: https://metrics.torproject.org/exonerator.html"
}

func (m *Tor) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *Tor) Run(ctx *common.ModuleContext) error {
	ip, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	// The latest accepted data is the day before yesterday.
	dt := time.Now().AddDate(0, 0, -2)
	date := dt.Format("2006-01-02")
	url := "https://metrics.torproject.org/exonerator.html?ip=" + ip + "&timestamp=" + date + "&lang=en"

	_, body, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request Exonerator")
		logrus.Warn(err)
		return err
	}

	// The response is like "<h3 class="panel-title">VALUE</h3>"
	bodyString := string(*body)
	h3StartIndex := strings.Index(bodyString, "<h3 class=\"panel-title\">")
	if h3StartIndex == -1 {
		err := fmt.Errorf("Could not parse Exonerator response.")
		return err
	}
	h3StartIndex += 24

	h3EndIndex := strings.Index(bodyString[h3StartIndex:], "</h3>")
	if h3EndIndex == -1 {
		err := fmt.Errorf("Could not parse Exonerator response.")
		return err
	}

	answer := bodyString[h3StartIndex : h3StartIndex+h3EndIndex]
	logrus.Debugf("Exonerator answer is: %s", answer)

	if answer == "Result is positive" {
		err = ctx.AddTag("tor")
		if err != nil {
			logrus.Info("Something went wrong creating the new tag.")
			logrus.Debug(err)
		}
	} else if answer != "Result is negative" {
		err = fmt.Errorf("Exonerator returns an unexpected answer: %s", answer)
		return err
	}

	return nil
}
