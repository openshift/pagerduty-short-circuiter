package ui

import "github.com/gdamore/tcell/v2"

const (

	// Text Format
	TitleFmt = " [lightcyan::b]%s "

	// Table Titles
	AlertsTableTitle          = "[ ALERTS ]"
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
	FooterTextAlerts          = "[R] Refresh Alerts | [1] Trigerred Alerts | [2] Acknowledged Incidents | [3] Trigerred Incidents\n" + FooterText
	FooterTextTrigerredAlerts = "[1] Trigerred Alerts | [2] Acknowledged Incidents | [3] Trigerred Incidents\n" + FooterTextStatus + FooterText
	FooterTextIncidents       = "[ENTER] Select Incident  | [CTRL+A] Acknowledge Incidents\n" + FooterText
	FooterTextOncall          = "[N] Your Next Oncall Schedule | [A] All Teams Oncall\n" + FooterText

	// Colors
	TableTitleColor = tcell.ColorLightCyan
	BorderColor     = tcell.ColorLightGray
	FooterTextColor = tcell.ColorGray
	InfoTextColor   = tcell.ColorLightSlateGray
	ErrorTextColor  = tcell.ColorRed
	PromptTextColor = tcell.ColorLightGreen
	LoggerTextColor = tcell.ColorGreen
)
