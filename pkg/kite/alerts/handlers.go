package kite

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

type User struct {
	UserID     string
	Name       string
	Role       string
	Team       string
	Email      string
	AssignedTo string
}

// GetIncidents returns a slice of pagerduty incidents.
func GetIncidents(c client.PagerDutyClient, opts *pdApi.ListIncidentsOptions) ([]pdApi.Incident, error) {
	var aerr pdApi.APIError

	// Get incidents via pagerduty API
	incidents, err := c.ListIncidents(*opts)

	if err != nil {
		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				return nil, fmt.Errorf("API rate limited")
			}
			return nil, fmt.Errorf("status code: %d, error: %s", aerr.StatusCode, err)
		}
	}

	return incidents.Incidents, nil
}

// GetIncidentAlerts returns all the alerts belonging to a particular incident.
func GetIncidentAlerts(c client.PagerDutyClient, incident pdApi.Incident) ([]Alert, error) {
	var alerts []Alert

	// Fetch alerts related to an incident via pagerduty API
	incidentAlerts, err := c.ListIncidentAlerts(incident.Id)

	if err != nil {
		var aerr pdApi.APIError

		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				return nil, fmt.Errorf("API rate limited")
			}

			return nil, fmt.Errorf("status code: %d, error: %s", aerr.StatusCode, err)
		}
	}

	for _, alert := range incidentAlerts.Alerts {
		if alert.Status != constants.StatusResolved {
			tempAlertObj := Alert{}

			err = tempAlertObj.ParseAlertData(c, &alert)

			if err != nil {
				return nil, err
			}

			// Fetch incident Urgency
			tempAlertObj.Severity = incident.Urgency

			if tempAlertObj.Severity == "" {
				tempAlertObj.Severity = alert.Severity
			}

			alerts = append(alerts, tempAlertObj)
		}
	}

	return alerts, nil
}

// GetClusterName interacts with the PD service endpoint and returns the cluster name string.
func GetClusterName(servideID string, c client.PagerDutyClient) (string, error) {
	service, err := c.GetService(servideID, &pdApi.GetServiceOptions{})

	if err != nil {
		return "", err
	}

	ClusterName := strings.Split(service.Description, " ")[0]

	return ClusterName, nil
}

// AcknowledgeIncidents acknowledges incidents for the given incident IDs
// and retuns the acknowledged incidents.
func AcknowledgeIncidents(c client.PagerDutyClient, incidentIDs []string, pdUser User) ([]pdApi.Incident, error) {
	var incidents []pdApi.ManageIncidentsOptions
	var opts pdApi.ManageIncidentsOptions

	var response *pdApi.ListIncidentsResponse

	for _, id := range incidentIDs {
		opts.ID = id
		opts.Type = "incident"
		opts.Status = constants.StatusAcknowledged

		incidents = append(incidents, opts)
	}

	response, err := c.ManageIncidents(pdUser.Email, incidents)

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

	// If there's no cluster ID related to the given alert
	if a.ClusterID == "" {
		a.ClusterID = "N/A"
	}

	return nil
}

// ParseAlertMetaData parses the given alert metadata into a string and returns it.
func ParseAlertMetaData(alert Alert) string {
	var alertData string

	if alert.ClusterID != "" {
		data := fmt.Sprintf("* Cluster ID: %s\n", alert.ClusterID)
		alertData = alertData + data
	}

	if alert.ClusterName != "" {
		data := fmt.Sprintf("* Cluster Name: %s\n", alert.ClusterName)
		alertData = alertData + data
	}

	if alert.Console != "" {
		data := fmt.Sprintf("* Console: %s\n", alert.Console)
		alertData = alertData + data
	}

	if alert.Hostname != "" {
		data := fmt.Sprintf("* Hostname: %s\n", alert.Hostname)
		alertData = alertData + data
	}

	if alert.IP != "" {
		data := fmt.Sprintf("* IP: %s\n", alert.IP)
		alertData = alertData + data
	}

	if alert.LastCheckIn != "" {
		data := fmt.Sprintf("* Last Healthy Check-in: %s\n", alert.LastCheckIn)
		alertData = alertData + data
	}

	if alert.Tags != "" {
		data := fmt.Sprintf("* Tags: %s\n", alert.Tags)
		alertData = alertData + data
	}

	if alert.Token != "" {
		data := fmt.Sprintf("* Token: %s\n", alert.Token)
		alertData = alertData + data
	}

	if alert.Labels != "" {
		data := fmt.Sprintf("* %s", alert.Labels)
		alertData = alertData + data
	}

	if alert.Sop != "" {
		data := fmt.Sprintf("* SOP: %s\n", alert.Sop)
		alertData = alertData + data
	}

	if alert.WebURL != "" {
		data := fmt.Sprintf("* Web URL: %s\n", alert.WebURL)
		alertData = alertData + data
	}

	return alertData
}

// getTableData parses and returns tabular data for the given alerts, i.e table headers and rows.
func GetAlertsTableData(alerts []Alert) ([]string, [][]string) {
	headers := []string{"INCIDENT ID", "ALERT ID", "ALERT", "CLUSTER", "CLUSTER ID", "STATUS", "SEVERITY"}

	var tableData [][]string

	for _, alert := range alerts {
		tableData = append(tableData, []string{
			alert.IncidentID,
			alert.AlertID,
			alert.Name,
			alert.ClusterName,
			alert.ClusterID,
			alert.Status,
			alert.Severity,
		})
	}

	return headers, tableData
}

// getTableData parses and returns tabular data for the given incidents, i.e table rows.
func GetIncidentsTableData(ackIncidents []pdApi.Incident, triggeredIncidents []pdApi.Incident) (ackTableData, triggeredTableData [][]string) {
	for _, incident := range ackIncidents {
		ackTableData = append(ackTableData, []string{
			incident.Id,
			incident.Title,
			incident.Urgency,
			incident.Status,
			incident.Service.Summary,
			incident.Acknowledgements[0].Acknowledger.Summary,
		})
	}

	for _, incident := range triggeredIncidents {
		triggeredTableData = append(triggeredTableData, []string{
			incident.Id,
			incident.Title,
			incident.Urgency,
			incident.Status,
			incident.Service.Summary,
			incident.Assignments[0].Assignee.Summary,
		})
	}

	return ackTableData, triggeredTableData
}

// FilterAlertsByStatus filters the given alerts based on its status and returns low, high alerts.
func FilterAlertsByStatus(alerts []Alert) (low, high []Alert) {
	for _, alert := range alerts {
		if alert.Severity == constants.StatusHigh {
			high = append(high, alert)
		}

		if alert.Severity == constants.StatusLow {
			low = append(low, alert)
		}
	}

	return low, high
}

// FilterIncidentsByStatus filters the given incidents based on its urgency
// and returns acknowledged and trigerred (un-acknowledged) incidents tabular data.
func FilterIncidentsByStatus(incidents []pdApi.Incident) (ackIncidents, trigerredIncidents []pdApi.Incident) {
	for _, incident := range incidents {
		if incident.Status == constants.StatusAcknowledged {
			ackIncidents = append(ackIncidents, incident)
		}

		if incident.Status == constants.StatusTriggered {
			trigerredIncidents = append(trigerredIncidents, incident)
		}
	}

	return ackIncidents, trigerredIncidents
}
