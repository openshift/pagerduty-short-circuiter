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
	"errors"
	"fmt"
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

const APIKeyURL = "https://support.pagerduty.com/docs/generating-api-keys#section-generating-a-general-access-rest-api-key"

var loginArgs struct {
	apiKey string
}

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
	Cmd.Flags().StringVar(&loginArgs.apiKey, "key", "", "Access API key/token generated from "+APIKeyURL+"\nUse this option to overwrite the existing API key.")
}

// loginHandler handles the login flow into pdcli.
func loginHandler(cmd *cobra.Command, args []string) error {

	// load configuration info
	cfg, err := config.Fetch()

	// if the config file is located
	// check if config file is empty, initialize a new config struct
	if cfg == nil {
		cfg = new(config.Config)
	}

	// if no config file can be located
	// initialize a new config struct to parse and save it to a new config file
	if err != nil {
		cfg = new(config.Config)
	}

	// if the key arg is not-empty
	if loginArgs.apiKey != "" {
		cfg.ApiKey, err = config.ValidateKey(loginArgs.apiKey)

		if err != nil {
			return err
		}

		// Save the key in the config file
		err = config.Save(cfg)

		if err != nil {
			return err
		}

		err = login(cfg.ApiKey)

		if err != nil {
			return err
		}

		return nil
	}

	// API key is not found in the config file
	if len(cfg.ApiKey) == 0 {

		// Create a new API key and store it in the config file
		err = generateNewKey(cfg)

		if err != nil {
			return err
		}

		// Login using the newly generated API Key
		err = login(cfg.ApiKey)

		if err != nil {
			return err
		}

	} else {

		// Login using the existing API key in the configuration file
		err = login(cfg.ApiKey)

		if err != nil {
			return err
		}
	}

	return nil
}

// generateNewKey prompts the user to create a new API key and saves it to the config file.
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

	cfg.ApiKey, err = config.ValidateKey(apiKey)

	if err != nil {
		return err
	}

	err = config.Save(cfg)

	if err != nil {
		return err
	}

	return nil
}

// login handles PagerDuty REST API authentication via an user API token.
// Requests that cannot be authenticated will return a `401 Unauthorized` error response.
func login(apiKey string) error {

	// PagerDuty client object is created
	client, err := pdcli.NewConnection().Build()

	if err != nil {
		return err
	}

	user, err := client.GetCurrentUser(pagerduty.GetCurrentUserOptions{})

	if err != nil {
		var apiError pagerduty.APIError

		//`401 Unauthorized` error response
		if errors.As(err, &apiError) {
			err = fmt.Errorf("login failed\n%v Unauthorized", apiError.StatusCode)
			return err
		}

		return err
	} else {
		fmt.Printf("Successfully logged in as user: %s\n", user.Name)
	}

	return nil
}
