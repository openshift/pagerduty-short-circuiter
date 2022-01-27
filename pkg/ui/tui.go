package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/rivo/tview"
)

type TUI struct {
	App               *tview.Application
	Table             *tview.Table
	AckIncidentsTable *tview.Table
	IncidentsTable    *tview.Table
	AlertMetadata     *tview.TextView
	SecondaryWindow   *tview.TextView
	LogWindow         *tview.TextView
	Footer            *tview.TextView
	Pages             *tview.Pages
	Layout            *tview.Flex
	FrontPage         string
}

// initFooter initializes the footer text depending on the page currently visible.
func (t *TUI) initFooter() {
	name, _ := t.Pages.GetFrontPage()

	switch name {
	case AlertsPageTitle:
		t.Footer.SetText(FooterTextAlerts)

	case OncallPageTitle:
		t.Footer.SetText(FooterTextOncall)

	default:
		t.Footer.SetText(FooterText)
	}
}

// Init initializes all the TUI main elements.
func (tui *TUI) Init() {
	tui.App = tview.NewApplication()
	tui.Pages = tview.NewPages()
	tui.SecondaryWindow = tview.NewTextView()
	tui.LogWindow = tview.NewTextView()
	tui.Footer = tview.NewTextView()
	tui.AlertMetadata = tview.NewTextView()

	tui.SecondaryWindow.
		SetChangedFunc(func() { tui.App.Draw() }).
		SetTextColor(InfoTextColor).
		SetScrollable(true).
		ScrollToEnd().
		SetBorder(true).
		SetBorderColor(BorderColor).
		SetBorderAttributes(tcell.AttrDim).
		SetBorderPadding(1, 1, 1, 1)

	tui.LogWindow.
		SetChangedFunc(func() { tui.App.Draw() }).
		SetScrollable(true).
		ScrollToEnd().
		SetBorder(true).
		SetBorderColor(BorderColor).
		SetBorderAttributes(tcell.AttrDim).
		SetBorderPadding(0, 0, 1, 1)

	tui.Footer.
		SetTextAlign(tview.AlignLeft).
		SetTextColor(FooterTextColor).
		SetBorderPadding(1, 0, 1, 1)

	tui.AlertMetadata.
		SetScrollable(true).
		SetBorder(true).
		SetBorderColor(BorderColor).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderAttributes(tcell.AttrDim).
		SetTitle(fmt.Sprintf(TitleFmt, AlertMetadataViewTitle))

	// Initialize logger to output to log view
	utils.InitLogger(tui.LogWindow)

	// Create the main layout
	tui.Layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.Pages, 0, 6, true).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(tui.SecondaryWindow, 0, 1, false).
				AddItem(tui.LogWindow, 0, 2, false),
			0, 2, false).
		AddItem(tui.Footer, 0, 1, false)
}

// StartApp sets the UI layout and renders all the TUI elements.
func (t *TUI) StartApp() error {
	t.initFooter()
	return t.App.SetRoot(t.Layout, true).EnableMouse(true).Run()
}
