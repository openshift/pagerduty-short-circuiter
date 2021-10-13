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
	"regexp"
	"sort"
	"strconv"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/spf13/cobra"
)

var options struct {
	high       bool
	low        bool
	assignment string
	columns    string
	incidentID bool
	status     string
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
		"incident.id,alert.id,cluster.name,alert,cluster.id,status,severity",
		"Specify which columns to display separated by commas without any space in between",
	)

	// Alerts status
	Cmd.Flags().StringVar(
		&options.status,
		"status",
		"trigerred",
		"Filter alerts by status",
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

	var tui ui.TUI

	// Create a new pagerduty client
	client, err := client.NewClient().Connect()

	if err != nil {
		return err
	}

	// Fetch the currently logged in user's ID.
	user, err := client.GetCurrentUser(pdApi.GetCurrentUserOptions{})

	if err != nil {
		return err
	}

	// UI internals
	tui.Client = client
	tui.Username = user.Name

	// Check for incident ID argument
	if len(args) > 0 {
		incidentID = strings.TrimSpace(args[0])

		// Validate the incident ID
		match, _ := regexp.MatchString(constants.IncidentIdRegex, incidentID)

		if !match {
			return fmt.Errorf("invalid incident ID")
		}

		alerts, err := pdcli.GetIncidentAlerts(client, incidentID)

		if err != nil {
			return err
		}

		tui.FetchedAlerts = strconv.Itoa(len(alerts))

		tui.Init()

		initAlertsUI(&tui, alerts, incidentID+" "+ui.AlertsTableTitle)

		err = tui.StartApp()

		if err != nil {
			return err
		}

		return nil
	}

	// Set the limit on incidents fetched
	incidentOpts.Limit = constants.IncidentsLimit

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
		incidentOpts.UserIDs = append(users, user.ID)

	default:
		return fmt.Errorf("please enter a valid assigned-to option")
	}

	// Check urgency
	if options.low {
		incidentOpts.Urgencies = []string{"low"}
	} else if options.high {
		incidentOpts.Urgencies = []string{"high"}
	}

	// Check the status flag
	switch options.status {

	case "trigerred":
		// Fetch trigerred incidents
		incidentOpts.Statuses = append(status, constants.StatusTriggered)

	case "ack":
		// Fetch incidents that have been acknowledged
		incidentOpts.Statuses = append(users, constants.StatusAcknowledged)

	case "resolved":
		// Fetch resolved incidents
		incidentOpts.Statuses = append(users, constants.StatusResolved)

	default:
		return fmt.Errorf("please enter a valid status")
	}

	// Fetch incidents
	incidents, err := pdcli.GetIncidents(client, &incidentOpts)

	if err != nil {
		return err
	}

	// Check if there are no incidents returned
	if len(incidents) == 0 {
		fmt.Println("Currently there are no alerts assigned to " + options.assignment)
		os.Exit(0)
	}

	// Parse incident data to TUI
	for _, i := range incidents {
		if i.Status != constants.StatusAcknowledged {
			incident := []string{i.Id, i.Title, i.Urgency, i.Status, i.Service.Summary}
			tui.Incidents = append(tui.Incidents, incident)
		}
	}

	// Get incident alerts
	for _, incident := range incidents {

		// An incident can have more than one alert
		incidentAlerts, err = pdcli.GetIncidentAlerts(client, incident.Id)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	tui.AssginedTo = options.assignment
	tui.FetchedAlerts = strconv.Itoa(len(alerts))

	// Setup TUI
	tui.Init()
	initAlertsUI(&tui, alerts, ui.AlertsTableTitle)
	initIncidentsUI(&tui, client)

	err = tui.StartApp()

	if err != nil {
		return err
	}

	return nil
}

// initAlertsUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func initAlertsUI(tui *ui.TUI, alerts []pdcli.Alert, title string) {
	headers, data := getTableData(alerts)
	tui.Table = tui.InitTable(headers, data, true, false, title)
	tui.SetAlertsTableEvents(alerts)
	tui.SetAlertsSecondaryData()

	tui.Pages.AddPage(ui.AlertsPageTitle, tui.Table, true, true)
}

// initIncidentsUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func initIncidentsUI(tui *ui.TUI, c client.PagerDutyClient) {
	incidentHeaders := []string{"INCIDENT ID", "NAME", "SEVERITY", "STATUS", "SERVICE"}
	tui.IncidentsTable = tui.InitTable(incidentHeaders, tui.Incidents, true, true, ui.IncidentsTableTitle)
	tui.SetIncidentsTableEvents()

	tui.Pages.AddPage(ui.AckIncidentsPageTitle, tui.IncidentsTable, true, false)
}

// getTableData parses and returns tabular data for the given alerts, i.e table headers and rows.
func getTableData(alerts []pdcli.Alert) ([]string, [][]string) {
	var headers []string
	var tableData [][]string

	// columns returned by the columns flag
	columns := strings.Split(options.columns, ",")

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
