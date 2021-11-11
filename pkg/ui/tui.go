package ui

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/gdamore/tcell/v2"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/alerts"
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
	Info                *tview.TextView
	Layout              *tview.Flex
	Footer              *tview.TextView
	FrontPage           string

	// Misc. UI elements
	secondaryText string

	// API related
	Client       client.PagerDutyClient
	IncidentOpts pagerduty.ListIncidentsOptions
	Alerts       []pdcli.Alert

	// Internals
	SelectedIncidents map[string]string
	Incidents         [][]string
	AckIncidents      []string
	Username          string
	AssginedTo        string
	Columns           string
	FetchedAlerts     string
	ClusterID         string
	Primary           string
	Secondary         string
	HasEmulator       bool
}

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

// InitAlertsUI initializes TUI table component.
// It adds the returned table as a new TUI page view.
func (tui *TUI) InitAlertsUI(alerts []pdcli.Alert, tableTitle string, pageTitle string) {
	headers, data := pdcli.GetTableData(alerts, tui.Columns)
	tui.Table = tui.InitTable(headers, data, true, false, tableTitle)
	tui.SetAlertsTableEvents(alerts)

	tui.SetAlertsSecondaryData()

	if len(alerts) == 0 {
		tui.showSecondaryView("No alerts to display")
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
	incidentHeaders := []string{"INCIDENT ID", "NAME", "SEVERITY", "STATUS", "SERVICE"}

	if isSelectable {
		tui.IncidentsTable = tui.InitTable(incidentHeaders, incidents, true, true, tableTitle)
		tui.SetIncidentsTableEvents()
	} else {
		tui.IncidentsTable = tui.InitTable(incidentHeaders, incidents, false, false, tableTitle)
	}

	tui.Pages.AddPage(pageTitle, tui.IncidentsTable, true, false)
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
		SetBorderPadding(1, 0, 1, 1)

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

// StartApp sets the UI layout and renders all the TUI elements.
func (t *TUI) StartApp() error {
	t.initFooter()
	t.initKeyboard()

	return t.App.SetRoot(t.Layout, true).EnableMouse(false).Run()
}
