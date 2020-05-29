package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	FRINGE_RESSOURCE_DIRECTORY = "ressources"
	FRINGE_RESSOURCE_FILENAMES = "takeover_providers.json, wappalyzer.json"
	FRINGE_RESSOURCE_URL       = "https://static.fringeproject.com/fringe-runner/ressources/"
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

func getRessourceFilenames() []string {
	return strings.Split(FRINGE_RESSOURCE_FILENAMES, ", ")
}

func updateModuleRessource(filename string, ressourceFolder string, config *FringeConfig) error {
	url := FRINGE_RESSOURCE_URL + filename
	filePath := path.Join(ressourceFolder, filename)

	opt := HTTPOptions{
		Proxy:      config.Proxy,
		VerifyCert: config.VerifyCert,
	}
	client, err := NewHTTPClient(context.Background(), &opt)
	if err != nil {
		return err
	}

	statusCode, _, err := client.DownloadFile(http.MethodGet, url, "", "", filePath, nil)
	if err != nil {
		return err
	}

	if *statusCode != 200 {
		os.Remove(filePath)
		return fmt.Errorf("The server returns the status code %d", *statusCode)
	}

	return nil
}

func UpdateModuleRessources(config *FringeConfig) error {
	resourceDirectory, err := getResourceDirectory(config.HomeDirectory)
	if err != nil {
		return err
	}

	filenames := getRessourceFilenames()
	errorFiles := []string{}

	for _, filename := range filenames {
		err := updateModuleRessource(filename, resourceDirectory, config)
		if err != nil {
			errorFiles = append(errorFiles, filename)
		}
	}

	if len(errorFiles) > 0 {
		err = fmt.Errorf("Could not download the following files: \"%s\"", strings.Join(errorFiles, ", "))
		return err
	}

	return nil
}
