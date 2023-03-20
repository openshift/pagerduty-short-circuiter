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

	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
)

// Configuration struct to store user configuration.
type Config struct {
	ApiKey   string `json:"api_key,omitempty"`
	TeamID   string `json:"team_id,omitempty"`
	Team     string `json:"team,omitempty"`
	Terminal string `json:"terminal,omitempty"`
}

// Find returns the pdcli configuration filepath.
// If the config filepath doesn't exist, the desired config filepath string is created and returned.
func Find() (string, error) {

	// Return the test configuration filepath
	if kiteConfig := os.Getenv("KITE_CONFIG"); kiteConfig != "" {
		return kiteConfig, nil
	}

	// Locate the standard configuration directory
	configDir, err := os.UserConfigDir()

	if err != nil {
		return configDir, fmt.Errorf("cannot locate the user configuration directory")
	}

	configPath := filepath.Join(configDir, constants.ConfigFilepath)

	return configPath, nil
}

// Save saves the given configuration data to the config file.
// It creates a new directory to store the config file.
func Save(cfg *Config) error {

	file, err := Find()

	if err != nil {
		return err
	}

	// Check if the API key is valid
	cfg.ApiKey, err = validateKey(cfg.ApiKey)

	if err != nil {
		return err
	}

	if cfg.TeamID != "" {

		// Check if the team ID is valid
		cfg.TeamID, err = validateTeamID(cfg.TeamID)

		if err != nil {
			return err
		}
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

// Load loads the configuration file and parses it.
func Load() (config *Config, err error) {
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

	_, err = validateKey(config.ApiKey)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// validateKey sanitizes and validates the API key string.
func validateKey(apiKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)

	//Compare string with regex
	match, _ := regexp.MatchString(constants.APIKeyRegex, apiKey)

	if !match {
		return "", fmt.Errorf("invalid API key")
	}

	return apiKey, nil
}

// validateTeamID sanitizes and validates the team ID string.
func validateTeamID(teamID string) (string, error) {
	teamID = strings.TrimSpace(teamID)

	match, _ := regexp.MatchString(constants.TeamIdRegex, teamID)

	if !match {
		return "", fmt.Errorf("invalid Team ID")
	}

	return teamID, nil
}
