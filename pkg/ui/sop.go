package ui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/rivo/tview"
)

// App Setup
func ViewAlertSOP(tui *TUI, URL string) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			tui.App.Draw()
		})
	utils.FetchHTMLContent(URL, textView)
	name := "SOP - " + utils.GetReadmePath(URL)
	textView.Highlight("0").SetBorder(true).SetTitle(name)
	AddSOPSlide(name, textView, tui)
	// Input Handling
	textView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := textView.GetHighlights()
		if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyEnter {
				// TODO: Update with the Link
				url := textView.GetRegionText(currentSelection[0])
				utils.FetchHTMLContent(url, textView)
				fmt.Println(url)
			}
			if key == tcell.KeyTab {
				index = (index + 1) % tui.NumLinks
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + tui.NumLinks) % tui.NumLinks
			} else {
				return
			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})

}
