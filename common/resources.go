package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
)

const (
	FRINGE_RESSOURCE_DIRECTORY = "ressources"
)

func getResourceDirectory(homeDirectory string) (string, error) {
	resourceDirectory := path.Join(homeDirectory, FRINGE_RESSOURCE_DIRECTORY)

	err := os.MkdirAll(resourceDirectory, 0700)
	if err != nil {
		err = fmt.Errorf("Cannot create the directory for the ressources.")
		return "", err
	}

	return resourceDirectory, nil
}

func GetRessourceFile(config *FringeConfig, filename string) (string, error) {
	resourceDirectory, err := getResourceDirectory(config.HomeDirectory)

	if err != nil {
		return "", nil
	}

	return path.Join(resourceDirectory, filename), nil
}

func DownloadResource(resource ModuleResource, config *FringeConfig) error {
	resourceDirectory, err := getResourceDirectory(config.HomeDirectory)
	if err != nil {
		return err
	}

	opt := HTTPOptions{
		Proxy:      config.Proxy,
		VerifyCert: config.VerifyCert,
	}
	client, err := NewHTTPClient(context.Background(), &opt)
	if err != nil {
		return err
	}

	filePath := path.Join(resourceDirectory, resource.Name)
	statusCode, _, err := client.DownloadFile(http.MethodGet, resource.URL, "", "", filePath, nil)
	if err != nil {
		return err
	}

	if *statusCode != 200 {
		os.Remove(filePath)
		return fmt.Errorf("The server returns the status code %d", *statusCode)
	}

	return nil
}
