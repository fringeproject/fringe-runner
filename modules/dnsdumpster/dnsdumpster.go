package dnsdumpster

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Dnsdumpster struct {
}

func NewDnsdumpster() *Dnsdumpster {
	mod := &Dnsdumpster{}

	return mod
}

func (m *Dnsdumpster) Name() string {
	return "DNSDumpster"
}

func (m *Dnsdumpster) Slug() string {
	return "dnsdumpster"
}

func (m *Dnsdumpster) Description() string {
	return "Requests dnsdumpster.com website to get informations about hostname."
}

func (m *Dnsdumpster) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *Dnsdumpster) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	dnsdumpsterURL := "https://dnsdumpster.com"
	_, body, _, err := ctx.HttpRequest(http.MethodGet, dnsdumpsterURL, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error fetching DNSDumpster URL for csrfToken")
		logrus.Warn(err)
		return err
	}

	csrfRegexp, err := regexp.Compile("<input type=\"hidden\" name=\"csrfmiddlewaretoken\" value=\"(.*?)\">")
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot compile DNSDumpster CSRF regexp")
		logrus.Warn(err)
		return err
	}

	csrfMatches := csrfRegexp.FindStringSubmatch(string(*body))
	if len(csrfMatches) != 2 {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot find CSRF token")
		logrus.Warn(err)
		return err
	}

	csrfToken := csrfMatches[1]

	// Set the URL data
	urlData := url.Values{}
	urlData.Set("csrfmiddlewaretoken", csrfToken)
	urlData.Set("targetip", hostname)
	urlDataBytes := bytes.NewBufferString(urlData.Encode())

	// Get the default option and update timeout and headers fields
	httpOpts := ctx.GetDefaultHTTPOptions()
	httpOpts.Timeout = time.Second * 20
	httpOpts.Headers = []common.HTTPHeader{
		{Name: "Referer", Value: "https://dnsdumpster.com"},
		{Name: "Cookie", Value: "csrftoken=" + csrfToken},
		{Name: "Content-Type", Value: "application/x-www-form-urlencoded"},
	}

	statusCode, body, _, err := ctx.HttpRequest(http.MethodPost, dnsdumpsterURL, urlDataBytes, httpOpts)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot get DNSDumpster result")
		logrus.Warn(err)
		return err
	}

	if *statusCode != http.StatusOK {
		logrus.Debug(err)
		err = fmt.Errorf("DNSDumpster server returns status code: %d", *statusCode)
		logrus.Warn(err)
		return err
	}

	// Check if we have been throttled
	ipLimitRegexp, err := regexp.Compile(`Too many requests from your IP address`)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot compile IP limit regexp")
		logrus.Warn(err)
		return err
	}
	ipLimitMatches := ipLimitRegexp.FindAllStringSubmatch(string(*body), -1)
	if len(ipLimitMatches) > 0 {
		err = fmt.Errorf("The IP is temporarily throttled from DNSDumpster, please use a pro account.")
		logrus.Warn(err)
		return err
	}

	tableRegexp, err := regexp.Compile(`(?s)<table.*?>(.*?)</table>`)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot compile table regexp")
		logrus.Warn(err)
		return err
	}

	tableMatches := tableRegexp.FindAllStringSubmatch(string(*body), -1)
	if len(tableMatches) > 2 && len(tableMatches[3]) != 2 {
		logrus.Debug(err)
		err = fmt.Errorf("Can't find tables in response")
		logrus.Warn(err)
		return err
	}

	recordTable := tableMatches[3][1]
	recordRegexp, err := regexp.Compile(`(?s)<td class="col-md-4">(.*?)<br>`)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot compile record regexp")
		logrus.Warn(err)
		return err
	}

	recordMatches := recordRegexp.FindAllStringSubmatch(recordTable, -1)
	for _, record := range recordMatches {
		host := record[1]
		err = ctx.CreateNewAssetAsHostname(host)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
