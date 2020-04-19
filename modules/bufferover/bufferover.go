package bufferover

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type bufferOverResponse struct {
	FDNS_A []string `json:"FDNS_A"`
}

type BufferOver struct {
}

func NewBufferOver() *BufferOver {
	mod := &BufferOver{}

	return mod
}

func (m *BufferOver) Name() string {
	return "BufferOver"
}

func (m *BufferOver) Slug() string {
	return "bufferover"
}

func (m *BufferOver) Description() string {
	return "Requests BufferOver API. Ref: https://blog.erbbysam.com/index.php/2019/02/09/dnsgrep/. Ref: https://github.com/erbbysam/DNSGrep"
}

func (m *BufferOver) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	// Add a dot on the query to look up for subdomains
	url := "http://dns.bufferover.run/dns?q=." + hostname
	var res bufferOverResponse
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &res, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request BufferOver API.")
		logrus.Warn(err)
		return err
	}

	for _, fdns := range res.FDNS_A {
		// fdns line is <ip>,<hostname>
		parts := strings.Split(fdns, ",")

		if len(parts) != 2 {
			logrus.Info("Something went wrong parsing BufferOver FDNS response.")
		} else {
			err = ctx.CreateNewAssetAsHostname(parts[1])

			if err != nil {
				logrus.Info("Something went wrong creating the new asset.")
				logrus.Debug(err)
			}
		}
	}

	return nil
}
