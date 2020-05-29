package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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
