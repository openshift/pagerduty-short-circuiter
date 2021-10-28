package pdcli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

var Terminal string

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

// InitTerminalEmulator tries to set a terminal emulator by trying some known terminal emulators.
func InitTerminalEmulator() {
	emulators := []string{
		"x-terminal-emulator",
		"mate-terminal",
		"gnome-terminal",
		"terminator",
		"xfce4-terminal",
		"urxvt",
		"rxvt",
		"termit",
		"Eterm",
		"aterm",
		"uxterm",
		"xterm",
		"roxterm",
		"termite",
		"kitty",
		"hyper",
	}

	for _, t := range emulators {
		cmd := exec.Command("command", "-v", t)

		output, _ := cmd.CombinedOutput()

		cmd.ProcessState.Exited()

		term := string(output)

		term = strings.TrimSpace(term)

		if term != "" {
			Terminal = term
		}
	}
}

// ClusterLoginEmulator spawns an instance of ocm-container in a new terminal.
func ClusterLoginEmulator(clusterID string) error {

	var cmd *exec.Cmd

	// Check if ocm-container is installed locally
	ocmContainer, err := exec.LookPath("ocm-container")

	if err != nil {
		return errors.New("ocm-container is not found.\nPlease install it via: " + constants.OcmContainerURL)
	}

	// OCM container command to be executed for cluster login
	ocmCommand := ocmContainer + " " + clusterID

	cmd = exec.Command(Terminal, "-e", ocmCommand)

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// ClusterLoginShell spawns an instance of ocm-container in the same shell.
func ClusterLoginShell(clusterID string) *exec.Cmd {

	// Check if ocm-container is installed locally
	ocmContainer, err := exec.LookPath("ocm-container")

	if err != nil {
		fmt.Println("ocm-container is not found.\nPlease install it via:", constants.OcmContainerURL)
	}

	cmd := exec.Command(ocmContainer, clusterID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
