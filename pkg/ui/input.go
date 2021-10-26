package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
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

			return nil
		}

		if event.Rune() == 'q' || event.Rune() == 'Q' {
			tui.App.Stop()
		}

		if event.Rune() == 'm' || event.Rune() == 'M' {
			tui.toggleMouse()
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
					tui.AckIncidents = []string{}
				}
			}

			return event
		})
	}

	tui.AlertMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'Y' || event.Rune() == 'y' {
			tui.App.Stop()

			hasExited, _ := pdcli.ClusterLogin(tui.ClusterID)

			if hasExited {
				tui.Init()
				tui.Pages.AddAndSwitchToPage(AlertsPageTitle, tui.Table, true)
				err := tui.StartApp()

				if err != nil {
					panic(err)
				}
			}
		}

		return event
	})

}

// toggleMouse enables & disables mouse events in TUI.
func (tui *TUI) toggleMouse() {
	if tui.isMouseEnabled {
		tui.App.EnableMouse(false)
	} else {
		tui.App.EnableMouse(true)
	}

	tui.isMouseEnabled = !tui.isMouseEnabled
}
