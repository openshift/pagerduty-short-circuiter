package kite

import (
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

func InitAlertsKeyboard(tui *ui.TUI, client *client.PDClient, lowAlerts, highAlerts []Alert, pdUser User) {
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// If escape key is pressed
		if event.Key() == tcell.KeyEscape {
			if page, _ := tui.Pages.GetFrontPage(); page != ui.AlertsPageTitle {

				if page == ui.HighAlertsPageTitle || page == ui.LowAlertsPageTitle {
					tui.FrontPage = ui.AlertsPageTitle
				}

				// If the user is viewing the alert metadata
				if page == ui.AlertDataPageTitle {
					InitAlertsSecondaryView(tui, pdUser)
				}

				tui.Footer.SetText(ui.FooterTextAlerts)
				tui.Pages.SwitchToPage(tui.FrontPage)
			}
		}

		if event.Rune() == 'q' || event.Rune() == 'Q' {
			utils.InfoLogger.Println("Exiting kite")
			tui.App.Stop()
		}

		setupAlertsPageInput(tui, lowAlerts, highAlerts)
		setupIncidentsPageInput(tui, client, pdUser)
		setupAlertDataPageInput(tui)

		return event
	})
}

func setupAlertsPageInput(tui *ui.TUI, lowAlerts, highAlerts []Alert) {
	if page, _ := tui.Pages.GetFrontPage(); page == ui.AlertsPageTitle {
		tui.Pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			if event.Rune() == '1' {
				if page, _ := tui.Pages.GetFrontPage(); page != ui.AckIncidentsPageTitle {
					tui.Footer.SetText(ui.FooterText)
					utils.InfoLogger.Print("Viewing acknowledged incidents")
				}

				tui.Pages.SwitchToPage(ui.AckIncidentsPageTitle)
			}

			if event.Rune() == '2' {
				if page, _ := tui.Pages.GetFrontPage(); page != ui.IncidentsPageTitle {
					tui.Footer.SetText(ui.FooterTextIncidents)
					utils.InfoLogger.Print("Viewing triggerred incidents")
				}

				tui.Pages.SwitchToPage(ui.IncidentsPageTitle)
			}

			// Filter by alerts severity
			if event.Rune() == 'H' || event.Rune() == 'h' {
				if page, _ := tui.Pages.GetFrontPage(); page != ui.HighAlertsPageTitle {
					utils.InfoLogger.Print("Switching to high alerts view")
				}

				InitAlertsUI(highAlerts, ui.HighAlertsTableTitle, ui.HighAlertsPageTitle, tui)
			}

			if event.Rune() == 'L' || event.Rune() == 'l' {
				if page, _ := tui.Pages.GetFrontPage(); page != ui.LowAlertsPageTitle {
					utils.InfoLogger.Print("Switching to low alerts view")
				}

				InitAlertsUI(lowAlerts, ui.LowAlertsTableTitle, ui.LowAlertsPageTitle, tui)
			}

			return event
		})
	}
}

func setupAlertDataPageInput(tui *ui.TUI) {
	tui.AlertMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'Y' || event.Rune() == 'y' {

			if utils.Emulator != "" {
				err := utils.ClusterLoginEmulator(ClusterID)

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

			} else {
				tui.App.Stop()

				cmd := utils.ClusterLoginShell(ClusterID)

				err := cmd.Start()

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

				err = cmd.Wait()

				tui.Init()

				if err != nil {
					utils.ErrorLogger.Print(err)
				}

				err = tui.StartApp()

				if err != nil {
					panic(err)
				}
			}
		}

		return event
	})
}

func setupIncidentsPageInput(tui *ui.TUI, client *client.PDClient, pdUser User) {
	if page, _ := tui.Pages.GetFrontPage(); page == ui.IncidentsPageTitle {
		tui.Pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlA {
				for _, v := range SelectedIncidents {
					if v != "" {
						AckIncidents = append(AckIncidents, v)
					}
				}

				if len(AckIncidents) == 0 {
					utils.ErrorLogger.Print("Please select atleast one incident to acknowledge")
				} else {
					AckowledgeSelectedIncidents(tui, client, pdUser)
				}
			}

			return event
		})
	}
}
