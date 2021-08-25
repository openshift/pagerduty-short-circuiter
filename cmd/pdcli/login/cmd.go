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

	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/spf13/cobra"
)

const APIKeyURL = "https://support.pagerduty.com/docs/generating-api-keys#section-generating-a-general-access-rest-api-key"

var userKey string

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "PagerDuty CLI login",
	Long:  "Logs a user into the pdcli provided the user has a valid API key",
	Args:  cobra.NoArgs,
	RunE:  loginHandler,
}

func init() {
	//flags
	Cmd.Flags().StringVar(&userKey, "token", "", "Access API key/token generated from "+APIKeyURL)
}

func loginHandler(cmd *cobra.Command, args []string) error {

	//load configuration info
	cfg, err := config.Fetch()

	if err != nil {
		return fmt.Errorf("can't load config file: %v", err)
	}

	if cfg == nil {
		cfg = new(config.Config)
	}

	if cfg.ApiKey == "" {
		generateNewKey(cfg)
	} else {
		fmt.Println("Login Successful.")
	}

	return nil

}

func generateNewKey(cfg *config.Config) (string, error) {
	//prompts the user to generate an API Key
	fmt.Println("In order to login it is mandatory to provide an API key.\nThe recommended way is to generate an API key via: " + APIKeyURL)

	//Takes standard input from the user and stores it in a variable
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("API Key: ")
	apiKey, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	cfg.ApiKey = apiKey

	err = config.Save(cfg)

	if err != nil {
		return "", err
	}

	return apiKey, nil
}
