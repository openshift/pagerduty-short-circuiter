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
	kite "github.com/openshift/pagerduty-short-circuiter/pkg/kite/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/spf13/cobra"
)

var options struct {
	assignment string
	incidentID bool
}

var Cmd = &cobra.Command{
	Use:   "alerts",
	Short: "This command will list all the open high alerts assigned to self.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  alertsHandler,
}

func init() {
	// Incident Assignment
	Cmd.Flags().StringVar(
		&options.assignment,
		"assigned-to",
		"self",
		"Filter alerts assigned-to self/team/silentTest",
	)
}

// alertsHandler is the main alerts command handler.
func alertsHandler(cmd *cobra.Command, args []string) error {
	var (
		// Internals
		alerts       []kite.Alert
		incidentID   string
		incidentOpts pdApi.ListIncidentsOptions

		//UI
		tui ui.TUI
	)

	// Setup TUI
	tui.Init()
	fmt.Println("Initializing terminal UI")
	utils.InfoLogger.Print("Initialized terminal UI")

	// Determine terminal emulator for cluster login
	utils.InitTerminalEmulator()

	if utils.Emulator != "" {
		utils.InfoLogger.Printf("Terminal emulator for cluster login set to: %s", utils.Emulator)
	} else {
		utils.ErrorLogger.Printf("No terminal emulator found")
	}

	// Create a new pagerduty client
	fmt.Println("Connecting to PagerDuty API")
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

	// Initialize a new user object to be used across
	pdUser := kite.User{
		UserID: user.ID,
		Name:   user.Name,
		Role:   user.Role,
		Email:  user.Email,
	}

	fmt.Println("Fetching alerts...")

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
			Id: incidentID,
		}

		utils.InfoLogger.Printf("GET: fetching incident alerts for incident ID: %s", incident.Id)
		alerts, err := kite.GetIncidentAlerts(client, incident)
		if err != nil {
			return err
		}

		utils.InfoLogger.Print("Initializing alerts view")
		kite.InitAlertsUI(alerts, ui.AlertsTableTitle, ui.AlertsPageTitle, &tui)

		err = tui.StartApp()
		if err != nil {
			return err
		}

		return nil
	}

	// Set incidents urgency
	urgency := []string{constants.StatusLow, constants.StatusHigh}
	incidentOpts.Urgencies = append(incidentOpts.Urgencies, urgency...)
	utils.InfoLogger.Printf("Retrieving incidents with urgency: %v", incidentOpts.Urgencies)

	// Set incidents status
	status := []string{constants.StatusAcknowledged, constants.StatusTriggered}
	incidentOpts.Statuses = append(incidentOpts.Statuses, status...)
	utils.InfoLogger.Printf("Retrieving incidents with status: %v", incidentOpts.Statuses)

	// Set the limit on incidents fetched
	utils.InfoLogger.Printf("Incidents limit set to: %d", constants.IncidentsLimit)
	incidentOpts.Limit = constants.IncidentsLimit

	// Sort incidents by urgency
	incidentOpts.SortBy = "urgency:DESC"

	// Check the assigned-to flag
	switch options.assignment {

	case "team":
		// Load the configuration file
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		pdUser.AssignedTo = cfg.Team

		if cfg.TeamID == "" {
			return fmt.Errorf("no team selected, please run 'kite teams' to set a team")
		}

		// Fetch incidents belonging to a specific team
		utils.InfoLogger.Printf("Retrieving incidents assigned to team: %s", cfg.Team)
		incidentOpts.TeamIDs = []string{cfg.TeamID}

	case "silentTest":
		// Fetch incidents assigned to silent test
		pdUser.AssignedTo = "Silent Test"
		utils.InfoLogger.Printf("Retrieving incidents assigned to Silent Test")
		incidentOpts.UserIDs = []string{constants.SilentTest}

	case "self":
		// Fetch incidents only assigned to self
		pdUser.AssignedTo = "Self"
		utils.InfoLogger.Printf("Retrieving incidents assigned to: %s", pdUser.Name)
		incidentOpts.UserIDs = []string{pdUser.UserID}

	default:
		return fmt.Errorf("please enter a valid assigned-to option")
	}

	// Get incidents
	utils.InfoLogger.Printf("GET: fetching incidents")
	incidents, err := kite.GetIncidents(client, &incidentOpts)
	if err != nil {
		return err
	}

	// Filter incidents by status
	utils.InfoLogger.Print("Filtering incidents by status")
	ackIncidents, triggeredIncidents := kite.FilterIncidentsByStatus(incidents)

	// Format incidents into tabular rows
	ackIncidentsData, triggeredIncidentsData := kite.GetIncidentsTableData(ackIncidents, triggeredIncidents)

	// If viewing alerts assigned to logged-in user fetch only ack'd incident alerts
	if options.assignment == "self" {
		incidents = ackIncidents
	}

	utils.InfoLogger.Print("GET: fetching incident alerts")
	for _, incident := range incidents {
		// An incident can have more than one alert
		incidentAlerts, err := kite.GetIncidentAlerts(client, incident)
		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	// Filter alerts by status
	utils.InfoLogger.Print("Filtering alerts by status")
	lowAlerts, highAlerts := kite.FilterAlertsByStatus(alerts)

	// Get incident alerts & filter incidents by status
	utils.InfoLogger.Print("GET: fetching incident alerts")
	for _, incident := range incidents {
		// An incident can have more than one alert
		incidentAlerts, err := kite.GetIncidentAlerts(client, incident)
		if err != nil {
			return err
		}

		alerts = append(alerts, incidentAlerts...)
	}

	if len(alerts) == 0 && options.assignment == "self" {
		utils.InfoLogger.Printf("No acknowledged alerts for user %s found", pdUser.Name)
	}

	kite.InitAlertsUI(alerts, ui.AlertsTableTitle, ui.AlertsPageTitle, &tui)
	kite.InitAckIncidentsUI(ackIncidentsData, &tui)
	kite.InitIncidentsUI(triggeredIncidentsData, &tui)
	kite.InitAlertsSecondaryView(&tui, pdUser)
	kite.InitAlertsKeyboard(&tui, client, lowAlerts, highAlerts, pdUser)

	// Start TUI
	utils.InfoLogger.Print("Initializing TUI")
	fmt.Println("Starting TUI")
	err = tui.StartApp()
	if err != nil {
		return err
	}

	return nil
}
