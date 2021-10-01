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
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/manifoldco/promptui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/output"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

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

	var (
		incidentAlerts []pdcli.Alert
		incidentID     string
		incidentOpts   pdApi.ListIncidentsOptions
		alerts         []pdcli.Alert
		teams          []string
		users          []string
		status         []string
	)

	// Create a new pagerduty client
	client, err := client.NewClient().Connect()

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
		selectAlert(client, incidentID, &incidentOpts)
	}

	// Set the limit on incidents fetched
	incidentOpts.Limit = constants.AlertsLimit

	// Fetch only triggered, acknowledged incidents (not resolved ones)
	incidentOpts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged)

	// Fetch the currently logged in user's ID.
	userID, err := pdcli.GetCurrentUserID(client)

	if err != nil {
		return err
	}

	// Check the assigned-to flag
	switch options.assignment {

	case "team":
		// Fetch incidents belonging to a specific team
		incidentOpts.TeamIDs = append(teams, constants.TeamID)

	case "silentTest":
		// Fetch incidents assigned to silent test
		incidentOpts.UserIDs = append(users, constants.SilentTest)

	case "self":
		// Fetch incidents only assigned to self
		incidentOpts.UserIDs = append(users, userID)
	}

	// Check urgency
	if options.low {
		incidentOpts.Urgencies = []string{"low"}
	} else if options.high {
		incidentOpts.Urgencies = []string{"high"}
	}

	// Fetch all incidents
	incidents, err := pdcli.GetIncidents(client, &incidentOpts)

	if err != nil {
		return err
	}

	// Check if there are no incidents returned
	if len(incidents) == 0 {
		fmt.Println("Currently there are no alerts assigned to " + options.assignment)
		os.Exit(0)
	}

	for _, incident := range incidents {

		// An incident can have more than one alert
		incidentAlerts, err = pdcli.GetIncidentAlerts(client, incident.Id)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	// Check for interactive mode
	if options.interactive {
		err = selectIncident(client, &incidentOpts)

		if err != nil {
			return err
		}

	} else {
		printAlerts(alerts)
	}

	return nil
}

// selectIncident lists incidents in interactive mode.
func selectIncident(c client.PagerDutyClient, opts *pdApi.ListIncidentsOptions) error {
	var items []string

	incidents, err := pdcli.GetIncidents(c, opts)

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
	err = selectAlert(c, incidentID, opts)

	if err != nil {
		return err
	}

	return nil
}

// selectAlert prompts the user to select an alert in interactive mode.
func selectAlert(c client.PagerDutyClient, incidentID string, opts *pdApi.ListIncidentsOptions) error {
	var items []string
	var alertData pdcli.Alert

	alerts, err := pdcli.GetIncidentAlerts(c, incidentID)

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

	// Fetch the selected alert data from the already fetched alerts slice
	for _, v := range alerts {
		if v.AlertID == alertID {
			alertData = v
		}
	}

	printAlertMetadata(&alertData)

	promptClusterLogin(c, &alertData, opts)

	return nil
}

// promptClusterLogin prompts the user for cluster login, if yes spawns an instance of ocm-container
func promptClusterLogin(c client.PagerDutyClient, alert *pdcli.Alert, opts *pdApi.ListIncidentsOptions) error {

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
		selectIncident(c, opts)
	}

	return nil
}

// printAlerts prints all the alerts to the console in a tabular format.
func printAlerts(alerts []pdcli.Alert) {
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
func printAlertMetadata(alert *pdcli.Alert) {

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

	if alert.WebURL != "" {
		fmt.Printf("* Web URL: %s\n", alert.WebURL)
	}
}
