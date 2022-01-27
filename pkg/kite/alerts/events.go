package kite

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

var (
	ClusterID         string
	AckIncidents      []string
	SelectedIncidents map[string]string
)

// SetAlertsTableEvents is the event handler for the alerts table.
// It handles the program flow when a table selection is made.
func SetAlertsTableEvents(alerts []Alert, tui *ui.TUI) {
	tui.Table.SetSelectedFunc(func(row int, column int) {
		alertID := tui.Table.GetCell(row, 1).Text

		for _, alert := range alerts {
			if alertID == alert.AlertID {
				utils.InfoLogger.Printf("GET: fetching alert metadata for alert ID: %s", alertID)
				alertData := ParseAlertMetaData(alert)
				ClusterID = alert.ClusterID

				// Do not prompt for cluster login if there's no cluster ID associated with the alert (v3 clusters)
				if ClusterID != "N/A" && ClusterID != "" && alertData != "" {
					tui.SecondaryWindow.SetText(fmt.Sprintf("Press 'Y' to log into the cluster: %s", alert.ClusterName)).SetTextColor(ui.PromptTextColor)
				}

				tui.AlertMetadata.SetText(alertData)
				break
			}
		}

		tui.Pages.AddAndSwitchToPage(ui.AlertDataPageTitle, tui.AlertMetadata, true)
		tui.Footer.SetText(ui.FooterText)
	})
}

// SetIncidentsTableEvents is the event handler for the incidents table in ack mode.
// It handles the program flow when a table selection is made.
func SetIncidentsTableEvents(tui *ui.TUI) {
	SelectedIncidents = make(map[string]string)
	tui.IncidentsTable.SetSelectedFunc(func(row, column int) {

		incidentID := tui.IncidentsTable.GetCell(row, 0).Text

		if _, ok := SelectedIncidents[incidentID]; !ok || SelectedIncidents[incidentID] == "" {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorLimeGreen)
			SelectedIncidents[incidentID] = incidentID
			utils.InfoLogger.Printf("Selected incident: %s", incidentID)
		} else {
			tui.IncidentsTable.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
			SelectedIncidents[incidentID] = ""
			utils.InfoLogger.Printf("Deselected incident: %s", incidentID)
		}
	})
}

// AcknowledgeSelectedIncidents acknowledges the selected incidents.
// All the incidents that have been acknowledged are printed to the secondary view.
func AckowledgeSelectedIncidents(tui *ui.TUI, client *client.PDClient, pdUser User) {
	utils.InfoLogger.Printf("PUT: acknowledging incidents: %v", AckIncidents)
	ackIncidents, err := AcknowledgeIncidents(client, AckIncidents, pdUser)

	if err != nil {
		utils.ErrorLogger.Printf("%v", err)
		return
	}

	for _, v := range ackIncidents {
		utils.InfoLogger.Printf("Incident %s has been acknowledged", v.Id)
	}

	var i int

	// Remove ack'ed alerts from table
	for i < tui.IncidentsTable.GetRowCount() {
		for _, v := range AckIncidents {
			if tui.IncidentsTable.GetCell(i, 0).Text == v {
				tui.IncidentsTable.RemoveRow(i)
			}
		}

		i++
	}

	AckIncidents = []string{}

	// Refresh Page
	SetIncidentsTableEvents(tui)
	tui.Pages.SwitchToPage(ui.IncidentsPageTitle)
}
