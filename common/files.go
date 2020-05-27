package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	FRINGE_HOME_NAME           = ".fringe-runner"
	FRINGE_RESSOURCE_DIRECTORY = "ressources"
	FRINGE_RESSOURCE_FILENAMES = "takeover_providers.json, wappalyzer.json"
	FRINGE_RESSOURCE_URL       = "https://static.fringeproject.com/fringe-runner/ressources/"
)

// Return the folder path to store the wordlists
func getFringeHomePath() (string, error) {
	// Enumerate the following directories in this order:
	// - The path set in the .env file `LISTS`
	// - The user home directory (HOME/%USERPROFILE%)
	// - The binary location

	value, exist := os.LookupEnv("FRINGE_HOME")
	if exist {
		return value, nil
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		return path.Join(homeDir, FRINGE_HOME_NAME), nil
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		return path.Join(dir, FRINGE_HOME_NAME), nil
	}

	return "", fmt.Errorf("Couldn't find fringe-runner home directory. Please specify a home dirctory for your user or the `FRINGE_HOME` environment variable.")
}

func getRessourceDirectory() (string, error) {
	fringeHome, err := getFringeHomePath()
	if err != nil {
		return "", err
	}

	ressourceDirectory := path.Join(fringeHome, FRINGE_RESSOURCE_DIRECTORY)

	err = os.MkdirAll(ressourceDirectory, 0700)
	if err != nil {
		err = fmt.Errorf("Cannot create the directory for the ressources.")
		return "", err
	}

	return ressourceDirectory, nil
}

func GetRessourceFile(filename string) (string, error) {
	ressourceDirectory, err := getRessourceDirectory()

	if err != nil {
		return "", nil
	}

	return path.Join(ressourceDirectory, filename), nil
}

func getRessourceFilenames() []string {
	return strings.Split(FRINGE_RESSOURCE_FILENAMES, ", ")
}

func updateModuleRessource(filename string, ressourceFolder string) error {
	url := FRINGE_RESSOURCE_URL + filename
	filePath := path.Join(ressourceFolder, filename)

	opt := HTTPOptions{}
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

func UpdateModuleRessources() error {
	ressourceDirectory, err := getRessourceDirectory()
	if err != nil {
		return err
	}

	filenames := getRessourceFilenames()
	errorFiles := []string{}

	for _, filename := range filenames {
		err := updateModuleRessource(filename, ressourceDirectory)
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

func ReadJSONFile(path string, result interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(byteValue), result)
}
