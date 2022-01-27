package kite

import (
	"fmt"

	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
)

// InitOnCallSecondaryView initializes the text content to be rendered in the secondary view.
func InitOnCallSecondaryView(user string, primary string, secondary string, tui *ui.TUI) {
	tui.SecondaryWindow.SetText(
		fmt.Sprintf("Logged in user: %s\n\nPrimary on-call: %s\n\nSecondary on-call: %s",
			user,
			primary,
			secondary),
	)
}

// InitOncallUI initializes TUI table component.
// It adds the returned table as a new TUI page component.
func InitOncallUI(onCallData []OncallUser, tableTitle string, pageTitle string, tui *ui.TUI) {
	headers, data := GetOncallTableData(onCallData)
	tui.Table = tui.InitTable(headers, data, false, false, tableTitle)
	tui.Pages.AddPage(pageTitle, tui.Table, true, true)
}

// GetOncallTableData parses and returns tabular data for the given oncall data, i.e table headers and rows.
func GetOncallTableData(oncallData []OncallUser) ([]string, [][]string) {
	var tableData [][]string

	for _, v := range oncallData {
		tableData = append(tableData, []string{v.EscalationPolicy, v.Name, v.OncallRole, v.Start, v.End})
	}

	headers := []string{"Escalation Policy", "Name", "Oncall Role", "From", "To"}

	return headers, tableData
}
