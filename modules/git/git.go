package git

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

const (
	HTTPGitURI = ".git/HEAD"
)

type HttpGit struct {
}

func NewHttpGit() *HttpGit {
	mod := &HttpGit{}

	return mod
}

func (m *HttpGit) Name() string {
	return "HTTP git"
}

func (m *HttpGit) Slug() string {
	return "http-git"
}

func (m *HttpGit) Description() string {
	return "Test if the web server exposes the .git folder."
}

func (m *HttpGit) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *HttpGit) Run(ctx *common.ModuleContext) error {
	rawURL, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot parse URL")
		logrus.Warn(err)
		return err
	}

	// baseURL.Path = path.Join(baseURL.Path, HTTPGitURI)
	// relativeGitURL := baseURL.String()
	baseURL.Path = path.Join(baseURL.Path, "/"+HTTPGitURI)
	rootGitURL := baseURL.String()

	urls := []string{rootGitURL}
	for _, url := range urls {
		statusCode, bodyBytes, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
		if err != nil {
			logrus.Warn("Error fetching URL ", url)
			logrus.Debug(err)
			continue
		}

		if *statusCode == http.StatusOK && strings.HasPrefix(string(*bodyBytes), "ref:") {
			err = ctx.CreateNewAssetAsRaw("The .git folder is exposed.")
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create vulnerability.")
			}

			err = ctx.AddTag("git")
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not add tag.")
			}
		}
	}

	return nil
}
