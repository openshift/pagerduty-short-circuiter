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

package alerts

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/manifoldco/promptui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/output"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

type Alert struct {
	IncidentID  string
	AlertID     string
	ClusterID   string
	ClusterName string
	Name        string
	Console     string
	Labels      string
	LastCheckIn string
	Severity    string
	Status      string
	Sop         string
	Token       string
	Tags        string
}

var options struct {
	high        bool
	low         bool
	assignment  string
	columns     string
	interactive bool
	incidentID  bool
}

var Cmd = &cobra.Command{
	Use:   "alerts",
	Short: "This command will list all the open high alerts assigned to self.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  alertsHandler,
}

func init() {

	// Urgency
	Cmd.Flags().BoolVar(&options.low, "low", false, "View all low alerts")
	Cmd.Flags().BoolVar(&options.high, "high", true, "View all high alerts")

	// Incident Assignment
	Cmd.Flags().StringVar(
		&options.assignment,
		"assigned-to",
		"self",
		"Filter alerts based on user or team",
	)

	// Columns displayed
	Cmd.Flags().StringVar(
		&options.columns,
		"columns",
		"incident.id,cluster.name,alert,cluster.id,status,severity",
		"Specify which columns to display separated by commas without any space in between",
	)

	// Interactive mode
	Cmd.Flags().BoolVarP(
		&options.interactive,
		"interactive",
		"i",
		false,
		"Use interactive mode",
	)

}

// alertsHandler is the main alerts command handler.
func alertsHandler(cmd *cobra.Command, args []string) error {
	var incidentAlerts []Alert
	var alerts []Alert
	var incidentID string

	// Create a new pagerduty client
	client, err := pdcli.NewConnection().Build()

	if err != nil {
		return err
	}

	// Check for incident ID argument
	if len(args) > 0 {
		incidentID = strings.TrimSpace(args[0])

		// Validate the incident ID
		match, _ := regexp.MatchString(constants.IncidentIdRegex, incidentID)

		if !match {
			return fmt.Errorf("invalid incident ID")
		}

		// Show alerts for the given incident
		selectAlert(client, incidentID)

	}

	// Fetch all incidents
	incidents, err := GetIncidents(client)

	if err != nil {
		return err
	}

	for _, incident := range incidents {

		// An incident can have more than one alert
		incidentAlerts, err = GetIncidentAlerts(client, incident.Id)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	// Check for interactive mode
	if options.interactive {

		err = selectIncident(client)

		if err != nil {
			return err
		}

	} else {
		printAlerts(alerts)
	}

	return nil
}

// GetIncidents returns an array pagerduty incidents.
func GetIncidents(c client.PagerDutyClient) ([]pdApi.Incident, error) {

	var status []string
	var teams []string
	var users []string

	var opts pdApi.ListIncidentsOptions

	var aerr pdApi.APIError

	// Get current user details
	user, err := c.GetCurrentUser(pdApi.GetCurrentUserOptions{})

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

	// Check the assigned-to flag
	switch options.assignment {

	case "team":
		// Fetch incidents belonging to a specific team
		opts.TeamIDs = append(teams, constants.TeamID)

	case "silentTest":
		// Fetch incidents assigned to silent test
		opts.UserIDs = append(users, constants.SilentTest)

	case "self":
		// Fetch incidents only assigned to self
		opts.UserIDs = append(users, user.ID)
	}

	// Fetch only triggered, acknowledged incidents (not resolved ones)
	opts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged)

	// Let the number of incidents fetched
	opts.Limit = constants.AlertsLimit

	// Check urgency
	if options.low {
		opts.Urgencies = []string{"low"}
	} else if options.high {
		opts.Urgencies = []string{"high"}
	}

	// Get incidents via pagerduty API
	incidents, err := c.ListIncidents(opts)

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

	// Check if there are no incidents returned
	if len(incidents.Incidents) == 0 {
		fmt.Println("Currently there are no alerts assigned to " + options.assignment)
		os.Exit(1)
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

	for _, p := range incidentAlerts.Alerts {

		tempAlertObj := Alert{}

		// Check if the alert is not resolved
		if p.Status != constants.StatusResolved {
			tempAlertObj.ParseAlertData(&p)
			alerts = append(alerts, tempAlertObj)
		}

	}

	return alerts, nil
}

// GetAlertData parses a pagerduty alert data into the Alert struct.
func (a *Alert) ParseAlertData(alert *pdApi.IncidentAlert) (err error) {

	// check if the alert is of type 'missing cluster'
	isCHGM := alert.Body["details"].(map[string]interface{})["notes"]

	if isCHGM != nil {
		notes := strings.Split(fmt.Sprint(alert.Body["details"].(map[string]interface{})["notes"]), "\n")

		a.ClusterID = strings.Replace(notes[0], "cluster_id: ", "", 1)
		a.ClusterName = strings.Split(fmt.Sprint(alert.Body["details"].(map[string]interface{})["name"]), ".")[0]

		lastCheckIn := fmt.Sprint(alert.Body["details"].(map[string]interface{})["last healthy check-in"])
		a.LastCheckIn, err = formatTimestamp(lastCheckIn)

		if err != nil {
			return err
		}

		a.Token = fmt.Sprint(alert.Body["details"].(map[string]interface{})["token"])
		a.Tags = fmt.Sprint(alert.Body["details"].(map[string]interface{})["tags"])
		a.Sop = strings.Replace(notes[1], "runbook: ", "", 1)

	} else {
		a.ClusterID = fmt.Sprint(alert.Body["details"].(map[string]interface{})["cluster_id"])
		a.ClusterName = strings.Split(fmt.Sprint(alert.Service.Summary), ".")[0]
		a.Console = fmt.Sprint(alert.Body["details"].(map[string]interface{})["console"])
		a.Labels = fmt.Sprint(alert.Body["details"].(map[string]interface{})["firing"])
		a.Sop = fmt.Sprint(alert.Body["details"].(map[string]interface{})["link"])
	}

	a.IncidentID = alert.Incident.ID
	a.AlertID = alert.ID
	a.Name = alert.Summary
	a.Severity = alert.Severity
	a.Status = alert.Status

	return nil

}

// GetAlertMetadata returns the alert details of a particular incident and alert.
func GetAlertMetadata(c client.PagerDutyClient, incidentID, alertID string) (*Alert, error) {
	var alertData Alert

	alert, response, err := c.GetIncidentAlert(incidentID, alertID)

	if err != nil {
		return nil, err
	}

	// check for http status code error
	if response.StatusCode != 200 {
		err = fmt.Errorf("error: %v, Status Code: %v", response.Body, response.StatusCode)
		return nil, err
	}

	alertData.ParseAlertData(alert.IncidentAlert)

	return &alertData, nil
}

// selectIncident lists incidents in interactive mode.
func selectIncident(c client.PagerDutyClient) error {
	var items []string

	incidents, err := GetIncidents(c)

	if err != nil {
		return err
	}

	for _, i := range incidents {
		incident := i.Id + " " + i.Title + " " + i.Urgency
		items = append(items, incident)
	}

	items = append(items, "exit")

	// list incidents in interactive mode
	prompt := promptui.Select{
		Label: "Select incident",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	// Exit from interactive mode
	if result == "exit" {
		os.Exit(0)
	}

	// Fetch only the incident ID
	incidentID := strings.Split(result, " ")[0]

	// list alerts for a selected incident
	err = selectAlert(c, incidentID)

	if err != nil {
		return err
	}

	return nil
}

// selectAlert prompts the user to select an alert in interactive mode.
func selectAlert(c client.PagerDutyClient, incidentID string) error {
	var items []string
	var alertData *Alert

	alerts, err := GetIncidentAlerts(c, incidentID)

	if err != nil {
		return err
	}

	for _, v := range alerts {
		alert := v.AlertID + " " + "[" + v.ClusterName + "]" + " " + v.Name + " " + v.Severity
		items = append(items, alert)
	}

	items = append(items, "exit")

	// list alerts in interactive mode
	prompt := promptui.Select{
		Label: "Select alert",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	// Exit from interactive mode
	if result == "exit" {
		os.Exit(0)
	}

	// Fetch only the alert ID
	alertID := strings.Split(result, " ")[0]

	// Fetch the metadata of a given alert
	alertData, err = GetAlertMetadata(c, incidentID, alertID)

	if err != nil {
		return err
	}

	printAlertMetadata(alertData)

	promptClusterLogin(c, alertData)

	return nil
}

// promptClusterLogin prompts the user for cluster login, if yes spawns an instance of ocm-container
func promptClusterLogin(c client.PagerDutyClient, alert *Alert) error {

	template := &promptui.SelectTemplates{
		Label: "{{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     "Do you want to proceed with cluster login?",
		Items:     []string{"Yes", "No"},
		Templates: template,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	// If a user chooses not to log into the cluster
	if result == "No" {
		os.Exit(0)
	}

	// Check if ocm-container is installed locally
	ocmContainer, err := exec.LookPath("ocm-container")

	if err != nil {
		fmt.Println("ocm-container is not found.\nPlease install it via:", constants.OcmContainerURL)
	}

	cmd := exec.Command(ocmContainer, alert.ClusterID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	// If the command exits, switch control flow back to incident selection
	if cmd.ProcessState.Exited() {
		selectIncident(c)
	}

	return nil
}

// printAlerts prints all the alerts to the console in a tabular format.
func printAlerts(alerts []Alert) {
	var headers []string

	// columns returned by the columns flag
	columns := strings.Split(options.columns, ",")

	columnsMap := make(map[string]bool)

	for _, c := range columns {
		columnsMap[c] = true
	}

	// Initializing a new table printer
	table := output.NewTable(true)

	headersMap := make(map[int]string)

	for _, alert := range alerts {

		var values []string

		var i int

		if columnsMap["incident.id"] {
			i++
			headersMap[i] = "INCIDENT ID"
			values = append(values, alert.IncidentID)
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

		table.AddRow(values)
	}

	keys := make([]int, 0)

	for k := range headersMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, v := range keys {
		headers = append(headers, headersMap[v])
	}

	table.SetHeaders(headers)
	table.SetData()
	table.Print()
}

// printAlertMetadata prints the alert data to the console.
func printAlertMetadata(alert *Alert) {

	if alert.ClusterID != "" {
		fmt.Printf("* Cluster ID: %s\n", alert.ClusterID)
	}

	if alert.ClusterName != "" {
		fmt.Printf("* Cluster Name: %s\n", alert.ClusterName)
	}

	if alert.Console != "" {
		fmt.Printf("* Console: %s\n", alert.Console)
	}

	if alert.LastCheckIn != "" {
		fmt.Printf("* Last Healthy Check-in: %s\n", alert.LastCheckIn)
	}

	if alert.Tags != "" {
		fmt.Printf("* Tags: %s\n", alert.Tags)
	}

	if alert.Token != "" {
		fmt.Printf("* Token: %s\n", alert.Token)
	}

	if alert.Labels != "" {
		fmt.Printf("* %s", alert.Labels)
	}

	if alert.Sop != "" {
		fmt.Printf("* SOP: %s\n", alert.Sop)
	}
}

// formatTimestamp formats a given timestamp into a UTC format time and returns the string.
func formatTimestamp(timestamp string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", timestamp)

	if err != nil {
		return "", err
	}

	return t.Format("01-02-2006 15:04 UTC"), nil
}
