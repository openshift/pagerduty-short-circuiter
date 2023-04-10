package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"git.sr.ht/~rockorager/tterm"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/rivo/tview"
)

// Declares the tab struct
type TerminalTab struct {
	index     int
	title     string
	primitive tview.Primitive
}

var CurrentActivePage int = 0
var TotalPageCount int = -1

// Creates and return a new tab
func InitKiteTab(tui *TUI, layout *tview.Flex) *TerminalTab {
	TotalPageCount += 1
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs, TotalPageCount)
	index := len(tui.TerminalTabs)
	if len(tui.TerminalTabs) == 0 {
		index = 0
	}
	return &TerminalTab{
		index:     index,
		title:     "kite",
		primitive: layout,
	}
}

// Creates and return a new tab
func NewTab(name string, command string, tui *TUI) *TerminalTab {
	TotalPageCount += 1
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs, TotalPageCount)
	index := len(tui.TerminalTabs)
	if len(tui.TerminalTabs) == 0 {
		index = 0
	}
	return &TerminalTab{
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

// Move to the previous slide
func PreviousSlide(tui *TUI) {
	CurrentActivePage = (CurrentActivePage - 1 + len(tui.TerminalTabs)) % len(tui.TerminalTabs)
	tui.TerminalPageBar.Highlight(strconv.Itoa(tui.TerminalUIRegionIDs[CurrentActivePage])).
		ScrollToHighlight()
	tui.TerminalInputBuffer = []rune{}
}

// Move to the next slide
func NextSlide(tui *TUI) {
	CurrentActivePage = (CurrentActivePage + 1) % len(tui.TerminalTabs)
	tui.TerminalPageBar.Highlight(strconv.Itoa(tui.TerminalUIRegionIDs[CurrentActivePage])).
		ScrollToHighlight()
	tui.TerminalInputBuffer = []rune{}
}

func indexOf(arr []int, ele int) int {
	for index, item := range arr {
		if item == ele {
			return index
		}
	}
	return -1
}

// Remove the slide with the given index
func RemoveSlide(s int, tui *TUI) {
	index := indexOf(tui.TerminalUIRegionIDs, s)
	tui.TerminalTabs = append(tui.TerminalTabs[:index], tui.TerminalTabs[index+1:]...)
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs[:index], tui.TerminalUIRegionIDs[index+1:]...)
	tui.TerminalPageBar.Clear()
	for index, tabSlide := range tui.TerminalTabs {
		oldIndex := tabSlide.index
		tabSlide.index = index
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, oldIndex, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
	}
	tui.TerminalPages.RemovePage(strconv.Itoa(s))
	PreviousSlide(tui)
}

// Adds a slide to the end of currently present slides
func AddSlide(tui *TUI, name string) {
	var tabSlide TerminalTab
	if name == constants.Bash {
		tabSlide = *NewTab(name, os.Getenv("SHELL"), tui)
	} else if name == constants.OcmContainer {
		tabSlide = *NewTab(name, "ocm-container", tui)
	}
	tui.TerminalTabs = append(tui.TerminalTabs, tabSlide)
	tui.TerminalPages.AddPage(strconv.Itoa(tabSlide.index), tabSlide.primitive, true, tabSlide.index == 0)
	fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, tabSlide.index, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
	CurrentActivePage = tabSlide.index
	tui.TerminalPageBar.Highlight(strconv.Itoa(CurrentActivePage)).
		ScrollToHighlight()
	tui.TerminalInputBuffer = []rune{}
}

func InitTerminalMux(tui *TUI, kiteTab *TerminalTab) *tview.Flex {
	// Initial Slides
	tui.TerminalTabs = append(tui.TerminalTabs, *kiteTab)

	// Set the bottom navigation bar
	tui.TerminalPageBar.
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			tui.TerminalPages.SwitchToPage(added[0])
		})

	for _, slide := range tui.TerminalTabs {
		tui.TerminalPages.AddPage(strconv.Itoa(slide.index), slide.primitive, true, slide.index == 0)
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, slide.index, fmt.Sprintf("%d %s", slide.index+1, slide.title))
	}
	tui.TerminalPageBar.Highlight("0")

	// Returns the main view & layout for the app
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.TerminalPages, 0, 1, true).
		AddItem(tui.TerminalPageBar, 1, 1, false)
}
