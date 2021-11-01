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
	"strconv"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
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
}

// alertsHandler is the main alerts command handler.
func alertsHandler(cmd *cobra.Command, args []string) error {

	var (
		// Internals
		incidentAlerts []pdcli.Alert
		alerts         []pdcli.Alert
		incidentID     string
		incidentOpts   pdApi.ListIncidentsOptions
		teams          []string
		users          []string
		status         []string

		//UI
		tui ui.TUI
	)

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
	tui.Columns = options.columns

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

		tui.InitAlertsUI(alerts, ui.AlertsTableTitle, ui.AlertsPageTitle)

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
		// Load the configuration file
		cfg, err := config.Load()

		if err != nil {
			return err
		}

		teamID := cfg.TeamID
		tui.AssginedTo = cfg.Team

		if teamID == "" {
			return fmt.Errorf("no team selected, please run 'pdcli teams' to set a team")
		}

		// Fetch incidents belonging to a specific team
		incidentOpts.TeamIDs = append(teams, teamID)

		// Fetch incidents with the following statuses
		incidentOpts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged, constants.StatusResolved)

	case "silentTest":
		// Fetch incidents assigned to silent test
		incidentOpts.UserIDs = append(users, constants.SilentTest)
		tui.AssginedTo = "Silent Test"

		// Fetch incidents with the following statuses
		incidentOpts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged, constants.StatusResolved)

	case "self":
		// Fetch incidents only assigned to self
		incidentOpts.UserIDs = append(users, user.ID)
		tui.AssginedTo = user.Name

		// Fetch only acknowledged incidents when option is self (default)
		incidentOpts.Statuses = append(status, constants.StatusAcknowledged, constants.StatusResolved)

	default:
		return fmt.Errorf("please enter a valid assigned-to option")
	}

	// Check urgency
	if options.low {
		incidentOpts.Urgencies = []string{"low"}
	} else if options.high {
		incidentOpts.Urgencies = []string{"high"}
	}

	tui.IncidentOpts = incidentOpts

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

	// Get incident alerts
	for _, incident := range incidents {

		// An incident can have more than one alert
		incidentAlerts, err = pdcli.GetIncidentAlerts(client, incident.Id)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	// Total alerts retreived
	tui.FetchedAlerts = strconv.Itoa(len(alerts))

	// Determine terminal emulator for cluster login
	pdcli.InitTerminalEmulator()

	if pdcli.Terminal != "" {
		tui.HasEmulator = true
	}

	// Setup TUI
	tui.Init()
	tui.InitAlertsUI(alerts, ui.AlertsTableTitle, ui.AlertsPageTitle)

	err = tui.StartApp()

	if err != nil {
		return err
	}

	return nil
}
