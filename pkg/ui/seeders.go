package ui

import (
	"strconv"

	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
)

// SeedHighAlertsUI fetches trigerred alerts with status high and initializes a TUI table/page component.
func (tui *TUI) SeedHighAlertsUI() {
	var triggeredHigh []pdcli.Alert

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

	tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged}
	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		tui.showError(err.Error())
	}

	for _, i := range incidents {
		incident := []string{i.Id, i.Title, i.Urgency, i.Status, i.Service.Summary}
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

	tui.IncidentOpts.Statuses = []string{constants.StatusTriggered}
	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		tui.showError(err.Error())
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

	//Refresh triggered and resolved alerts
	pdcli.TrigerredAlerts = []pdcli.Alert{}
	pdcli.ResolvedAlerts = []pdcli.Alert{}

	if tui.AssginedTo == tui.Username {
		tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged}
	} else {
		tui.IncidentOpts.Statuses = []string{constants.StatusAcknowledged, constants.StatusTriggered}
	}

	incidents, err := pdcli.GetIncidents(tui.Client, &tui.IncidentOpts)

	if err != nil {
		tui.showError(err.Error())
	}

	for _, incident := range incidents {
		incidentAlerts, err = pdcli.GetIncidentAlerts(tui.Client, incident)

		if err != nil {
			tui.showError(err.Error())
		}

		alerts = append(alerts, incidentAlerts...)
	}

	// Total alerts retreived
	tui.FetchedAlerts = strconv.Itoa(len(alerts))

	tui.Alerts = alerts

	tui.InitAlertsUI(tui.Alerts, AlertsTableTitle, AlertsPageTitle)
}
