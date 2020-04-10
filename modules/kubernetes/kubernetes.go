package kubernetes

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
)

type KubernetesAPI struct {
}

func NewKubernetesAPI() *KubernetesAPI {
	mod := &KubernetesAPI{}

	return mod
}

func (m *KubernetesAPI) Name() string {
	return "Kubernetes API"
}

func (m *KubernetesAPI) Slug() string {
	return "kubernetes-api"
}

func (m *KubernetesAPI) Description() string {
	return "Test if a kubernetes API is exposed on port 10250."
}

func (m *KubernetesAPI) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	urls := []string{
		"https://" + hostname + ":10250/pods",
	}

	for _, url := range urls {
		statusCode, _, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
		if err != nil {
			logrus.Warnf("Error fetching URL %s", url)
			logrus.Debug(err)

			// We stop the iteration because there is a HTTP error and continue
			// on others URL
			continue
		}

		if *statusCode == http.StatusOK {
			err = ctx.CreateNewAsset("Kubernetes API is exposed.", assets.AssetTypes["raw"])
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create vulnerability.")
			}
		}
	}

	return nil
}
