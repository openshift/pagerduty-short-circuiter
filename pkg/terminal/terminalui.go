package terminal

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EUI struct {
	app  *tview.Application
	list *tview.List
}

var (
	selectedTerminal string
	selected         bool
)

func (eui *EUI) UiEmulator(terminals []string) string {
	eui.app = tview.NewApplication()

	eui.list = tview.NewList().ShowSecondaryText(false)

	// Add the available terminal emulators to the list.
	for i, t := range terminals {
		func(t string) {
			eui.list.AddItem(fmt.Sprintf("%d. %s", i+1, t), "", 0, func() {
				selectedTerminal = t
				selected = true
				eui.app.Stop()
			})
		}(t)
	}

	eui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {

			selected = false
			eui.app.Stop()
			return nil
		}
		return event
	})

	// Set the list as the root of the application and focus on it.
	eui.app.SetRoot(eui.list, true).SetFocus(eui.list)

	// Add a footer with a quit message.
	footer := tview.NewTextView().SetText("Q [Quit]")
	eui.app.SetRoot(
		tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(eui.list, 0, 1, true).
			AddItem(footer, 1, 0, false),
		true,
	)

	// Run the application and wait for the user to make a selection.
	if err := eui.app.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	if selected {
		return selectedTerminal
	}

	// Return an empty string if the user quits.
	return ""
}
