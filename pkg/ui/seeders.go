package ui

import (
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
)

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
