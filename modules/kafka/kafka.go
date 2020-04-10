package kafka

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
)

type KafkaAPI struct {
}

func NewKafkaAPI() *KafkaAPI {
	mod := &KafkaAPI{}

	return mod
}

func (m *KafkaAPI) Name() string {
	return "Kafka REST API"
}

func (m *KafkaAPI) Slug() string {
	return "kafka-api"
}

func (m *KafkaAPI) Description() string {
	return "Test if a Kafka REST interface is exposed on port 8083. Ref: https://docs.confluent.io/current/connect/references/restapi.html"
}

func (m *KafkaAPI) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	// Try both HTTP and HTTPS on this interface
	urls := []string{
		"http://" + hostname + ":8083/",
		"https://" + hostname + ":8083/",
	}

	for _, url := range urls {
		statusCode, body, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
		if err != nil {
			logrus.Warnf("Error fetching URL %s", url)
			logrus.Debug(err)

			// We stop the iteration because there is a HTTP error and continue
			// on others URL
			continue
		}

		if *statusCode == http.StatusOK && strings.Contains(string(*body), "kafka_cluster_id") {
			err = ctx.CreateNewAsset("Kafka REST interface is exposed.", assets.AssetTypes["raw"])
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create vulnerability.")
			}
		}
	}

	return nil
}