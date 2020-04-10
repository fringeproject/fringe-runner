package elasticsearch

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
)

type ElasticSearchAPI struct {
}

func NewElasticSearchAPI() *ElasticSearchAPI {
	mod := &ElasticSearchAPI{}

	return mod
}

func (m *ElasticSearchAPI) Name() string {
	return "Elasticsearch API"
}

func (m *ElasticSearchAPI) Slug() string {
	return "elasticsearch-api"
}

func (m *ElasticSearchAPI) Description() string {
	return "Test if an elasticsearch API is exposed on port 9200 or 9300. Ref: https://www.elastic.co/guide/en/elasticsearch/reference/current/docs.html"
}

func (m *ElasticSearchAPI) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	urls := []string{
		"http://" + hostname + ":9200/_nodes",
		"https://" + hostname + ":9200/_nodes",
	}

	for _, url := range urls {
		statusCode, _, headers, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
		if err != nil {
			logrus.Warnf("Error fetching URL %s", url)
			logrus.Debug(err)

			// We stop the iteration because there is a HTTP error and continue
			// on others URL
			continue
		}

		if *statusCode == http.StatusOK {
			err = ctx.CreateNewAsset("An unauthenticated ElasticSearch database is exposed.", assets.AssetTypes["raw"])
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create vulnerability.")
			}
		} else if *statusCode == http.StatusUnauthorized {
			authenticateHeader := strings.ToLower((*headers).Get("WWW-Authenticate"))

			if strings.Contains(authenticateHeader, "elasticsearch") {
				err = ctx.CreateNewAsset("An authenticated ElasticSearch database is exposed.", assets.AssetTypes["raw"])
				if err != nil {
					logrus.Debug(err)
					logrus.Warn("Could not create vulnerability.")
				}
			}
		}
	}

	return nil
}
