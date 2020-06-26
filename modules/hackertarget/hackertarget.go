package hackertarget

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type HackerTarget struct {
}

func NewHackerTarget() *HackerTarget {
	mod := &HackerTarget{}

	return mod
}

func (m *HackerTarget) Name() string {
	return "HackerTarget"
}

func (m *HackerTarget) Slug() string {
	return "hackertarget"
}

func (m *HackerTarget) Description() string {
	return "Requests HackerTarget API to get informations about hostname."
}

func (m *HackerTarget) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *HackerTarget) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}
	hostname = common.CleanHostname(hostname)

	// https://api.hackertarget.com/aslookup/?q=
	url := "https://api.hackertarget.com/hostsearch/?q=" + hostname
	_, body, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error fetching HackerTarget API.")
		logrus.Warn(err)
		return err
	}

	// The response is a csv of <domain>,<ip>
	r := csv.NewReader(bytes.NewReader(*body))
	records, err := r.ReadAll()
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot parse HackerTarget response.")
		logrus.Warn(err)
		return err
	}

	for _, record := range records {
		hostname = record[0]
		// ip := record[1]

		err = ctx.CreateNewAssetAsHostname(hostname)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
