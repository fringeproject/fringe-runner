package github

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type githubSearchResponseItem struct {
	HtmlUrl string `json:"html_url"`
}

type githubSearchResponse struct {
	TotalCount int64                      `json:"total_count"`
	Items      []githubSearchResponseItem `json:"items"`
	Message    string                     `json:"message"`
}

func GithubSearch(ctx *common.ModuleContext, search, sort, order string, page int) (*githubSearchResponse, error) {
	token, err := ctx.GetConfigurationValue("github_api_token")
	if err != nil {
		err := fmt.Errorf("You must provide a Github API token \"github_api_token\".")
		return nil, err
	}

	logrus.Debugf("Search %s [sort:%s][order:%s][page:%d]", search, sort, order, page)
	url := fmt.Sprintf("https://api.github.com/search/code?per_page=100&s=%s&type=Code&o=%s&q=%s&page=%d", sort, order, search, page)

	httpOpts := ctx.GetDefaultHTTPOptions()
	httpOpts.Headers = []common.HTTPHeader{
		{Name: "Authorization", Value: "token " + token},
	}

	res := &githubSearchResponse{}
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, res, httpOpts)
	if err != nil {
		return nil, err
	}

	if res.Message != "" {
		return nil, fmt.Errorf("The Github API returns the following error: \"%s\"", res.Message)
	}

	return res, nil
}

func GithubReadCode(ctx *common.ModuleContext, url string) (string, error) {
	rawURL := strings.Replace(url, "https://github.com/", "https://raw.githubusercontent.com/", 1)
	rawURL = strings.Replace(rawURL, "/blob/", "/", 1)

	logrus.Debugf("Get code: \"%s\"", rawURL)
	_, code, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		return "", err
	}

	return string(*code), nil
}
