package ui

import (
	"github.com/gdamore/tcell/v2"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
)

// initKeyboard initializes the keyboard event handlers for all the TUI components.
func (tui *TUI) initKeyboard() {

	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {

			if tui.Pages.HasPage(AlertsPageTitle) {

				tui.Footer.SetText(FooterTextAlerts)

				tui.Pages.SwitchToPage(AlertsPageTitle)

				tui.showDefaultSecondaryView()
			}

			if tui.Pages.HasPage(OncallPageTitle) {
				tui.Pages.SwitchToPage(OncallPageTitle)
				tui.Footer.SetText(FooterTextOncall)
			}

			return nil
		}

		if event.Rune() == 'q' || event.Rune() == 'Q' {
			tui.App.Stop()
		}

		return event
	})

	tui.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if tui.IncidentsTable != nil {
			if event.Rune() == 'A' || event.Rune() == 'a' {
				if len(tui.Incidents) == 0 {
					tui.showSecondaryView("No incidents assigned to " + tui.Username + " found.")
				}

				tui.Footer.SetText(FooterTextAck)

				tui.Pages.SwitchToPage(AckIncidentsPageTitle)
			}
		} else {
			tui.Footer.SetText(FooterText)
		}

		if tui.NextOncallTable != nil {
			if event.Rune() == 'N' || event.Rune() == 'n' {
				tui.Pages.SwitchToPage(NextOncallPageTitle)

				if len(tui.AckIncidents) == 0 {
					tui.showSecondaryView("You are not scheduled for any oncall duties for the next 3 months. Cheer up!")
				}
			}
		}

		if tui.AllTeamsOncallTable != nil {
			if event.Rune() == 'A' || event.Rune() == 'a' {
				tui.Pages.SwitchToPage(AllTeamsOncallPageTitle)
			}
		}

		return event
	})

	if tui.IncidentsTable != nil {
		tui.IncidentsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlA {

				for _, v := range tui.SelectedIncidents {
					if v != "" {
						tui.AckIncidents = append(tui.AckIncidents, v)
					}
				}

				if len(tui.AckIncidents) == 0 {
					tui.showSecondaryView("Please select atleast one incident to acknowledge")
				} else {
					tui.ackowledgeSelectedIncidents()
				}
			}

			return event
		})
	}

	tui.AlertMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'Y' || event.Rune() == 'y' {

			// if tui.HasEmulator {
			// 	err := pdcli.ClusterLoginEmulator(tui.ClusterID)

			// 	if err != nil {
			// 		tui.showError(err.Error())
			// 	}
			// } else {
			tui.App.Stop()

			cmd := pdcli.ClusterLoginShell(tui.ClusterID)

			err := cmd.Start()

			if err != nil {
				tui.showError(err.Error())
			}

			err = cmd.Wait()

			tui.Init()

			if err != nil {
				tui.showError(err.Error())
			}

			tui.Pages.AddAndSwitchToPage(AlertsPageTitle, tui.Table, true)
			tui.Pages.AddPage(AckIncidentsPageTitle, tui.IncidentsTable, true, false)

			err = tui.StartApp()

			if err != nil {
				panic(err)
			}

		}
		//}

		return event
	})

}
