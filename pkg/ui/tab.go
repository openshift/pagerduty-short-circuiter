package ui

import (
	"fmt"
	"os/exec"
	"strconv"

	"git.sr.ht/~rockorager/tterm"
	"github.com/rivo/tview"
)

// Declares the tab struct
type TerminalTab struct {
	index     int
	regionID  int
	title     string
	primitive tview.Primitive
}

var CurrentActivePage int = 0
var TotalPageCount int = 0
var CursorPos int

// Creates and return a new kite tab
func InitKiteTab(tui *TUI, layout *tview.Flex) *TerminalTab {
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs, TotalPageCount)
	return &TerminalTab{
		index:     0,
		regionID:  0,
		title:     "kite",
		primitive: layout,
	}
}

// Creates and returns a SOP tab
func InitSOPTab(name string, layout *tview.TextView, tui *TUI) *TerminalTab {
	TotalPageCount += 1
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs, TotalPageCount)
	index := len(tui.TerminalTabs)
	return &TerminalTab{
		index:     index,
		regionID:  TotalPageCount,
		title:     name,
		primitive: layout,
	}
}

// Creates and return a new tab
func NewTab(name string, command string, args []string, tui *TUI) *TerminalTab {
	TotalPageCount += 1
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs, TotalPageCount)
	index := len(tui.TerminalTabs)
	return &TerminalTab{
		index:     index,
		regionID:  TotalPageCount,
		title:     name,
		primitive: newTabPrimitive(command, args),
	}
}

// Returns the primitive for a new tab
func newTabPrimitive(command string, args []string) (content tview.Primitive) {
	cmd := exec.Command(command, args...)
	term := tterm.NewTerminal(cmd)
	term.SetBorder(true)
	term.SetTitle(fmt.Sprintf(" Welcome to %s ", command))
	return term
}

// Move to the previous slide
func PreviousSlide(tui *TUI) {
	CurrentActivePage = (CurrentActivePage - 1 + len(tui.TerminalTabs)) % len(tui.TerminalTabs)
	tui.TerminalPageBar.Highlight(strconv.Itoa(tui.TerminalTabs[CurrentActivePage].regionID)).
		ScrollToHighlight()
	tui.TerminalInputBuffer = []rune{}
}

// Move to the next slide
func NextSlide(tui *TUI) {
	CurrentActivePage = (CurrentActivePage + 1 + len(tui.TerminalTabs)) % len(tui.TerminalTabs)
	tui.TerminalPageBar.Highlight(strconv.Itoa(tui.TerminalTabs[CurrentActivePage].regionID)).
		ScrollToHighlight()
	tui.TerminalInputBuffer = []rune{}
}

func indexOf(ele int, tabs []TerminalTab) int {
	for index, item := range tabs {
		if item.regionID == ele {
			return index
		}
	}
	return -1
}

// Remove the slide with the given index
// Exit the app if only one slide is present
func RemoveSlide(s int, tui *TUI) {
	if len(tui.TerminalTabs) == 1 {
		tui.App.Stop()
		return
	}
	index := indexOf(s, tui.TerminalTabs)
	tui.TerminalTabs = append(tui.TerminalTabs[:index], tui.TerminalTabs[index+1:]...)
	tui.TerminalUIRegionIDs = append(tui.TerminalUIRegionIDs[:index], tui.TerminalUIRegionIDs[index+1:]...)
	tui.TerminalPageBar.Clear()
	for index, tabSlide := range tui.TerminalTabs {
		tabSlide.index = index
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, tabSlide.regionID, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
	}
	tui.TerminalPages.RemovePage(strconv.Itoa(s))
	PreviousSlide(tui)
	tui.TerminalInputBuffer = []rune{}
}

// Adds a slide to the end of currently present slides
func AddNewSlide(tui *TUI, name string, command string, args []string, isCluster bool) {
	if len(tui.TerminalTabs) < 9 {
		if isCluster {
			for i, tab := range tui.TerminalTabs {
				if tab.primitive != nil && tab.title == args[0] {
					tui.TerminalPageBar.Highlight(strconv.Itoa(i)).
						ScrollToHighlight()
					return
				}
			}
		}
		tabSlide := NewTab(name, command, args, tui)
		tui.TerminalTabs = append(tui.TerminalTabs, *tabSlide)
		tui.TerminalPages.AddPage(strconv.Itoa(tabSlide.regionID), tabSlide.primitive, true, true)
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, tabSlide.regionID, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
		CurrentActivePage = tabSlide.index
		tui.TerminalPageBar.Highlight(strconv.Itoa(tabSlide.regionID)).
			ScrollToHighlight()
	}
	tui.TerminalInputBuffer = []rune{}
}

// Adds a SOP slide to the end of currently present slides
func AddSOPSlide(name string, textView *tview.TextView, tui *TUI) {
	if len(tui.TerminalTabs) < 9 {
		for i, tab := range tui.TerminalTabs {
			if tab.primitive != nil && tab.title == name {
				CurrentActivePage = tab.index
				tui.TerminalPageBar.Highlight(strconv.Itoa(i)).
					ScrollToHighlight()
				return
			}
		}
		tabSlide := InitSOPTab(name, textView, tui)
		tui.TerminalTabs = append(tui.TerminalTabs, *tabSlide)
		tui.TerminalPages.AddPage(strconv.Itoa(tabSlide.regionID), tabSlide.primitive, true, true)
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, tabSlide.regionID, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
		CurrentActivePage = tabSlide.index
		tui.TerminalPageBar.Highlight(strconv.Itoa(tabSlide.regionID)).
			ScrollToHighlight()
	}
	tui.TerminalInputBuffer = []rune{}
}

// Navigate to the specified slide
func SwitchToSlide(slideNum int, tui *TUI) {
	if slideNum > 0 && slideNum <= len(tui.TerminalTabs) {
		slideNum = slideNum - 1
		regionID := tui.TerminalTabs[slideNum].regionID
		tui.TerminalPageBar.Highlight(strconv.Itoa(regionID)).
			ScrollToHighlight()
		CurrentActivePage = slideNum
	}
}

// Init the Layout for Terminal Multiplexer
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
		tui.TerminalPages.AddPage(strconv.Itoa(slide.index), slide.primitive, true, true)
		fmt.Fprintf(tui.TerminalPageBar, `["%d"]%s[white][""]  `, slide.index, fmt.Sprintf("%d %s", slide.index+1, slide.title))
	}
	tui.TerminalPageBar.Highlight("0")
	tui.TerminalFixedFooter.SetText(TerminalFooterText)

	// Returns the main view & layout for the app
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.TerminalPages, 0, 1, true).
		AddItem(tui.TerminalPageBar, 1, 1, false).
		AddItem(tui.TerminalFixedFooter, 1, 1, false)
}
