/*
Copyright © 2021 Red Hat, Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oncall

import (
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/oncall"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "oncall to the PagerDuty CLI",
	Long:  "Running the pdcli oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	RunE:  oncallHandler,
}

// oncallHandler is the main handler for pdcli oncall.
func oncallHandler(cmd *cobra.Command, args []string) (err error) {

	var onCallUsers []pdcli.OncallUser
	var allTeamsOncall []pdcli.OncallUser
	var nextOncall []pdcli.OncallUser

	// Establish a secure connection with the PagerDuty API
	client, err := client.NewClient().Connect()

	if err != nil {
		return err
	}

	user, err := client.GetCurrentUser(pagerduty.GetCurrentUserOptions{})

	if err != nil {
		return err
	}

	// UI
	var tui ui.TUI

	tui.Username = user.Name

	// Initialize TUI
	tui.Init()

	// Fetch oncall data from Platform-SRE team
	onCallUsers, err = pdcli.TeamSREOnCall(client)

	if err != nil {
		return err
	}

	for _, v := range onCallUsers {
		if strings.Contains(v.OncallRole, "Primary") {
			tui.Primary = v.Name
		}

		if strings.Contains(v.OncallRole, "Secondary") {
			tui.Secondary = v.Name
		}
	}

	// Fetch oncall data from all teams
	allTeamsOncall, err = pdcli.AllTeamsOncall(client)

	if err != nil {
		return err
	}

	// Fetch the current user's oncall schedule
	nextOncall, err = pdcli.UserNextOncallSchedule(client, user.ID)

	if err != nil {
		return err
	}

	initOncallUI(&tui, onCallUsers)
	initAllTeamsOncallUI(&tui, allTeamsOncall)
	initNextOncallUI(&tui, nextOncall)

	tui.SetOncallSecondaryData()

	err = tui.StartApp()

	if err != nil {
		return err
	}

	return nil
}

// initOncallUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func initOncallUI(tui *ui.TUI, onCallData []pdcli.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.Table = tui.InitTable(headers, data, false, false, ui.OncallTableTitle)
	tui.Pages.AddPage(ui.OncallPageTitle, tui.Table, true, true)
}

// initOncallUI initializes TUI NextOncall table component.
// It adds the returned table as a new TUI page view.
func initNextOncallUI(tui *ui.TUI, onCallData []pdcli.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.NextOncallTable = tui.InitTable(headers, data, false, false, ui.NextOncallTableTitle)
	tui.Pages.AddPage(ui.NextOncallPageTitle, tui.NextOncallTable, true, false)
}

// initOncallUI initializes TUI AllTeamsOncall table component.
// It adds the returned table as a new TUI page view.
func initAllTeamsOncallUI(tui *ui.TUI, onCallData []pdcli.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.AllTeamsOncallTable = tui.InitTable(headers, data, false, false, ui.AllTeamsOncallTableTitle)
	tui.Pages.AddPage(ui.AllTeamsOncallPageTitle, tui.AllTeamsOncallTable, true, false)
}

// getOncallTableData parses and returns tabular data for the given oncall data, i.e table headers and rows.
func getOncallTableData(oncallData []pdcli.OncallUser) ([]string, [][]string) {

	var tableData [][]string

	for _, v := range oncallData {

		var data []string

		if v.EscalationPolicy != "" {
			data = append(data, v.EscalationPolicy)
		} else {
			data = append(data, "N/A")
		}

		if v.Name != "" {
			data = append(data, v.Name)
		} else {
			data = append(data, "N/A")
		}

		if v.OncallRole != "" {
			data = append(data, v.OncallRole)
		} else {
			data = append(data, "N/A")
		}

		if v.Start != "" {
			data = append(data, v.Start)
		} else {
			data = append(data, "N/A")
		}

		if v.End != "" {
			data = append(data, v.End)
		} else {
			data = append(data, "N/A")
		}

		tableData = append(tableData, data)
	}

	headers := []string{"Escalation Policy", "Name", "Oncall Role", "From", "To"}

	return headers, tableData
}
