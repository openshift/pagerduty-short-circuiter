package ui

import (
	"github.com/gdamore/tcell/v2"

	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

// initKeyboard initializes the keyboard event handlers for all the TUI components.
func (tui *TUI) initKeyboard() {
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {
			// Check if alerts command is executed
			if tui.Pages.HasPage(AlertsPageTitle) {
				tui.InitAlertsSecondaryView()
				page, _ := tui.Pages.GetFrontPage()

				// If the user is viewing the alert metadata
				if page == AlertDataPageTitle {
					tui.Pages.SwitchToPage(tui.FrontPage)
				} else if page == AlertMetadata {
					tui.Pages.SwitchToPage(IncidentsPageTitle)
				} else {
					tui.InitAlertsUI(tui.Alerts, AlertsTableTitle, AlertsPageTitle)
					tui.Pages.SwitchToPage(AlertsPageTitle)
					tui.Footer.SetText(FooterTextAlerts)
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
			utils.InfoLogger.Println("Exiting kite")
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
				utils.InfoLogger.Print("Switching to acknowledged incidents view")
				tui.SeedAckIncidentsUI()

				if len(tui.Incidents) == 0 {
					utils.InfoLogger.Printf("No acknowledged incidents assigned found")
				}

				tui.Pages.SwitchToPage(AckIncidentsPageTitle)
			}

			if event.Rune() == '2' {
				utils.InfoLogger.Print("Switching to incidents view")
				tui.SeedIncidentsUI()

				if len(tui.Incidents) == 0 {
					utils.InfoLogger.Printf("No trigerred incidents assigned to found")
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
				utils.InfoLogger.Print("Refreshing alerts...")
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
					utils.ErrorLogger.Print("Please select atleast one incident to acknowledge")
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

			if utils.Emulator != "" {
				err := utils.ClusterLoginEmulator(tui.ClusterID)

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

			} else {
				tui.App.Stop()

				cmd := utils.ClusterLoginShell(tui.ClusterID)

				err := cmd.Start()

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

				err = cmd.Wait()

				tui.Init()

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

				// Refresh alerts table
				tui.SeedAlertsUI()
				utils.InfoLogger.Print("Switching back to alerts view")
				tui.Pages.SwitchToPage(AlertsPageTitle)

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
					utils.InfoLogger.Print("Viewing user next on-call schedule")
					tui.Pages.SwitchToPage(NextOncallPageTitle)

					if len(tui.AckIncidents) == 0 {
						utils.InfoLogger.Print("You are not scheduled for any oncall duties for the next 3 months. Cheer up!")
					}
				}
			}

			if tui.AllTeamsOncallTable != nil {
				if event.Rune() == 'A' || event.Rune() == 'a' {
					utils.InfoLogger.Print("Switching to all team on-call view")
					tui.Pages.SwitchToPage(AllTeamsOncallPageTitle)
				}
			}

			return event
		})
	}
}
