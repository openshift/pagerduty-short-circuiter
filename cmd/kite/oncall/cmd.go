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
	kite "github.com/openshift/pagerduty-short-circuiter/pkg/kite/oncall"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "oncall to the PagerDuty CLI",
	Long:  "Running the kite oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	RunE:  oncallHandler,
}

// oncallHandler is the main handler for kite oncall.
func oncallHandler(cmd *cobra.Command, args []string) (err error) {
	var (
		onCallUsers    []kite.OncallUser
		allTeamsOncall []kite.OncallUser
		nextOncall     []kite.OncallUser
		primary        string
		secondary      string
		tui            ui.TUI
	)

	// Initialize TUI
	tui.Init()
	utils.InfoLogger.Print("Initialized terminal UI")

	// Establish a secure connection with the PagerDuty API
	utils.InfoLogger.Print("Connecting to PagerDuty API")
	client, err := client.NewClient().Connect()

	if err != nil {
		return err
	}
	utils.InfoLogger.Print("Connection successful")

	// Fetch the currently logged in user's ID.
	utils.InfoLogger.Print("GET: fetching logged in user data")
	user, err := client.GetCurrentUser(pagerduty.GetCurrentUserOptions{})

	if err != nil {
		return err
	}

	tui.Username = user.Name

	// Fetch oncall data from Platform-SRE team
	utils.InfoLogger.Print("GET: fetching on-call data of current user team")
	onCallUsers, err = kite.TeamSREOnCall(client)

	if err != nil {
		return err
	}

	for _, v := range onCallUsers {
		if strings.Contains(v.OncallRole, "Primary") {
			primary = v.Name
		}

		if strings.Contains(v.OncallRole, "Secondary") {
			secondary = v.Name
		}
	}

	// Fetch oncall data from all teams
	utils.InfoLogger.Print("GET: fetching on-call data of all teams")
	allTeamsOncall, err = kite.AllTeamsOncall(client)

	if err != nil {
		return err
	}

	// Fetch the current user's oncall schedule
	utils.InfoLogger.Print("GET: fetching next on-call schedule of logged in user")
	nextOncall, err = kite.UserNextOncallSchedule(client, user.ID)

	if err != nil {
		return err
	}

	utils.InfoLogger.Print("Initializing on-call view")
	initOncallUI(&tui, onCallUsers)

	utils.InfoLogger.Print("Initializing all teams on-call view")
	initAllTeamsOncallUI(&tui, allTeamsOncall)

	utils.InfoLogger.Print("Initializing next on-call schedule view")
	initNextOncallUI(&tui, nextOncall)

	utils.InfoLogger.Print("Initializing secondary view")
	tui.InitOnCallSecondaryView(user.Name, primary, secondary)

	err = tui.StartApp()

	if err != nil {
		return err
	}

	return nil
}

// initOncallUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func initOncallUI(tui *ui.TUI, onCallData []kite.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.Table = tui.InitTable(headers, data, false, false, ui.OncallTableTitle)
	tui.Pages.AddPage(ui.OncallPageTitle, tui.Table, true, true)
}

// initOncallUI initializes TUI NextOncall table component.
// It adds the returned table as a new TUI page view.
func initNextOncallUI(tui *ui.TUI, onCallData []kite.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.NextOncallTable = tui.InitTable(headers, data, false, false, ui.NextOncallTableTitle)
	tui.Pages.AddPage(ui.NextOncallPageTitle, tui.NextOncallTable, true, false)
}

// initOncallUI initializes TUI AllTeamsOncall table component.
// It adds the returned table as a new TUI page view.
func initAllTeamsOncallUI(tui *ui.TUI, onCallData []kite.OncallUser) {
	headers, data := getOncallTableData(onCallData)
	tui.AllTeamsOncallTable = tui.InitTable(headers, data, false, false, ui.AllTeamsOncallTableTitle)
	tui.Pages.AddPage(ui.AllTeamsOncallPageTitle, tui.AllTeamsOncallTable, true, false)
}

// getOncallTableData parses and returns tabular data for the given oncall data, i.e table headers and rows.
func getOncallTableData(oncallData []kite.OncallUser) ([]string, [][]string) {
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
