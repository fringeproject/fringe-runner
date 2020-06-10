package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	FRINGE_HOME_DIRECTORY = ".fringe-runner"
)

type FringeConfig struct {
	LogLevel            string            `yaml:"log_level"`
	Proxy               string            `yaml:"proxy"`
	VerifyCert          bool              `yaml:"verify_cert"`
	HomeDirectory       string            `yaml:"home_directory"`
	FringeCoordinator   string            `yaml:"fringe_coordinator"`
	FringePerimeter     string            `yaml:"fringe_perimeter"`
	FringeRunnerId      string            `yaml:"fringe_runner_id"`
	FringeRunnerToken   string            `yaml:"fringe_runner_token"`
	ModuleConfiguration map[string]string `yaml:"module_configuration"`
}

// Return the folder path to store the wordlists
func findFringeHomePath() (string, error) {
	// Enumerate the following directories in this order:
	// - The path set the `FRINGE_HOME` env variable
	// - The user home directory (HOME/%USERPROFILE%)
	// - The binary location

	value, exist := os.LookupEnv("FRINGE_HOME")
	if exist {
		return value, nil
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		return path.Join(homeDir, FRINGE_HOME_DIRECTORY), nil
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		return path.Join(dir, FRINGE_HOME_DIRECTORY), nil
	}

	return "", fmt.Errorf("Couldn't find fringe-runner home directory. Please specify a home dirctory for your user or the `FRINGE_HOME` environment variable.")
}

func ReadConfigFile(configPath string) (*FringeConfig, error) {
	if strings.HasPrefix(configPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(homeDir, configPath[2:])
		}
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Warning(err)
		err = fmt.Errorf("The configuration file does not exist at %s.", configPath)
		return nil, err
	}

	// Instantiate the structure with default values
	config := FringeConfig{
		LogLevel:            "info",
		Proxy:               "",
		VerifyCert:          true,
		HomeDirectory:       "", // set the empty string for futur check
		ModuleConfiguration: map[string]string{},
	}
	err = yaml.Unmarshal([]byte(configFile), &config)
	if err != nil {
		err = fmt.Errorf("The configuration file is not a valid YAML file.")
		return nil, err
	}

	// Check some config values
	if config.HomeDirectory == "" {
		defaultHome, err := findFringeHomePath()
		if err != nil {
			err = fmt.Errorf("Cannot find a default home directory for the runner, please specify configuration key \"home_directory\".")
			return nil, err
		}
		config.HomeDirectory = defaultHome
	}

	return &config, nil
}
