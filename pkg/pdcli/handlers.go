package pdcli

import (
	"errors"
	"fmt"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

type Alert struct {
	IncidentID  string
	AlertID     string
	ClusterID   string
	ClusterName string
	Name        string
	Console     string
	Hostname    string
	IP          string
	Labels      string
	LastCheckIn string
	Severity    string
	Status      string
	Sop         string
	Token       string
	Tags        string
	WebURL      string
}

// GetIncidents returns a slice of pagerduty incidents.
func GetIncidents(c client.PagerDutyClient, opts *pdApi.ListIncidentsOptions) ([]pdApi.Incident, error) {

	var aerr pdApi.APIError

	// Get incidents via pagerduty API
	incidents, err := c.ListIncidents(*opts)

	if err != nil {
		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				fmt.Println("rate limited")
				return nil, err
			}

			fmt.Println("status code:", aerr.StatusCode)

			return nil, err
		}
	}

	return incidents.Incidents, nil
}

// GetIncidentAlerts returns all the alerts belong to a particular incident.
func GetIncidentAlerts(c client.PagerDutyClient, incidentID string) ([]Alert, error) {

	var alerts []Alert

	// Fetch alerts related to an incident via pagerduty API
	incidentAlerts, err := c.ListIncidentAlerts(incidentID)

	if err != nil {
		var aerr pdApi.APIError

		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				fmt.Println("rate limited")
				return nil, err
			}

			fmt.Println("status code:", aerr.StatusCode)

			return nil, err
		}
	}

	for _, alert := range incidentAlerts.Alerts {

		// Check if the alert is not resolved
		if alert.Status != constants.StatusResolved {
			tempAlertObj := Alert{}
			err = tempAlertObj.ParseAlertData(c, &alert)

			if err != nil {
				return nil, err
			}

			alerts = append(alerts, tempAlertObj)
		}

	}

	return alerts, nil
}

// GetClusterName interacts with the PD service endpoint and returns the cluster name as a string.
func GetClusterName(servideID string, c client.PagerDutyClient) (string, error) {

	service, err := c.GetService(servideID, &pdApi.GetServiceOptions{})

	if err != nil {
		return "", err
	}

	clusterName := strings.Split(service.Description, " ")[0]

	return clusterName, nil
}

// GetCurrentUserID returns the ID of the currently logged in user.
func GetCurrentUserID(c client.PagerDutyClient) (string, error) {
	var aerr pdApi.APIError

	// Get current user details
	user, err := c.GetCurrentUser(pdApi.GetCurrentUserOptions{})

	if err != nil {
		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				fmt.Println("rate limited")
				return "", err
			}

			fmt.Println("status code:", aerr.StatusCode)

			return "", err
		}
	}

	return user.ID, nil
}

// AcknowledgeIncidents acknowledges incidents for the given incident IDs
// and retuns the acknowledged incidents
func AcknowledgeIncidents(c client.PagerDutyClient, incidentIDs []string) ([]pdApi.Incident, error) {
	var incidents []pdApi.ManageIncidentsOptions
	var opts pdApi.ManageIncidentsOptions

	var response *pdApi.ListIncidentsResponse

	for _, id := range incidentIDs {
		opts.ID = id
		opts.Type = "incident"
		opts.Status = constants.StatusAcknowledged

		incidents = append(incidents, opts)
	}

	user, err := c.GetCurrentUser(pdApi.GetCurrentUserOptions{})

	if err != nil {
		return nil, err
	}

	response, err = c.ManageIncidents(user.Email, incidents)

	if err != nil {
		return nil, err
	}

	return response.Incidents, nil
}

// ParseAlertData parses a pagerduty alert data into the Alert struct.
func (a *Alert) ParseAlertData(c client.PagerDutyClient, alert *pdApi.IncidentAlert) (err error) {

	a.IncidentID = alert.Incident.ID
	a.AlertID = alert.ID
	a.Name = alert.Summary
	a.Severity = alert.Severity
	a.Status = alert.Status
	a.WebURL = alert.HTMLURL

	// Check if the alert is of type 'Missing cluster'
	isCHGM := alert.Body["details"].(map[string]interface{})["notes"]

	// Check if the alert is of type 'Certificate is expiring'
	isCertExpiring := alert.Body["details"].(map[string]interface{})["hostname"]

	if isCHGM != nil {
		notes := strings.Split(fmt.Sprint(alert.Body["details"].(map[string]interface{})["notes"]), "\n")

		a.ClusterID = strings.Replace(notes[0], "cluster_id: ", "", 1)
		a.ClusterName = strings.Split(fmt.Sprint(alert.Body["details"].(map[string]interface{})["name"]), ".")[0]

		lastCheckIn := fmt.Sprint(alert.Body["details"].(map[string]interface{})["last healthy check-in"])
		a.LastCheckIn, err = utils.FormatTimestamp(lastCheckIn)

		if err != nil {
			return err
		}

		a.Token = fmt.Sprint(alert.Body["details"].(map[string]interface{})["token"])
		a.Tags = fmt.Sprint(alert.Body["details"].(map[string]interface{})["tags"])
		a.Sop = strings.Replace(notes[1], "runbook: ", "", 1)

	} else if isCertExpiring != nil {
		a.Hostname = fmt.Sprint(alert.Body["details"].(map[string]interface{})["hostname"])
		a.IP = fmt.Sprint(alert.Body["details"].(map[string]interface{})["ip"])
		a.Sop = fmt.Sprint(alert.Body["details"].(map[string]interface{})["url"])
		a.Name = strings.Split(alert.Summary, " on ")[0]
		a.ClusterName = "N/A"

	} else {
		a.ClusterID = fmt.Sprint(alert.Body["details"].(map[string]interface{})["cluster_id"])
		a.ClusterName, err = GetClusterName(alert.Service.ID, c)

		// If the service mapped to the current incident is not available (404)
		if err != nil {
			a.ClusterName = "N/A"
		}

		a.Console = fmt.Sprint(alert.Body["details"].(map[string]interface{})["console"])
		a.Labels = fmt.Sprint(alert.Body["details"].(map[string]interface{})["firing"])
		a.Sop = fmt.Sprint(alert.Body["details"].(map[string]interface{})["link"])
	}

	return nil
}
