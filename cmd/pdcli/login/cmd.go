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
package login

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/spf13/cobra"
)

const APIKeyURL = "https://support.pagerduty.com/docs/generating-api-keys#section-generating-a-general-access-rest-api-key"

var userKey string

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the PagerDuty CLI",
	Long: `Running the pdcli login command will send a request to PagerDuty REST API provided a valid API key.
The PagerDuty REST API supports authenticating via the user API token.`,
	Args: cobra.NoArgs,
	RunE: loginHandler,
}

func init() {
	//flags
	Cmd.Flags().StringVar(&userKey, "key", "", "Access API key/token generated from "+APIKeyURL+"\nUse this option to overwrite the existing API key.")
}

// loginHandler handles the login flow into pdcli
func loginHandler(cmd *cobra.Command, args []string) error {

	//load configuration info
	cfg, err := config.Fetch()

	if err != nil {
		return fmt.Errorf("cannot load config file: %v", err)
	}

	// check if config file is empty
	if cfg == nil {
		cfg = new(config.Config)
	}

	// if the key flag is given
	if userKey != "" {
		cfg.ApiKey, err = validateKey(userKey)

		if err != nil {
			return err
		}

		config.Save(cfg)

		return nil
	}

	// API key is not found in the config file
	if len(cfg.ApiKey) == 0 {
		err = generateNewKey(cfg)

		if err != nil {
			return err
		}

	} else {
		fmt.Println("Login Successful")
	}

	return nil
}

// generateNewKey prompts the user to create a new API key and saves it to the config file
func generateNewKey(cfg *config.Config) error {
	//prompts the user to generate an API Key
	fmt.Println("In order to login it is mandatory to provide an API key.\nThe recommended way is to generate an API key via: " + APIKeyURL)

	//Takes standard input from the user and stores it in a variable
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("API Key: ")
	apiKey, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	cfg.ApiKey, err = validateKey(apiKey)

	if err != nil {
		return err
	}

	err = config.Save(cfg)

	if err != nil {
		return err
	}

	return nil
}

// validateKey sanitizes and validates the API key string
func validateKey(apiKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)

	//compare string with regex
	match, _ := regexp.MatchString("^[a-z|A-Z0-9_-]{20}$", apiKey)

	if !match {
		return "", fmt.Errorf("invalid API key")
	}

	return apiKey, nil
}
