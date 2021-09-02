/*
Copyright © 2021 Red Hat, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Configuration struct to store user configuration.
type Config struct {
	ApiKey string `json:"api_key,omitempty"`
}

// Find returns the pdcli configuration filepath string.
// If the config filepath doesn't exist, the desired config filepath string is returned.
func Find() (string, error) {
	//locates the user home directory
	homedir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("cannot locate user home directory")
	}

	configPath := filepath.Join(homedir, ".config/pagerduty-cli/config.json")

	//Check if the path exists
	_, err = os.Stat(configPath)

	//If the config path doesn't exist
	if os.IsNotExist(err) {
		configDir, err := os.UserConfigDir()

		if err != nil {
			return configDir, err
		}

		configPath = filepath.Join(configDir, "pagerduty-cli/config.json")

		if err != nil {
			return "", err
		}
	}

	return configPath, nil
}

// Save saves the given configuration data to the config file.
// It creates a new directory to store the config file.
func Save(cfg *Config) error {
	file, err := Find()

	if err != nil {
		return err
	}

	// Create a new directory to store config file
	dir := filepath.Dir(file)
	err = os.MkdirAll(dir, os.FileMode(0755))

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return fmt.Errorf("cannot marshal configuration file: %v", err)
	}

	err = ioutil.WriteFile(file, data, 0600)

	if err != nil {
		return fmt.Errorf("cannot save configuration file '%s': %v", file, err)
	}

	return nil
}

// Fetch loads the configuration file and parses it.
func Fetch() (config *Config, err error) {
	//Locate the config filepath
	configFile, err := Find()

	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configFile)

	if os.IsNotExist(err) {
		return nil, err
	}

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		err = fmt.Errorf("cannot read config file")
		return nil, err
	}

	if len(configData) == 0 {
		err = fmt.Errorf("configuration file is empty")
		return nil, err
	}

	config = &Config{}
	err = json.Unmarshal(configData, config)

	if err != nil {
		err = fmt.Errorf("error parsing config file")
		return nil, err
	}

	return
}

// ValidateKey sanitizes and validates the API key string.
func ValidateKey(apiKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)

	//compare string with regex
	match, _ := regexp.MatchString("^[a-z|A-Z0-9+_-]{20}$", apiKey)

	if !match {
		return "", fmt.Errorf("invalid API key")
	}

	return apiKey, nil
}
