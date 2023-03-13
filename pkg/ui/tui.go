package ui

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/rivo/tview"
)

type TUI struct {

	// Main UI elements
	App                 *tview.Application
	AlertMetadata       *tview.TextView
	Table               *tview.Table
	IncidentsTable      *tview.Table
	NextOncallTable     *tview.Table
	AllTeamsOncallTable *tview.Table
	Pages               *tview.Pages
	SecondaryWindow     *tview.TextView
	LogWindow           *tview.TextView
	Layout              *tview.Flex
	Footer              *tview.TextView
	FrontPage           string

	// API related
	Client       client.PagerDutyClient
	IncidentOpts pagerduty.ListIncidentsOptions
	Alerts       []pdcli.Alert

	// Internals
	SelectedIncidents map[string]string
	Incidents         [][]string
	AckIncidents      []string
	AssignedTo        string
	Username          string
	Role              string
	Columns           string
	ClusterID         string
}

// InitAlertsUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func (tui *TUI) InitAlertsUI(alerts []pdcli.Alert, tableTitle string, pageTitle string) {
	headers, data := pdcli.GetTableData(alerts, tui.Columns)
	tui.Table = tui.InitTable(headers, data, true, false, tableTitle)
	tui.SetAlertsTableEvents(alerts)

	if len(alerts) == 0 && tui.Username == tui.AssignedTo {
		utils.InfoLogger.Printf("No acknowledged alerts for user %s found", tui.Username)
	}

	tui.Pages.AddPage(pageTitle, tui.Table, true, true)
	tui.FrontPage = pageTitle

	if pageTitle == TrigerredAlertsPageTitle {
		tui.Footer.SetText(FooterTextTrigerredAlerts)
	} else {
		tui.Footer.SetText(FooterTextAlerts)
	}
}

// InitIncidentsUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func (tui *TUI) InitIncidentsUI(incidents [][]string, tableTitle string, pageTitle string, isSelectable bool) {
	incidentHeaders := []string{"INCIDENT ID", "NAME", "SEVERITY", "STATUS", "SERVICE", "ASSIGNED TO"}

	if isSelectable {
		tui.IncidentsTable = tui.InitTable(incidentHeaders, incidents, true, true, tableTitle)
		tui.SetIncidentsTableEvents()
	} else {
		tui.IncidentsTable = tui.InitTable(incidentHeaders, incidents, false, false, tableTitle)
	}

	if !tui.Pages.HasPage(pageTitle) {
		tui.Pages.AddPage(pageTitle, tui.IncidentsTable, true, false)
	}
}

func (tui *TUI) InitAlertsSecondaryView() {
	tui.SecondaryWindow.SetText(
		fmt.Sprintf("Logged in user: %s\n\nViewing alerts assigned to: %s\n\nPagerDuty role: %s",
			tui.Username,
			tui.AssignedTo,
			tui.Role)).
		SetTextColor(InfoTextColor)
}

func (tui *TUI) InitOnCallSecondaryView(user string, primary string, secondary string) {
	tui.SecondaryWindow.SetText(
		fmt.Sprintf("Logged in user: %s\n\nPrimary on-call: %s\n\nSecondary on-call: %s",
			user,
			primary,
			secondary),
	)
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
	t.initKeyboard()

	return t.App.SetRoot(t.Layout, true).EnableMouse(true).Run()
}
