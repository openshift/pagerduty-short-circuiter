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

package pdcli

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
)

// ConnectionBuilder contains the information and logic needed to build a connection to PDCLI.
type ConnectionBuilder struct {
	cfg    *config.Config
	client *pagerduty.Client
}

// NewConnection creates a builder that can then be used to configure and build an PDCLI connection.
func NewConnection() *ConnectionBuilder {
	return &ConnectionBuilder{}
}

// Build uses the information stored in the builder to create a new PDCLI connection.
// Returns the PagerDuty API client object.
func (cb *ConnectionBuilder) Build() (client *pagerduty.Client, err error) {
	if cb.cfg == nil || cb.client == nil {
		// Load the configuration file
		cb.cfg, err = config.Fetch()

		if err != nil {
			err = fmt.Errorf("configuration file not found, run the 'pdcli login' command")
			return nil, err
		}

		if cb.cfg == nil {
			err = fmt.Errorf("not logged in, run the 'pdcli login' command")
			return nil, err
		}

		// If the config file exists, the API Key is validated
		_, err = config.ValidateKey(cb.cfg.ApiKey)

		if err != nil {
			err = fmt.Errorf("invalid API key, run the 'pdcli login' command")
			return nil, err
		}

		// Create a new PagerDuty API client
		client = pagerduty.NewClient(cb.cfg.ApiKey)
	}

	return client, nil
}
