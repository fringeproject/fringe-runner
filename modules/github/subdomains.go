package github

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type GithubSubdomains struct {
}

func NewGithubSubdomains() *GithubSubdomains {
	mod := &GithubSubdomains{}

	return mod
}

func (m *GithubSubdomains) Name() string {
	return "Github subdomains"
}

func (m *GithubSubdomains) Slug() string {
	return "github-subdomains"
}

func (m *GithubSubdomains) Description() string {
	return "Search the Github codebase for subdomains. Ref: https://developer.github.com/v3/search/#search-code"
}

func (m *GithubSubdomains) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *GithubSubdomains) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	search := fmt.Sprintf("\"%s\"", hostname)
	page := 1
	hostname = "." + hostname
	hostnames := []string{}

	for {
		results, err := GithubSearch(ctx, search, "indexed", "desc", page)
		if err != nil {
			// If there is an error on the first page, then we stop the module
			if page == 1 {
				return err
			} else {
				// Stop the loop if there is an error with the API (rate-limit...)
				break
			}
		}

		// There is no more results
		if len(results.Items) == 0 {
			break
		}

		for _, item := range results.Items {
			code, err := GithubReadCode(ctx, item.HtmlUrl)
			if err != nil {
				logrus.Warn(err)
				continue
			}

			subdomains := common.SearchAllHostname(code)
			for _, host := range subdomains {
				if strings.HasSuffix(host, hostname) && !common.StringInSlice(hostnames, host) {
					hostnames = append(hostnames, host)

					err = ctx.CreateNewAssetAsHostname(host)
					if err != nil {
						logrus.Info("Something went wrong creating the new asset.")
						logrus.Debug(err)
					}
				}
			}
		}

		page += 1
	}

	return nil
}
