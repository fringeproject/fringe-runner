package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Check if a file exist on the specify path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
