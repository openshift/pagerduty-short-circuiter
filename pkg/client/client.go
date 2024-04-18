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

package client

import (
	"fmt"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
)

// PagerDutyClient is an interface for the actual PD API
type PagerDutyClient interface {
	ListIncidents(pdApi.ListIncidentsOptions) (*pdApi.ListIncidentsResponse, error)
	ListIncidentAlerts(incidentId string) (*pdApi.ListAlertsResponse, error)
	GetCurrentUser(pdApi.GetCurrentUserOptions) (*pdApi.User, error)
	GetIncidentAlert(incidentID, alertID string) (*pdApi.IncidentAlertResponse, error)
	GetService(serviceID string, opts *pdApi.GetServiceOptions) (*pdApi.Service, error)
	ListOnCalls(opts pdApi.ListOnCallOptions) (*pdApi.ListOnCallsResponse, error)
	ManageIncidents(from string, incidents []pdApi.ManageIncidentsOptions) (*pdApi.ListIncidentsResponse, error)
}

type PDClient struct {
	cfg      *config.Config
	PdClient PagerDutyClient
}

// NewClient creates an instance of PDClient that is then used to connect to the actual pagerduty client.
func NewClient() *PDClient {
	return &PDClient{}
}

// Connect uses the information stored in new client to create a new PagerDuty connection.
// It returns the PDClient object with pagerduty API connection initialized.
func (pd *PDClient) Connect() (client *PDClient, err error) {

	if pd.cfg == nil {

		// Load the configuration file
		pd.cfg, err = config.Load()
		if err != nil {
			err = fmt.Errorf("configuration file not found, run the 'kite login' command")
			return nil, err
		}

		if pd.cfg == nil {
			err = fmt.Errorf("not logged in, run the 'kite login' command")
			return nil, err
		}

		if err != nil {
			err = fmt.Errorf("invalid API key, run the 'kite login' command")
			return nil, err
		}

		// Create a new PagerDuty API client
		pd.PdClient = pdApi.NewClient(pd.cfg.ApiKey)
	}

	return pd, nil
}

func (c *PDClient) ListIncidents(opts pdApi.ListIncidentsOptions) (*pdApi.ListIncidentsResponse, error) {
	return c.PdClient.ListIncidents(opts)
}

func (c *PDClient) ListIncidentAlerts(incidentID string) (*pdApi.ListAlertsResponse, error) {
	return c.PdClient.ListIncidentAlerts(incidentID)
}

func (c *PDClient) GetCurrentUser(opts pdApi.GetCurrentUserOptions) (*pdApi.User, error) {
	return c.PdClient.GetCurrentUser(opts)
}

func (c *PDClient) GetIncidentAlert(incidentID, alertID string) (*pdApi.IncidentAlertResponse, error) {
	return c.PdClient.GetIncidentAlert(incidentID, alertID)
}

func (c *PDClient) GetService(serviceID string, opts *pdApi.GetServiceOptions) (*pdApi.Service, error) {
	return c.PdClient.GetService(serviceID, opts)
}

func (c *PDClient) ListOnCalls(opts pdApi.ListOnCallOptions) (*pdApi.ListOnCallsResponse, error) {
	return c.PdClient.ListOnCalls(opts)
}

func (c *PDClient) ManageIncidents(from string, incidents []pdApi.ManageIncidentsOptions) (*pdApi.ListIncidentsResponse, error) {
	return c.PdClient.ManageIncidents(from, incidents)
}
