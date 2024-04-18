package ui

import (
	"fmt"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ocm"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

// SetAlertsTableEvents is the event handler for the alerts table.
// It handles the program flow when a table selection is made.
func (tui *TUI) SetAlertsTableEvents(alerts []pdcli.Alert) {
	tui.Table.SetSelectedFunc(func(row int, column int) {
		var alertData string

		alertID := tui.Table.GetCell(row, 1).Text

		for _, alert := range alerts {
			if alertID == alert.AlertID {
				utils.InfoLogger.Printf("GET: fetching alert metadata for alert ID: %s", alertID)
				alertData = pdcli.ParseAlertMetaData(alert)
				tui.ClusterName = alert.ClusterName
				tui.ClusterID = alert.ClusterID
				tui.SOPLink = alert.Sop
				break
			}
		}

		tui.AlertMetadata.SetText(alertData)
		tui.Pages.AddAndSwitchToPage(AlertDataPageTitle, tui.AlertMetadata, true)

		// Do not prompt for cluster login if there's no cluster ID associated with the alert (v3 clusters)
		if tui.ClusterID != "N/A" && tui.ClusterID != "" && alertData != "" {
			secondaryWindowText := fmt.Sprintf("Press 'Y' to log into the cluster: %s\nPress 'S' to view the SOP\nPress 'L' to view service logs", tui.ClusterName)
			tui.SecondaryWindow.SetText(secondaryWindowText).SetTextColor(PromptTextColor)
		}
	})
}

// SetAcknowledgeTableEvents is the event handler for the acknowledged incidents table.
// It handles the program flow when a Enter is pressed on a incident is made.
func (tui *TUI) SetAckTableEvents() {
	tui.SelectedIncidents = make(map[string]string)
	tui.IncidentsTable.SetSelectedFunc(func(row, column int) {
		var incident pdApi.Incident
		client, _ := client.NewClient().Connect()
		incidentID := tui.IncidentsTable.GetCell(row, 0).Text
		incident.APIObject.ID = incidentID
		var clusterName string
		var alertData string

		alerts, _ := pdcli.GetIncidentAlerts(client, incident)
		Alert := alerts[0]

		for _, alert := range alerts {
			if incidentID == alert.IncidentID {
				alertData = pdcli.ParseAlertMetaData(alert)
				clusterName = alert.ClusterName
				tui.ClusterID = alert.ClusterID
				break
			}
		}
		if len(alerts) == 1 {
			alertData = pdcli.ParseAlertMetaData(Alert)
			tui.AlertMetadata.SetText(alertData)
			tui.Pages.AddAndSwitchToPage(AckAlertDataPage, tui.AlertMetadata, true)
			tui.Footer.SetText(FooterText)

		} else {
			tui.SetAlertsTableEvents(alerts)
			tui.InitAlertsUI(alerts, AckAlertDataPage, AckAlertDataPage)

		}
		// Do not prompt for cluster login if there's no cluster ID associated with the alert (v3 clusters)
		if tui.ClusterID != "N/A" && tui.ClusterID != "" && alertData != "" {
			secondaryWindowText := fmt.Sprintf("Press 'Y' to log into the cluster: %s\nPress 'S' to view the SOP\nPress 'L' to view service logs", clusterName)
			tui.SecondaryWindow.SetText(secondaryWindowText)
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

// fetchClusterServiceLogs returns the given cluster's service logs
// It initializes a text view and displays the parsed service log data
func (tui *TUI) fetchClusterServiceLogs() {
	var responseStr string
	serviceLogsItems, err := ocm.GetClusterServiceLogs(tui.ClusterID)
	if err != nil {
		utils.InfoLogger.Printf("%v", err)
		return
	}

	if serviceLogsItems.Empty() {
		responseStr = fmt.Sprintf("No service logs found for the cluster: %s/%s", tui.ClusterID, tui.ClusterName)
	} else {
		responseStr = ocm.ParseServiceLogItems(serviceLogsItems)
	}

	tui.ServiceLogView.SetText(responseStr)
	tui.Pages.AddAndSwitchToPage(ServiceLogsPageTitle, tui.ServiceLogView, true)
}
