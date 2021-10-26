package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/rivo/tview"
)

type TUI struct {

	// Main UI elements
	App            *tview.Application
	AlertMetadata  *tview.TextView
	Table          *tview.Table
	IncidentsTable *tview.Table
	Pages          *tview.Pages
	Info           *tview.TextView
	Layout         *tview.Flex
	Footer         *tview.TextView

	// Misc. UI elements
	isMouseEnabled bool
	secondaryText  string

	// Internals
	Username      string
	AssginedTo    string
	FetchedAlerts string
	Incidents     [][]string
	AckIncidents  []string
	ClusterID     string
	Client        client.PagerDutyClient
}

const (
	// Table Titles
	AlertsTableTitle       = "[ Alerts ]"
	AlertMetadataViewTitle = "[ Alert Data ]"

	// Page Titles
	AlertsPageTitle       = "Alerts"
	AlertDataPageTitle    = "Metadata"
	AckIncidentsPageTitle = "Incidents"

	// Text Format
	TitleFmt = " [lightcyan::b]%s "

	// Footer
	FooterText       = "[Q] Quit | [Esc] Go Back | [M] Enable/Disable Mouse"
	FooterTextAlerts = FooterText + " | [A] Ack Mode"
	FooterTextAck    = FooterText + " | [ENTER] Select Incident  | [CTRL+A] Acknowledge Incidents"

	// Colors
	TableTitleColor = tcell.ColorLightCyan
	BorderColor     = tcell.ColorLightGray
	FooterTextColor = tcell.ColorGray
	InfoTextColor   = tcell.ColorLightSlateGray
	ErrorTextColor  = tcell.ColorRed
	PromptTextColor = tcell.ColorLightGreen
)

// queueUpdateDraw is used to synchronize access to primitives from non-main goroutines.
func (tui *TUI) queueUpdateDraw(f func()) {
	go func() {
		tui.App.QueueUpdateDraw(f)
	}()
}

// showDefaultSecondaryView sets the secondary textview component on every page for the given secondaryText.
func (tui *TUI) showDefaultSecondaryView() {
	tui.queueUpdateDraw(func() {
		tui.Info.SetText(tui.secondaryText).SetTextColor(InfoTextColor)
	})
}

// showSecondaryView sets the secondary textview component to render a specific string.
func (tui *TUI) showSecondaryView(msg string) {
	tui.queueUpdateDraw(func() {
		tui.Info.SetText(msg).SetTextColor(TableTitleColor)
	})

	go time.AfterFunc(4*time.Second, tui.showDefaultSecondaryView)
}

// promptSecondaryView sets the secondary textview component to a specific string with no timeout view refresh.
func (tui *TUI) promptSecondaryView(msg string) {
	tui.queueUpdateDraw(func() {
		tui.Info.SetText(msg).SetTextColor(PromptTextColor)
	})
}

// showError renders the program error string to the secondary textview component.
func (tui *TUI) showError(msg string) {
	tui.queueUpdateDraw(func() {
		tui.Info.SetText(msg).SetTextColor(ErrorTextColor)
	})

	go time.AfterFunc(5*time.Second, tui.showDefaultSecondaryView)
}

// Init initializes all the TUI main elements.
func (tui *TUI) Init() {

	tui.App = tview.NewApplication()
	tui.Pages = tview.NewPages()
	tui.Info = tview.NewTextView()
	tui.Footer = tview.NewTextView()
	tui.AlertMetadata = tview.NewTextView()

	tui.Info.
		SetBorder(true).
		SetBorderColor(BorderColor).
		SetBorderAttributes(tcell.AttrDim).
		SetBorderPadding(0, 0, 1, 1)

	tui.Footer.
		SetTextAlign(tview.AlignLeft).
		SetTextColor(FooterTextColor).
		SetBorderPadding(2, 2, 1, 1)

	tui.AlertMetadata.
		SetScrollable(true).
		SetBorder(true).
		SetBorderColor(BorderColor).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderAttributes(tcell.AttrDim).
		SetTitle(fmt.Sprintf(TitleFmt, AlertMetadataViewTitle))

	tui.showDefaultSecondaryView()

	// Create the main layout.
	tui.Layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.Pages, 0, 6, true).
		AddItem(tui.Info, 0, 1, false).
		AddItem(tui.Footer, 0, 1, false)
}

// initFooter initializes the footer text depending on the page currently visible.
func (t *TUI) initFooter() {
	name, _ := t.Pages.GetFrontPage()

	switch name {

	case AlertsPageTitle:
		t.Footer.SetText(FooterTextAlerts)

	default:
		t.Footer.SetText(FooterText)

	}
}

// StartApp sets the UI layout and renders all the TUI elements.
func (t *TUI) StartApp() error {
	t.initFooter()
	t.initKeyboard()

	return t.App.SetRoot(t.Layout, true).EnableMouse(false).Run()
}
