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
	"sort"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/olekukonko/tablewriter"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

type Alert struct {
	IncidentID string
	Name       string
	ClusterID  string
	Severity   string
	Status     string
}

var args struct {
	high       bool
	low        bool
	assignment string
	columns    string
}

var Cmd = &cobra.Command{
	Use:   "alerts",
	Short: "This command will list all the open high alerts assigned to self.",
	RunE:  alertsHandler,
}

func init() {

	// Urgency
	Cmd.Flags().BoolVar(&args.low, "low", false, "View all low alerts")
	Cmd.Flags().BoolVar(&args.high, "high", true, "View all high alerts")

	// Incident Assignment
	Cmd.Flags().StringVar(
		&args.assignment,
		"assigned-to",
		"self",
		"Filter alerts based on user or team",
	)

	// Columns displayed
	Cmd.Flags().StringVar(
		&args.columns,
		"columns",
		"incident.id,name,cluster.id,status,severity",
		"Specify which columns to display separated by commas without any space in between",
	)
}

// AlertsHandler is the main alerts command handler.
func alertsHandler(cmd *cobra.Command, args []string) error {
	var incidentAlerts []Alert
	var alerts []Alert

	client, err := pdcli.NewConnection().Build()

	if err != nil {
		return err
	}

	// Get incident IDs
	incidentIDs, err := GetIncidents(client)

	if err != nil {
		return err
	}

	for _, id := range incidentIDs {

		// An incident can have more than one alert
		incidentAlerts, err = GetIncidentAlerts(client, id)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	printAlerts(alerts)

	return nil
}

// getIncidents returns a string slice consisting IDs of the first 10 incidents.
func GetIncidents(c client.PagerDutyClient) ([]string, error) {

	var incidentIDs []string
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
	switch args.assignment {

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
	if args.low {
		opts.Urgencies = []string{"low"}
	} else if args.high {
		opts.Urgencies = []string{"high"}
	}

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
		fmt.Println("Currently there are no alerts assigned to " + args.assignment)
		os.Exit(1)
	}

	for _, incident := range incidents.Incidents {
		incidentIDs = append(incidentIDs, incident.Id)
	}

	return incidentIDs, nil
}

// getIncidentAlerts returns all the alerts belong to a particular incident.
func GetIncidentAlerts(c client.PagerDutyClient, incidentID string) ([]Alert, error) {

	var alerts []Alert

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

		tempAlertObj.IncidentID = incidentID
		tempAlertObj.Name = p.Summary
		tempAlertObj.Severity = p.Severity
		tempAlertObj.Status = p.Status
		tempAlertObj.ClusterID = fmt.Sprint(p.Body["details"].(map[string]interface{})["cluster_id"])

		alerts = append(alerts, tempAlertObj)
	}

	return alerts, nil
}

// printAlerts prints all the alerts to the console in a tabular format.
func printAlerts(alerts []Alert) {

	var tableData [][]string
	var headers []string

	columns := strings.Split(args.columns, ",")

	columnsMap := make(map[string]bool)
	headersMap := make(map[int]string)

	for _, c := range columns {
		columnsMap[c] = true
	}

	for _, alert := range alerts {

		var values []string

		var i int

		if columnsMap["incident.id"] {
			i++
			headersMap[i] = "INCIDENT ID"
			values = append(values, alert.IncidentID)
		}

		if columnsMap["name"] {
			i++
			headersMap[i] = "NAME"
			values = append(values, alert.Name)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(tableData)
	table.SetBorder(false)
	table.Render()
}
