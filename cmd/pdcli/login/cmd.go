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
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/spf13/cobra"
)

var loginArgs struct {
	apiKey string
}

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the PagerDuty CLI",
	Long: `The pdcli login command logs a user into PagerDuty CLI given a valid API key is provided. 
	You will have to login only once, all the pdcli commands are then available even if the terminal restarts.`,
	Args: cobra.NoArgs,
	RunE: loginHandler,
}

func init() {

	Cmd.Flags().StringVar(
		&loginArgs.apiKey,
		"api-key",
		"",
		"Access API key/token generated from "+constants.APIKeyURL+"\nUse this option to overwrite the existing API key.",
	)
}

// loginHandler handles the login flow into pdcli.
func loginHandler(cmd *cobra.Command, args []string) error {

	// Currently logged in user
	var user string
	var pdClient client.PagerDutyClient

	// load the configuration info
	cfg, err := config.Load()

	// If no config file can be located
	// Or if the config file has errors
	// Or if this is the first time a user is trying to login
	// A new configuration struct is initialized on login
	if err != nil {
		cfg = new(config.Config)
	}

	// If the key arg is not-empty
	if loginArgs.apiKey != "" {
		cfg.ApiKey = loginArgs.apiKey

		// Save the key in the config file
		err = config.Save(cfg)

		if err != nil {
			return err
		}

		// PagerDuty API client
		pdClient, err = client.NewClient().Connect()

		if err != nil {
			return err
		}

		user, err = Login(cfg.ApiKey, pdClient)

		if err != nil {
			return err
		}

		successMessage(user)

		return nil
	}

	// API key is not found in the config file
	if len(cfg.ApiKey) == 0 {

		// Create a new API key and store it in the config file
		err = generateNewKey(cfg)

		if err != nil {
			return err
		}

		// PagerDuty API client
		pdClient, err = client.NewClient().Connect()

		if err != nil {
			return err
		}

		// Login using the newly generated API Key
		user, err = Login(cfg.ApiKey, pdClient)

		if err != nil {
			return err
		}

		successMessage(user)

	} else {

		// PagerDuty API client
		pdClient, err = client.NewClient().Connect()

		if err != nil {
			return err
		}

		// Login using the existing API key in the configuration file
		user, err = Login(cfg.ApiKey, pdClient)

		if err != nil {
			return err
		}

		successMessage(user)
	}

	return nil
}

// generateNewKey prompts the user to create a new API key and saves it to the config file.
func generateNewKey(cfg *config.Config) (err error) {
	//prompts the user to generate an API Key
	fmt.Println("In order to login it is mandatory to provide an API key.\nThe recommended way is to generate an API key via: " + constants.APIKeyURL)

	//Takes standard input from the user and stores it in a variable
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("API Key: ")
	cfg.ApiKey, err = reader.ReadString('\n')

	if err != nil {
		return err
	}

	err = config.Save(cfg)

	if err != nil {
		return err
	}

	return nil
}

// Login handles PagerDuty REST API authentication via an user API token.
// Requests that cannot be authenticated will return a `401 Unauthorized` error response.
// It returns the username of the currently logged in user.
func Login(apiKey string, client client.PagerDutyClient) (string, error) {

	user, err := client.GetCurrentUser(pagerduty.GetCurrentUserOptions{})

	if err != nil {
		var apiError pagerduty.APIError

		//`401 Unauthorized` error response
		if errors.As(err, &apiError) {
			err = fmt.Errorf("login failed\n%v Unauthorized", apiError.StatusCode)
			return "", err
		}

		return "", err
	}

	return user.Name, nil
}

// successMessage prints the currently logged in username to the console.
// if pagerduty login is successful.
func successMessage(user string) {
	fmt.Printf("Successfully logged in as user: %s\n", user)
}
