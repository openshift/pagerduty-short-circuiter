package ui

import "github.com/gdamore/tcell/v2"

const (

	// Text Format
	TitleFmt = " [lightcyan::b]%s "

	// Table Titles
	AlertsTableTitle          = "[ ALERTS ]"
	ResolvedAlertsTableTitle  = "[ RESOLVED ALERTS ]"
	TrigerredAlertsTableTitle = "[ TRIGERRED ALERTS ]"
	HighAlertsTableTitle      = "[ TRIGERRED ALERTS - HIGH ]"
	LowAlertsTableTitle       = "[ TRIGERRED ALERTS - LOW ]"
	AlertMetadataViewTitle    = "[ ALERT DATA ]"
	IncidentsTableTitle       = "[ TRIGERRED INCIDENTS ]"
	AckIncidentsTableTitle    = "[ ACKNOWLEDGED INCIDENTS ]"
	OncallTableTitle          = "[ ONCALL ]"
	NextOncallTableTitle      = "[ NEXT ONCALL ]"
	AllTeamsOncallTableTitle  = "[ ALL TEAMS ONCALL ]"

	// Page Titles
	AlertsPageTitle          = "Alerts"
	AlertDataPageTitle       = "Metadata"
	ResolvedAlertsPageTitle  = "Resolved"
	TrigerredAlertsPageTitle = "Trigerred"
	HighAlertsPageTitle      = "High Alerts"
	LowAlertsPageTitle       = "Low Alerts"
	IncidentsPageTitle       = "Incidents"
	AckIncidentsPageTitle    = "AckIncidents"
	OncallPageTitle          = "Oncall"
	NextOncallPageTitle      = "Next Oncall"
	AllTeamsOncallPageTitle  = "All Teams Oncall"

	// Footer
	FooterText                = "[Q] Quit | [Esc] Go Back"
	FooterTextStatus          = "[H] High Alerts | [L] Low Alerts\n"
	FooterTextAlerts          = "[1] Resolved Alerts | [2] Trigerred Alerts | [3] Acknowledged Incidents | [4] Trigerred Incidents\n" + FooterText
	FooterTextTrigerredAlerts = "[1] Resolved Alerts | [2] Trigerred Alerts | [3] Acknowledged Incidents | [4] Trigerred Incidents\n" + FooterTextStatus + FooterText
	FooterTextIncidents       = "[ENTER] Select Incident  | [CTRL+A] Acknowledge Incidents\n" + FooterText
	FooterTextOncall          = "[N] Your Next Oncall Schedule | [A] All Teams Oncall\n" + FooterText

	// Colors
	TableTitleColor = tcell.ColorLightCyan
	BorderColor     = tcell.ColorLightGray
	FooterTextColor = tcell.ColorGray
	InfoTextColor   = tcell.ColorLightSlateGray
	ErrorTextColor  = tcell.ColorRed
	PromptTextColor = tcell.ColorLightGreen
)
