package ui

import (
	"fmt"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
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

		var incident pdApi.Incident
		client, _ := client.NewClient().Connect()
		incidentID := tui.IncidentsTable.GetCell(row, 0).Text
		incident.Id = incidentID
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
			tui.Pages.AddAndSwitchToPage(AlertMetadata, tui.AlertMetadata, true)

		} else {
			tui.SetAlertsTableEvents(alerts)
			tui.InitAlertsUI(alerts, AlertMetadata, AlertMetadata)

		}
		// Do not prompt for cluster login if there's no cluster ID associated with the alert (v3 clusters)
		if tui.ClusterID != "N/A" && tui.ClusterID != "" && alertData != "" {
			tui.SecondaryWindow.SetText(fmt.Sprintf("Press 'Y' to log into the cluster: %s", clusterName)).SetTextColor(PromptTextColor)
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
		utils.InfoLogger.Printf("Incident %s has been acknowledged", v.Id)
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
