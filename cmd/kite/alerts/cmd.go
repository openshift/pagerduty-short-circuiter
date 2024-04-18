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
	"regexp"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
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
	Short: "This command will list all the open high alerts assigned to self or selected team .",
	Args:  cobra.MaximumNArgs(1),
	RunE:  alertsHandler,
}

func init() {

	// Incident Assignment
	Cmd.Flags().StringVar(
		&options.assignment,
		"assigned-to",
		"self",
		"View alerts assigned to your logged in PagerDuty user account (default), selected team or Silent Test.",
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

	// Setup TUI
	tui.Init()
	utils.InfoLogger.Print("Initialized terminal UI")

	if utils.Emulator != "" {
		utils.InfoLogger.Printf("Terminal emulator for cluster login set to: %s", utils.Emulator)
	} else {
		utils.ErrorLogger.Printf("No terminal emulator found")
	}

	// Create a new pagerduty client
	utils.InfoLogger.Print("Connecting to PagerDuty API")
	client, err := client.NewClient().Connect()

	if err != nil {
		return err
	}
	utils.InfoLogger.Print("Connection successful")

	// Fetch the currently logged in user's ID.
	utils.InfoLogger.Print("GET: fetching logged in user data")
	user, err := client.GetCurrentUser(pdApi.GetCurrentUserOptions{})

	if err != nil {
		return err
	}

	// UI internals
	tui.Client = client
	tui.Username = user.Name
	tui.Columns = options.columns
	tui.Role = user.Role

	// Check for incident ID argument
	if len(args) > 0 {
		incidentID = strings.TrimSpace(args[0])

		// Validate the incident ID
		match, _ := regexp.MatchString(constants.IncidentIdRegex, incidentID)

		if !match {
			return fmt.Errorf("invalid incident ID")
		}

		// Create PD Incident Object with given ID
		incident := pdApi.Incident{
			APIObject: pdApi.APIObject{
				ID: incidentID,
			},
		}

		utils.InfoLogger.Printf("GET: fetching incident alerts for incident ID: %s", incident.APIObject.ID)
		alerts, err := pdcli.GetIncidentAlerts(client, incident)

		if err != nil {
			return err
		}

		tui.Alerts = alerts

		utils.InfoLogger.Print("Initializing alerts view")
		tui.InitAlertsUI(alerts, ui.AlertsTableTitle, ui.AlertsPageTitle)

		err = tui.StartApp()

		if err != nil {
			return err
		}

		return nil
	}

	// Set incidents urgency
	utils.InfoLogger.Printf("Incidents urgency set to: %s, %s", constants.StatusLow, constants.StatusHigh)
	incidentOpts.Urgencies = []string{constants.StatusLow, constants.StatusHigh}

	// Check the assigned-to flag
	switch options.assignment {

	case "team":
		// Load the configuration file
		cfg, err := config.Load()

		if err != nil {
			return err
		}

		teamID := cfg.TeamID
		tui.AssignedTo = cfg.Team

		if teamID == "" {
			return fmt.Errorf("no team selected, please run 'kite teams' to set a team")
		}

		// Fetch incidents belonging to a specific team
		incidentOpts.TeamIDs = append(teams, teamID)
		utils.InfoLogger.Printf("Retrieving incidents assigned to team: %s", cfg.Team)

		// Fetch incidents with the following statuses
		incidentOpts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged)
		utils.InfoLogger.Printf("Retrieving incidents with status: %s, %s", constants.StatusTriggered, constants.StatusAcknowledged)

	case "silentTest":
		// Fetch incidents assigned to silent test
		tui.AssignedTo = "Silent Test"
		incidentOpts.UserIDs = append(users, constants.SilentTest)
		utils.InfoLogger.Printf("Retrieving incidents assigned to Silent Test")

		// Fetch incidents with the following statuses
		incidentOpts.Statuses = append(status, constants.StatusTriggered, constants.StatusAcknowledged)
		utils.InfoLogger.Printf("Retrieving incidents with status: %s, %s", constants.StatusTriggered, constants.StatusAcknowledged)

	case "self":
		// Fetch incidents only assigned to self
		tui.AssignedTo = user.Name
		incidentOpts.UserIDs = append(users, user.ID)
		utils.InfoLogger.Printf("Retrieving incidents assigned to: %s", user.Name)

		// Fetch only acknowledged incidents when option is self (default)
		incidentOpts.Statuses = append(status, constants.StatusAcknowledged)
		utils.InfoLogger.Printf("Retrieving incidents with status: %s", constants.StatusAcknowledged)

	default:
		return fmt.Errorf("please enter a valid assigned-to option")
	}

	// Set the limit on incidents fetched
	utils.InfoLogger.Printf("Incidents limit set to: %d", constants.IncidentsLimit)
	incidentOpts.Limit = constants.IncidentsLimit

	// Fetch incidents
	utils.InfoLogger.Printf("GET: fetching incidents")
	incidents, err := pdcli.GetIncidents(client, &incidentOpts)

	if err != nil {
		return err
	}

	// Get incident alerts
	utils.InfoLogger.Printf("GET: fetching incident alerts")
	for _, incident := range incidents {
		// An incident can have more than one alert
		incidentAlerts, err = pdcli.GetIncidentAlerts(client, incident)

		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	tui.Alerts = alerts
	tui.IncidentOpts = incidentOpts

	tui.InitAlertsUI(tui.Alerts, ui.AlertsTableTitle, ui.AlertsPageTitle)
	tui.InitAlertsSecondaryView()
	// Start TUI
	err = tui.StartApp()

	if err != nil {
		return err
	}

	return nil
}
