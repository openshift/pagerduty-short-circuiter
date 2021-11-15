package ui

import (
	"github.com/gdamore/tcell/v2"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
)

// initKeyboard initializes the keyboard event handlers for all the TUI components.
func (tui *TUI) initKeyboard() {

	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {

			// Check if alerts command is executed
			if tui.Pages.HasPage(AlertsPageTitle) {
				page, _ := tui.Pages.GetFrontPage()

				// If the user is viewing the alert metadata
				if page == AlertDataPageTitle {
					tui.Pages.SwitchToPage(tui.FrontPage)
					tui.showDefaultSecondaryView()
				} else if page == HighAlertsPageTitle || page == LowAlertsPageTitle {
					tui.Pages.SwitchToPage(TrigerredAlertsPageTitle)
					tui.Footer.SetText(FooterTextTrigerredAlerts)
				} else {
					tui.InitAlertsUI(tui.Alerts, AlertsTableTitle, AlertsPageTitle)
					tui.Footer.SetText(FooterTextAlerts)
					tui.Pages.SwitchToPage(AlertsPageTitle)
					tui.showDefaultSecondaryView()
				}
			}

			// Check if oncall command is executed
			if tui.Pages.HasPage(OncallPageTitle) {
				tui.Pages.SwitchToPage(OncallPageTitle)
				tui.Footer.SetText(FooterTextOncall)
			}

			return nil
		}

		if event.Rune() == 'q' || event.Rune() == 'Q' {
			tui.App.Stop()
		}

		tui.setupAlertsPageInput()
		tui.setupIncidentsPageInput()
		tui.setupAlertDetailsPageInput()
		tui.setupOncallPageInput()

		return event
	})

}

func (tui *TUI) setupAlertsPageInput() {
	if title, _ := tui.Pages.GetFrontPage(); title == AlertsPageTitle {

		tui.Pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			if event.Rune() == '1' {
				tui.InitAlertsUI(pdcli.ResolvedAlerts, ResolvedAlertsTableTitle, ResolvedAlertsPageTitle)
			}

			if event.Rune() == '2' {
				tui.InitAlertsUI(pdcli.TrigerredAlerts, TrigerredAlertsTableTitle, TrigerredAlertsPageTitle)
			}

			if event.Rune() == '3' {
				tui.SeedAckIncidentsUI()

				if len(tui.Incidents) == 0 {
					tui.showSecondaryView("No acknowledged incidents assigned to " + tui.Username + " found.")
				}

				tui.Pages.SwitchToPage(AckIncidentsPageTitle)
			}

			if event.Rune() == '4' {

				tui.SeedIncidentsUI()

				if len(tui.Incidents) == 0 {
					tui.showSecondaryView("No triggered incidents assigned to " + tui.Username + " found.")
				}

				tui.Pages.SwitchToPage(IncidentsPageTitle)
			}

			// Filter by status
			if title, _ := tui.Pages.GetFrontPage(); title == TrigerredAlertsPageTitle {

				if event.Rune() == 'H' || event.Rune() == 'h' {

					tui.SeedHighAlertsUI()
					tui.Pages.SwitchToPage(HighAlertsPageTitle)
				}

				if event.Rune() == 'L' || event.Rune() == 'l' {

					tui.SeedHLowAlertsUI()
					tui.Pages.SwitchToPage(LowAlertsPageTitle)
				}
			}

			// Alerts refresh
			if event.Rune() == 'r' || event.Rune() == 'R' {
				tui.showSecondaryView("Fetching alerts...")
				tui.SeedAlertsUI()
			}

			return event
		})
	}
}

func (tui *TUI) setupIncidentsPageInput() {
	if title, _ := tui.Pages.GetFrontPage(); title == IncidentsPageTitle {
		tui.Pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
}

func (tui *TUI) setupAlertDetailsPageInput() {
	tui.AlertMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'Y' || event.Rune() == 'y' {

			if tui.HasEmulator {
				err := pdcli.ClusterLoginEmulator(tui.ClusterID)

				if err != nil {
					tui.showError(err.Error())
				}

			} else {
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
		}

		return event
	})
}

func (tui *TUI) setupOncallPageInput() {
	if title, _ := tui.Pages.GetFrontPage(); title == OncallPageTitle {
		tui.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

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
	}
}
