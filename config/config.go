// Package config provides functions related to stayback config file
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

// Config describes the data backup up to the s3 bucket
// directory entries that have leading slashes are assumed to be absolute : /opt/mystuff
// directories are assumed to be relative to $HOME
type Config struct {
	S3Bucket        string   `json:"s3Bucket"`
	EncryptedDirs   []string `json:"encryptedDirs"`
	UnEncryptedDirs []string `json:"unEncryptedDirs"`
	Timestamp       string   `json:"timestamp"`
}

// configFile returns a Config object from the json file path ptovided
// if the configPath is "", it uses $HOME/.stayback
func getConfig(configPath string) (config Config, err error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config, err
		}
		configPath = path.Join(homeDir, ".stayback")
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	return config, err
}

func configReport(config Config) (report string, err error) {
	return report, err
}

// validateDir returns the absolutel path of the directory if it exists
// given a relative dir, assume it's relative to the home directory
func validateDir(iDir string) (string, error) {
	var err error
	// if iDir is a relative path, assume it's relative to $HOME and expand it
	if !path.IsAbs(iDir) {

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		iDir = path.Join(homeDir, iDir)
	}
	//now iDir is an absolute path

	_, err = os.Stat(iDir)
	if err != nil {
		return "", err
	}
	return iDir, err
}

// sortUnique converts each directory entry to a validated absolute
// then it sorts the list and removes the duplicates
func sortUnique(iList []string) (oList []string, err error) {
	var absList []string

	for _, v := range iList {
		abs, err := validateDir(v)
		if err != nil {
			return []string{}, err
		}
		absList = append(absList, abs)
	}

	sort.Strings(absList)

	for i, v := range absList {
		if i == 0 {
			oList = append(oList, v)
		} else {
			if v != absList[i-1] {
				oList = append(oList, v)
			}
		}
	}

	return oList, err
}

// validateConfig given an existing config, reutn a valid config
// remove duplicated from the slices
// sort the slices
// validate the config
func validateConfig(rawConfig *Config) (validConfig Config, err error) {
	// render out and verify the directories before checking for duplicates so we're comparing absolutes
	// sort and remove
	return validConfig, err
}
