package kite

import (
	"fmt"

	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
)

// InitAlertsSecondaryView initializes the text content to be rendered in the secondary view.
func InitAlertsSecondaryView(tui *ui.TUI, pdUser User) {
	tui.SecondaryWindow.SetText(
		fmt.Sprintf("Logged in user: %s\n\nViewing alerts assigned to: %s\n\nPagerDuty role: %s",
			pdUser.Name,
			pdUser.AssignedTo,
			pdUser.Role)).
		SetTextColor(ui.InfoTextColor)
}

// InitAlertsUI initializes TUI table primitive.
// It adds the returned table as a new TUI page view.
func InitAlertsUI(alerts []Alert, tableTitle string, pageTitle string, tui *ui.TUI) {
	headers, data := GetAlertsTableData(alerts)
	tui.Table = tui.InitTable(headers, data, true, false, tableTitle)
	tui.Pages.AddPage(pageTitle, tui.Table, true, true)
	SetAlertsTableEvents(alerts, tui)
	tui.FrontPage = pageTitle
}

// InitIncidentsUI initializes TUI IncidentsTable primitive.
func InitIncidentsUI(incidents [][]string, tui *ui.TUI) {
	incidentHeaders := []string{"INCIDENT ID", "NAME", "SEVERITY", "STATUS", "SERVICE", "ASSIGNED TO"}
	tui.IncidentsTable = tui.InitTable(incidentHeaders, incidents, true, true, ui.IncidentsTableTitle)
	tui.Pages.AddPage(ui.IncidentsPageTitle, tui.IncidentsTable, true, false)
	SetIncidentsTableEvents(tui)
}

// InitIncidentsUI initializes TUI AckIncidentsTable primitive.
func InitAckIncidentsUI(incidents [][]string, tui *ui.TUI) {
	incidentHeaders := []string{"INCIDENT ID", "NAME", "SEVERITY", "STATUS", "SERVICE", "ACKNOWLEDGED BY"}
	tui.AckIncidentsTable = tui.InitTable(incidentHeaders, incidents, false, false, ui.AckIncidentsTableTitle)
	tui.Pages.AddPage(ui.AckIncidentsPageTitle, tui.AckIncidentsTable, true, false)
}
