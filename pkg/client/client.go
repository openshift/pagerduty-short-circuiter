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
	"net/http"

	pdApi "github.com/PagerDuty/go-pagerduty"
)

// PagerDutyClient is an interface for the actual PD API
type PagerDutyClient interface {
	ListIncidents(pdApi.ListIncidentsOptions) (*pdApi.ListIncidentsResponse, error)
	ListIncidentAlerts(incidentId string) (*pdApi.ListAlertsResponse, error)
	GetCurrentUser(pdApi.GetCurrentUserOptions) (*pdApi.User, error)
	GetIncidentAlert(incidentID, alertID string) (*pdApi.IncidentAlertResponse, *http.Response, error)
	GetService(serviceID string, opts *pdApi.GetServiceOptions) (*pdApi.Service, error)
}

// PDClient wraps pdApi.Client
type PDClient struct {
	APIKey   string
	PdClient PagerDutyClient
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

func (c *PDClient) GetIncidentAlert(incidentID, alertID string) (*pdApi.IncidentAlertResponse, *http.Response, error) {
	return c.PdClient.GetIncidentAlert(incidentID, alertID)
}

func (c *PDClient) GetService(serviceID string, opts *pdApi.GetServiceOptions) (*pdApi.Service, error) {
	return c.PdClient.GetService(serviceID, opts)
}
