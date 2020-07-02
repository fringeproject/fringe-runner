package common

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type StringArray []string

// From: https://github.com/go-yaml/yaml/issues/100
func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

// Parse a YAML file to a structure
func ParseYAMLFile(filePath string, out interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Warning(err)
		err = fmt.Errorf("Cannot read the YAML file: %s.", filePath)
		return err
	}

	err = yaml.Unmarshal([]byte(file), out)
	if err != nil {
		err = fmt.Errorf("The YAML file does not match the interface structure.")
		return err
	}

	return nil
}
