package ui

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/gdamore/tcell/v2"

	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
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

				// Handle page traversal
				switch page {
				case AlertDataPageTitle:
					tui.Pages.SwitchToPage(tui.FrontPage)
				case ServiceLogsPageTitle:
					tui.Pages.SwitchToPage(AlertDataPageTitle)
					tui.InitAlertDataSecondaryView()
				default:
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
		if event.Key() == tcell.KeyLeft && CursorPos > 0 {
			CursorPos--
			return event
		} else if event.Key() == tcell.KeyRight && CursorPos < len(tui.TerminalInputBuffer) {
			CursorPos++
			return event
		} else if event.Key() == tcell.KeyCtrlN {
			NextSlide(tui)
			return nil
			// Move to the Previous Slide
		} else if event.Key() == tcell.KeyCtrlP {
			PreviousSlide(tui)
			return nil
			// Add a new Slide - bash
		} else if event.Key() == tcell.KeyCtrlA {
			AddNewSlide(tui, constants.Shell, os.Getenv("SHELL"), []string{}, false)
			return nil
			// Add a new Slide - ocm-container
		} else if event.Key() == tcell.KeyCtrlO {
			OcmContainerPath, err := exec.LookPath(constants.OcmContainer)
			if err != nil {
				utils.ErrorLogger.Println("ocm-container is not found.\nPlease install it via:", constants.OcmContainerURL)
				return nil
			}
			AddNewSlide(tui, constants.OcmContainer, OcmContainerPath, []string{}, false)
			return nil

		} else if event.Key() == tcell.KeyCtrlB {
			// Reset the input buffer
			tui.TerminalInputBuffer = []rune{}
			return nil
		} else if event.Rune() >= '0' && event.Rune() <= '9' && len(tui.TerminalInputBuffer) < 2 {
			// Append the digit to the input buffer
			tui.TerminalInputBuffer = append(tui.TerminalInputBuffer, event.Rune())
			// Check if the input buffer is complete
			if len(tui.TerminalInputBuffer) == 1 {
				// Convert the input buffer to a slide number
				slideNum, err := strconv.Atoi(string(tui.TerminalInputBuffer))
				if err == nil && slideNum >= 1 && slideNum <= len(tui.TerminalTabs) {
					// Navigate to the specified slide
					tui.TerminalPageBar.Highlight(strconv.Itoa(tui.TerminalUIRegionIDs[slideNum-1])).
						ScrollToHighlight()
					CurrentActivePage = slideNum - 1
				}
				// Reset the input buffer
				tui.TerminalInputBuffer = []rune{}
			}
			return nil
			// Delete the current active Slide
		} else if event.Key() == tcell.KeyCtrlE {
			slideNum, _ := strconv.Atoi(tui.TerminalPageBar.GetHighlights()[0])
			RemoveSlide(slideNum, tui)
			tui.TerminalInputBuffer = []rune{}
			return nil
		} else if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if len(tui.TerminalInputBuffer) > 0 {
				if CursorPos > 0 {
					tui.TerminalInputBuffer = append(tui.TerminalInputBuffer[:CursorPos-1], tui.TerminalInputBuffer[CursorPos:]...)
					CursorPos--
				} else {
					tui.TerminalInputBuffer = tui.TerminalInputBuffer[:len(tui.TerminalInputBuffer)-1]
				}
			}
			if len(tui.TerminalLastChars) > len("exit") {
				tui.TerminalLastChars = tui.TerminalLastChars[len(tui.TerminalLastChars)-len("exit"):]
			}
			// Working on the input buffer
		} else if event.Key() == tcell.KeyRune {
			if CursorPos >= len(tui.TerminalInputBuffer) {
				// Append new rune to end of input buffer
				tui.TerminalInputBuffer = append(tui.TerminalInputBuffer, event.Rune())
				tui.TerminalLastChars = append(tui.TerminalLastChars, event.Rune())
				if len(tui.TerminalLastChars) > len("exit") {
					tui.TerminalLastChars = tui.TerminalLastChars[1:]
				}
			} else {
				// Insert new rune at cursor position in input buffer
				tui.TerminalInputBuffer = append(tui.TerminalInputBuffer[:CursorPos], append([]rune{event.Rune()}, tui.TerminalInputBuffer[CursorPos:]...)...)
				tui.TerminalLastChars = append(tui.TerminalLastChars[:CursorPos], append([]rune{event.Rune()}, tui.TerminalLastChars[CursorPos:]...)...)
				if len(tui.TerminalLastChars) > len("exit") {
					tui.TerminalLastChars = tui.TerminalLastChars[1:]
				}
			}
			CursorPos++

		} else if event.Key() == tcell.KeyEnter {
			if string(tui.TerminalInputBuffer) == "exit" || string(tui.TerminalLastChars) == "exit" {
				slideNum, _ := strconv.Atoi(tui.TerminalPageBar.GetHighlights()[0])
				RemoveSlide(slideNum, tui)
			}
			tui.TerminalInputBuffer = []rune{}
			tui.TerminalLastChars = []rune{}
			CursorPos = 0
		}

		// Override the default exit behaviour with Ctrl+C
		if event.Key() == tcell.KeyCtrlC {
			return nil
		}
		// Exit the App on Ctrl + Q
		if event.Key() == tcell.KeyCtrlQ {
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
				utils.InfoLogger.Print("Switching to trigerred alerts view")
				tui.InitAlertsUI(pdcli.TrigerredAlerts, TrigerredAlertsTableTitle, TrigerredAlertsPageTitle)
			}

			if event.Rune() == '2' {
				utils.InfoLogger.Print("Switching to acknowledged incidents view")
				tui.SeedAckIncidentsUI()

				if len(tui.Incidents) == 0 {
					utils.InfoLogger.Printf("No acknowledged incidents assigned found")
				}

				tui.Pages.SwitchToPage(AckIncidentsPageTitle)
			}

			if event.Rune() == '3' {
				utils.InfoLogger.Print("Switching to incidents view")
				tui.SeedIncidentsUI()

				if len(tui.Incidents) == 0 {
					utils.InfoLogger.Printf("No trigerred incidents assigned to found")
				}

				tui.Pages.SwitchToPage(IncidentsPageTitle)
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
			// Get ocm-conatiner executable from PATH
			ocmContainer, err := exec.LookPath("ocm-container")

			if err != nil {
				errMessage := "ocm-container is not found.\nPlease install it via: " + constants.OcmContainerURL
				utils.ErrorLogger.Print(errMessage)
				return nil
			}

			// Convert the ClusterID into args for ocm-container command
			clusterIDArgs := []string{tui.ClusterID}
			AddNewSlide(tui, tui.ClusterName, ocmContainer, clusterIDArgs, true)
		}

		if event.Rune() == 'L' || event.Rune() == 'l' {
			utils.InfoLogger.Print("Retrieving service logs for cluster")
			tui.fetchClusterServiceLogs()
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
					tui.Footer.SetText(FooterText)

					if len(tui.AckIncidents) == 0 {
						utils.InfoLogger.Print("You are not scheduled for any oncall duties for the next 3 months. Cheer up!")
					}
				}
			}

			if tui.AllTeamsOncallTable != nil {
				if event.Rune() == 'A' || event.Rune() == 'a' {
					utils.InfoLogger.Print("Switching to all team on-call view")
					tui.Pages.SwitchToPage(AllTeamsOncallPageTitle)
					tui.Footer.SetText(FooterText)
				}
			}

			return event
		})
	}
}
