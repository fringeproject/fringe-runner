package wayback

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Wayback struct {
}

func NewWayback() *Wayback {
	mod := &Wayback{}

	return mod
}

func (m *Wayback) Name() string {
	return "Wayback"
}

func (m *Wayback) Slug() string {
	return "wayback"
}

func (m *Wayback) Description() string {
	return "Request Wayback website to retreive old information about hostname."
}

func RequestPage(ctx *common.ModuleContext, hostname string, page int) ([][]string, error) {
	url := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&collapse=urlkey&page=%d", hostname, page)

	var waybackResponse [][]string
	_, _, _, err := ctx.HTTPRequestJson("GET", url, nil, &waybackResponse, nil)
	if err != nil {
		logrus.Warn("Cannot request Wayback.")
		return nil, err
	}

	// TODO: use the header to count pages
	// numPagesHeader := resp.Header["X-Cdx-Num-Pages"]
	// if len(numPagesHeader) > 0 {
	// 	numPages := resp.Header["X-Cdx-Num-Pages"][0]
	// 	if page == numPages {
	// 	}
	// }

	return waybackResponse, nil
}

func (m *Wayback) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	hostname = common.CleanHostname(hostname)
	hostnames := make([]string, 0)
	page := 0

	for {
		waybackResponse, err := RequestPage(ctx, hostname, page)
		if err != nil {
			break
		}

		// If there's not results then we stop the loop
		if len(waybackResponse) == 0 {
			break
		}

		page += 1
		skip := true
		for _, item := range waybackResponse {
			// the first line is a header
			if skip {
				skip = false
				continue
			}

			if len(item) < 3 {
				continue
			}

			u, err := url.Parse(item[2])
			if err != nil {
				continue
			}

			host := common.CleanHostname(u.Hostname())

			if !common.StringInSlice(hostnames, host) {
				hostnames = append(hostnames, host)

				err = ctx.CreateNewAssetAsHostname(host)
				if err != nil {
					logrus.Info("Something went wrong creating the new asset.")
					logrus.Debug(err)
				}
			}
		}
	}

	return nil
}
