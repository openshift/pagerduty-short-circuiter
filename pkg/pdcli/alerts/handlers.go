package pdcli

import (
	"errors"
	"fmt"
	"sort"
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
	Notes       string
}

var (
	TrigerredAlerts []Alert
)

// GetIncidents returns a slice of pagerduty incidents.
func GetIncidents(c client.PagerDutyClient, opts *pdApi.ListIncidentsOptions) ([]pdApi.Incident, error) {
	var aerr pdApi.APIError
	var incidents []pdApi.Incident

	// Check if incidents are fetched for a Team
	isTeam := len(opts.TeamIDs) > 0

	// Get incidents via pagerduty API
	incidentsList, err := c.ListIncidents(*opts)

	if err != nil {
		if errors.As(err, &aerr) {
			if aerr.RateLimited() {
				return nil, fmt.Errorf("API rate limited")
			}
			return nil, fmt.Errorf("status code: %d, error: %s", aerr.StatusCode, err)
		}
	}

	for _, incident := range incidentsList.Incidents {
		// When incidents are fetched for a team, do not include the incidents assigned to SilentTest
		if isTeam && (incident.EscalationPolicy.ID == constants.SilentTestEscalationPolicyID ||
			incident.EscalationPolicy.ID == constants.CADSilentTestEscalationPolicyID ||
			incident.EscalationPolicy.ID == constants.CADSilentTestStageEscalationPolicyID ||
			incident.Assignments[0].Assignee.ID == constants.SilentTest ||
			incident.Assignments[0].Assignee.ID == constants.NobodySREP) {
			continue
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

// GetIncidentAlerts returns all the alerts belonging to a particular incident.
func GetIncidentAlerts(c client.PagerDutyClient, incident pdApi.Incident) ([]Alert, error) {
	var alerts []Alert

	// Fetch alerts related to an incident via pagerduty API
	incidentAlerts, err := c.ListIncidentAlerts(incident.APIObject.ID)

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
		status := alert.Status

		tempAlertObj := Alert{}

		// Fetch incident Urgency
		tempAlertObj.Severity = incident.Urgency

		if tempAlertObj.Severity == "" {
			tempAlertObj.Severity = alert.Severity
		}

		if status == constants.StatusTriggered {
			err = tempAlertObj.ParseAlertData(c, &alert)

			if err != nil {
				return nil, err
			}

			TrigerredAlerts = append(TrigerredAlerts, tempAlertObj)
		}

		alerts = append(alerts, tempAlertObj)
	}

	return alerts, nil
}

// GetClusterName interacts with the PD service endpoint and returns the cluster name string.
func GetClusterName(servideID string, c client.PagerDutyClient) (string, error) {
	service, err := c.GetService(servideID, &pdApi.GetServiceOptions{})

	if err != nil {
		return "", err
	}

	clusterName := strings.Split(service.Description, " ")[0]

	return clusterName, nil
}

// AcknowledgeIncidents acknowledges incidents for the given incident IDs
// and retuns the acknowledged incidents.
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
func GetTableData(alerts []Alert, cols string) ([]string, [][]string) {
	var headers []string
	var tableData [][]string

	// columns returned by the columns flag
	columns := strings.Split(cols, ",")

	columnsMap := make(map[string]bool)

	for _, c := range columns {
		columnsMap[c] = true
	}

	headersMap := make(map[int]string)

	for _, alert := range alerts {

		var values []string

		var i int

		if columnsMap["incident.id"] {
			i++
			headersMap[i] = "INCIDENT ID"
			values = append(values, alert.IncidentID)
		}

		if columnsMap["alert.id"] {
			i++
			headersMap[i] = "ALERT ID"
			values = append(values, alert.AlertID)
		}

		if columnsMap["alert"] {
			i++
			headersMap[i] = "ALERT"
			values = append(values, alert.Name)
		}

		if columnsMap["cluster.name"] {
			i++
			headersMap[i] = "CLUSTER NAME"
			values = append(values, alert.ClusterName)
		}

		if columnsMap["cluster.id"] {
			i++
			headersMap[i] = "CLUSTER ID"
			values = append(values, alert.ClusterID)
		}

		if columnsMap["status"] {
			i++
			headersMap[i] = "STATUS"
			values = append(values, alert.Status)
		}

		if columnsMap["severity"] {
			i++
			headersMap[i] = "SEVERITY"
			values = append(values, alert.Severity)
		}

		tableData = append(tableData, values)
	}

	keys := make([]int, 0)

	for k := range headersMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, v := range keys {
		headers = append(headers, headersMap[v])
	}

	return headers, tableData
}
