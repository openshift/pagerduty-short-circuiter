package ui

import (
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

	// Fetch new incidents
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
}

// SeedIncidentsUI fetches trigerred incidents and initializes a TUI table/page component.
func (tui *TUI) SeedIncidentsUI() {

	var incidentsData [][]string

	tui.IncidentOpts.Statuses = []string{constants.StatusTriggered}

	// Override incidents limit when viewing triggered incidents assigned to self
	if tui.AssginedTo == tui.Username {
		tui.IncidentOpts.Limit = 25
	}

	// Fetch new incidents from PD
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
