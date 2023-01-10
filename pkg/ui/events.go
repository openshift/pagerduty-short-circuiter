package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

// SetAlertsTableEvents is the event handler for the alerts table.
// It handles the program flow when a table selection is made.
func (tui *TUI) SetAlertsTableEvents(alerts []pdcli.Alert) {
	tui.Table.SetSelectedFunc(func(row int, column int) {
		var clusterName string
		var alertData string

		alertID := tui.Table.GetCell(row, 1).Text

		for _, alert := range alerts {
			if alertID == alert.AlertID {
				utils.InfoLogger.Printf("GET: fetching alert metadata for alert ID: %s", alertID)
				alertData = pdcli.ParseAlertMetaData(alert)
				clusterName = alert.ClusterName
				tui.ClusterID = alert.ClusterID
				break
			}
		}

		tui.AlertMetadata.SetText(alertData)
		tui.Pages.AddAndSwitchToPage(AlertDataPageTitle, tui.AlertMetadata, true)

		// Do not prompt for cluster login if there's no cluster ID associated with the alert (v3 clusters)
		if tui.ClusterID != "N/A" && tui.ClusterID != "" && alertData != "" {
			tui.SecondaryWindow.SetText(fmt.Sprintf("Press 'Y' to log into the cluster: %s", clusterName)).SetTextColor(PromptTextColor)
		}
	})
}

// SetIncidentsTableEvents is the event handler for the incidents table in ack mode.
// It handles the program flow when a table selection is made.
func (tui *TUI) SetIncidentsTableEvents() {
	tui.SelectedIncidents = make(map[string]string)
	tui.IncidentsTable.SetSelectedFunc(func(row, column int) {

		incidentID := tui.IncidentsTable.GetCell(row, 0).Text

		if _, ok := tui.SelectedIncidents[incidentID]; !ok || tui.SelectedIncidents[incidentID] == "" {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorLimeGreen)
			tui.SelectedIncidents[incidentID] = incidentID
			utils.InfoLogger.Printf("Selected incident: %s", incidentID)
		} else {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
			tui.SelectedIncidents[incidentID] = ""
			utils.InfoLogger.Printf("Deselected incident: %s", incidentID)
		}
	})
}

// acknowledgeSelectedIncidents acknowledges the selected incidents.
// All the incidents that have been acknowledged are printed to the secondary view.
func (tui *TUI) ackowledgeSelectedIncidents() {
	utils.InfoLogger.Printf("PUT: acknowledging incidents: %v", tui.AckIncidents)
	ackIncidents, err := pdcli.AcknowledgeIncidents(tui.Client, tui.AckIncidents)

	if err != nil {
		utils.ErrorLogger.Printf("%v", err)
		return
	}

	for _, v := range ackIncidents {
		utils.InfoLogger.Printf("Incident %s has been acknowledged", v.APIObject.ID)
	}

	var i int

	// Remove ack'ed alerts from table
	for i < tui.IncidentsTable.GetRowCount() {
		for _, v := range tui.AckIncidents {
			if tui.IncidentsTable.GetCell(i, 0).Text == v {
				tui.IncidentsTable.RemoveRow(i)
			}
		}

		i++
	}

	tui.AckIncidents = []string{}

	// Refresh Page
	tui.SetIncidentsTableEvents()
	tui.Pages.SwitchToPage(IncidentsPageTitle)
}
