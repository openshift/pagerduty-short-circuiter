package terminal

import (
	"fmt"
	"os/exec"

	"git.sr.ht/~rockorager/tterm"
	"github.com/rivo/tview"
)

// Declares the tab struct
type tab struct {
	index     int
	title     string
	primitive tview.Primitive
}

var Tabs []tab
var uiRegionIds []int
var currentActivePage int = 0
var totalPageCount int = -1

// Creates and return a new tab
func newTab(name string, command string) *tab {
	totalPageCount += 1
	uiRegionIds = append(uiRegionIds, totalPageCount)
	index := len(Tabs)
	if len(Tabs) == 0 {
		index = 0
	}
	return &tab{
		index:     index,
		title:     name,
		primitive: newTabPrimitive(command),
	}
}

// Returns the primitive for a new tab
func newTabPrimitive(command string) (content tview.Primitive) {
	cmd := exec.Command(command)
	term := tterm.NewTerminal(cmd)
	term.SetBorder(true)
	term.SetTitle(fmt.Sprintf(" Welcome to %s ", command))
	return term
}
