package ui

import "github.com/gdamore/tcell/v2"

const (

	// Text Format
	TitleFmt = " [lightcyan::b]%s "

	// Table Titles
	AlertsTableTitle          = "[ ALERTS ]"
	ResolvedAlertsTableTitle  = "[ RESOLVED ALERTS ]"
	TrigerredAlertsTableTitle = "[ TRIGERRED ALERTS ]"
	AlertMetadataViewTitle    = "[ ALERT DATA ]"
	IncidentsTableTitle       = "[ INCIDENTS ]"
	AckIncidentsTableTitle    = "[ ACKNOWLEDGED INCIDENTS ]"
	OncallTableTitle          = "[ ONCALL ]"
	NextOncallTableTitle      = "[ NEXT ONCALL ]"
	AllTeamsOncallTableTitle  = "[ ALL TEAMS ONCALL ]"

	// Page Titles
	AlertsPageTitle          = "Alerts"
	AlertDataPageTitle       = "Metadata"
	ResolvedAlertsPageTitle  = "Resolved"
	TrigerredAlertsPageTitle = "Trigerred"
	IncidentsPageTitle       = "Incidents"
	AckIncidentsPageTitle    = "AckIncidents"
	OncallPageTitle          = "Oncall"
	NextOncallPageTitle      = "Next Oncall"
	AllTeamsOncallPageTitle  = "All Teams Oncall"

	// Footer
	FooterText          = "[Q] Quit | [Esc] Go Back"
	FooterTextAlerts    = FooterText + " | [A] Ack Mode | [1] View Resolved Alerts | [2] View Trigerred Alerts"
	FooterTextIncidents = FooterText + " | [ENTER] Select Incident  | [CTRL+A] Acknowledge Incidents | [1] View Acknowledged Incidents"
	FooterTextOncall    = FooterText + " | [N] View Your Next Oncall Schedule | [A] View All Teams Oncall"

	// Colors
	TableTitleColor = tcell.ColorLightCyan
	BorderColor     = tcell.ColorLightGray
	FooterTextColor = tcell.ColorGray
	InfoTextColor   = tcell.ColorLightSlateGray
	ErrorTextColor  = tcell.ColorRed
	PromptTextColor = tcell.ColorLightGreen
)
