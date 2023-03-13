package ui

import (
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

// SeedHighAlertsUI fetches trigerred alerts with status high and initializes a TUI table/page component.
func (tui *TUI) SeedHighAlertsUI() {
	var triggeredHigh []pdcli.Alert

	utils.InfoLogger.Print("Fetching high alerts")

	for _, alert := range pdcli.TrigerredAlerts {
		if alert.Severity == constants.StatusHigh {
			triggeredHigh = append(triggeredHigh, alert)
		}
	}

	tui.InitAlertsUI(triggeredHigh, HighAlertsTableTitle, HighAlertsPageTitle)
}

// SeedHLowAlertsUI fetches trigerred alerts with status low and initializes a TUI table/page component.
func (tui *TUI) SeedHLowAlertsUI() {
	var triggeredLow []pdcli.Alert

	utils.InfoLogger.Print("Fetching low alerts")

	for _, alert := range pdcli.TrigerredAlerts {
		if alert.Severity == constants.StatusLow {
			triggeredLow = append(triggeredLow, alert)
		}
	}

	tui.InitAlertsUI(triggeredLow, LowAlertsTableTitle, LowAlertsPageTitle)
}

// SeedAckIncidentsUI fetches acknlowedged incidents and initializes a TUI table/page component.
func (tui *TUI) SeedAckIncidentsUI() {
	var ackIncidents [][]string

	utils.InfoLogger.Printf("Incidents status set to: %s", constants.StatusAcknowledged)
	tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged}

	// Fetch new incidents via PD API
	utils.InfoLogger.Print("GET: fetching acknowledged incidents")
	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		utils.ErrorLogger.Print(err)
	}

	for _, i := range incidents {
		I := i.Assignments[0]
		incident := []string{i.Id, i.Title, i.Urgency, i.Status, i.Service.Summary, I.Assignee.Summary}
		ackIncidents = append(ackIncidents, incident)
	}

	tui.Incidents = ackIncidents

	tui.InitIncidentsUI(tui.Incidents, AckIncidentsTableTitle, AckIncidentsPageTitle, false)
	tui.Footer.SetText(FooterTextAlerts)
	tui.Pages.SwitchToPage(AckIncidentsPageTitle)
}

// SeedIncidentsUI fetches trigerred incidents and initializes a TUI table/page component.
func (tui *TUI) SeedIncidentsUI() {
	var incidentsData [][]string

	utils.InfoLogger.Printf("Incidents status set to: %s", constants.StatusTriggered)
	tui.IncidentOpts.Statuses = []string{constants.StatusTriggered}

	// Override incidents limit when viewing triggered incidents
	utils.InfoLogger.Printf("Incidents limit set to: %d", constants.TrigerredIncidentsLimit)
	tui.IncidentOpts.Limit = constants.TrigerredIncidentsLimit

	// Fetch new incidents via PD API
	utils.InfoLogger.Print("GET: fetching incidents")
	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		utils.ErrorLogger.Print(err)
	}

	for _, i := range incidents {
		incident := []string{i.Id, i.Title, i.Urgency, i.Status, i.Service.Summary}
		incidentsData = append(incidentsData, incident)
	}

	tui.Incidents = incidentsData

	tui.InitIncidentsUI(tui.Incidents, IncidentsTableTitle, IncidentsPageTitle, true)
	tui.Footer.SetText(FooterTextIncidents)
}

// SeedIncidentsUI fetches acknowledged incident alerts and initializes a TUI table/page component.
func (tui *TUI) SeedAlertsUI() {
	var incidentAlerts []pdcli.Alert
	var alerts []pdcli.Alert

	// Refresh triggered and resolved alerts
	pdcli.TrigerredAlerts = []pdcli.Alert{}
	pdcli.ResolvedAlerts = []pdcli.Alert{}

	if tui.AssignedTo == tui.Username {
		utils.InfoLogger.Printf("Incidents status set to: %s", constants.StatusAcknowledged)
		tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged}
	} else {
		utils.InfoLogger.Printf("Incidents status set to: %s, %s", constants.StatusTriggered, constants.StatusAcknowledged)
		tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged, constants.StatusTriggered}
	}

	// Fetch new incidents via PD API
	utils.InfoLogger.Print("GET: fetching incidents")
	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		utils.ErrorLogger.Print(err)
	}

	// Fetch new incident alerts via PD API
	utils.InfoLogger.Print("GET: fetching incident alerts")
	for _, incident := range incidents {
		incidentAlerts, err = pdcli.GetIncidentAlerts(tui.Client, incident)

		if err != nil {
			utils.ErrorLogger.Print(err)
		}

		alerts = append(alerts, incidentAlerts...)
	}

	tui.Alerts = alerts

	tui.InitAlertsUI(tui.Alerts, AlertsTableTitle, AlertsPageTitle)
}
