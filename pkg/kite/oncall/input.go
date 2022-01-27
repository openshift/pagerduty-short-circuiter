package kite

import (
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/ui"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

// InitOncallKeyboard initializes the keyboard event handlers for all the oncall TUI primitives.
func InitOncallKeyboard(tui *ui.TUI, allTeamsOncall, nextOncall []OncallUser) {
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {
			tui.Pages.SwitchToPage(ui.OncallPageTitle)
			return nil
		}

		if event.Rune() == 'q' || event.Rune() == 'Q' {
			utils.InfoLogger.Println("Exiting kite")
			tui.App.Stop()
		}

		return event
	})

	setupOncallPageInput(tui, allTeamsOncall, nextOncall)
}

func setupOncallPageInput(tui *ui.TUI, allTeamsOncall []OncallUser, nextOncall []OncallUser) {
	if title, _ := tui.Pages.GetFrontPage(); title == ui.OncallPageTitle {
		tui.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'N' || event.Rune() == 'n' {
				utils.InfoLogger.Print("Viewing user next on-call schedule")
				InitOncallUI(nextOncall, ui.NextOncallTableTitle, ui.NextOncallPageTitle, tui)

				if len(nextOncall) == 0 {
					utils.InfoLogger.Print("You are not scheduled for any oncall duties for the next 3 months. Cheer up!")
				}
			}

			if event.Rune() == 'A' || event.Rune() == 'a' {
				utils.InfoLogger.Print("Switching to all team on-call view")
				InitOncallUI(allTeamsOncall, ui.AllTeamsOncallTableTitle, ui.AllTeamsOncallPageTitle, tui)
			}

			return event
		})
	}
}
