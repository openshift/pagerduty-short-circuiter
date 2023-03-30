package terminal

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// The application.
var app = tview.NewApplication()
var pages = tview.NewPages()
var info = tview.NewTextView()
var inputBuffer []rune

func main() {
	// App Shorcuts Implementation
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Move to the next Slide
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
			return nil
			// Move to the Previous Slide
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
			return nil
			// Add a new Slide
		} else if event.Key() == tcell.KeyCtrlA {
			addSlide()
			return nil
			// Delete the current active Slide
		} else if event.Key() == tcell.KeyCtrlE {
			slideNum, _ := strconv.Atoi(info.GetHighlights()[0])
			removeSlide(slideNum)
			return nil
			// TODO : Handle the buffer with more edge cases
			// Handling Backspace with input buffer
		} else if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
			}
			// Working on the input buffer
		} else if event.Key() == tcell.KeyRune {
			inputBuffer = append(inputBuffer, event.Rune())
			// Exit the current slide when exit command is typed
		} else if event.Key() == tcell.KeyEnter {
			if string(inputBuffer) == "exit" {
				inputBuffer = []rune{}
				slideNum, _ := strconv.Atoi(info.GetHighlights()[0])
				removeSlide(slideNum)
			}
			inputBuffer = []rune{}
		}
		return event
	})

	// Get the initial Config
	layout := initTerminalMux()
	// Start the application.
	if err := app.SetRoot(layout, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
