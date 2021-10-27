package ui

import (
	"github.com/gdamore/tcell/v2"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
)

// SetAlertsTableEvents is the event handler for the alerts table.
// It handles the program flow when a table selection is made.
func (tui *TUI) SetAlertsTableEvents(alerts []pdcli.Alert) {
	tui.Table.SetSelectedFunc(func(row int, column int) {
		alertID := tui.Table.GetCell(row, 1).Text

		for _, alert := range alerts {
			if alertID == alert.AlertID {
				alertData := pdcli.ParseAlertMetaData(alert)
				tui.AlertMetadata.SetText(alertData)
				tui.ClusterID = alert.ClusterID
			}
		}

		tui.Pages.AddAndSwitchToPage(AlertDataPageTitle, tui.AlertMetadata, true)
		tui.promptSecondaryView("Press 'Y' to proceed with cluster login")
	})

}

// SetIncidentsTableEvents is the event handler for the incidents table in ack mode.
// It handles the program flow when a table selection is made.
func (tui *TUI) SetIncidentsTableEvents() {
	tui.SelectedIncidents = make(map[string]string)
	var text bool
	tui.IncidentsTable.SetSelectedFunc(func(row, column int) {

		incidentID := tui.IncidentsTable.GetCell(row, 0).Text

		text = !text

		if _, ok := tui.SelectedIncidents[incidentID]; !ok || tui.SelectedIncidents[incidentID] == "" {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorLimeGreen)
			tui.SelectedIncidents[incidentID] = incidentID
		} else {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
			tui.SelectedIncidents[incidentID] = ""
		}
	})
}

// acknowledgeSelectedIncidents acknowledges the selected incidents.
// All the incidents that have been acknowledged are printed to the secondary view.
func (tui *TUI) ackowledgeSelectedIncidents() {
	ackIncidents, err := pdcli.AcknowledgeIncidents(tui.Client, tui.AckIncidents)

	if err != nil {
		tui.showError(err.Error())
		return
	}

	text := "The following incidents have been acknowledged:\n"

	for _, v := range ackIncidents {
		text = text + v.ID + " - " + v.Title + "\n"
	}

	tui.showSecondaryView(text)
}
