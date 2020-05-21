package docker

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type DockerAPI struct {
}

func NewDockerAPI() *DockerAPI {
	mod := &DockerAPI{}

	return mod
}

func (m *DockerAPI) Name() string {
	return "Docker API"
}

func (m *DockerAPI) Slug() string {
	return "docker-api"
}

func (m *DockerAPI) Description() string {
	return "Test if a docker API (dockerd) is exposed on port 2375 and 2376 of the host. Ref: https://docs.docker.com/engine/reference/commandline/dockerd/"
}

func (m *DockerAPI) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	// The doc says:
	//   "It is conventional to use port 2375 for un-encrypted, and port 2376
	//   for encrypted communication with the daemon."
	urls := []string{
		"http://" + hostname + ":2375/version",
		"https://" + hostname + ":2376/version",
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
			err = ctx.CreateNewAssetAsRaw("Docker API is exposed.")
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create vulnerability.")
			}

			err = ctx.AddTag("docker")
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not add tag.")
			}
		}
	}

	return nil
}
