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
	FooterText                = "[Esc] Go Back"
	FooterTextAlerts          = "[R] Refresh Alerts | [1] Trigerred Alerts | [2] Acknowledged Incidents | [3] Trigerred Incidents\n" + FooterText
	FooterTextTrigerredAlerts = "[1] Trigerred Alerts | [2] Acknowledged Incidents | [3] Trigerred Incidents\n" + FooterText
	FooterTextIncidents       = "[ENTER] Select Incident | [CTRL+A] Acknowledge Incidents\n" + FooterText
	FooterTextOncall          = "[N] Your Next Oncall Schedule | [A] All Teams Oncall\n" + FooterText
	TerminalFooterText        = "[CTRL + N] Next Slide | [CTRL + P] Previous Slide | [CTRL + A] Add Slide | [CTRL + E] Exit Slide | [CTRL + B] + [Num] Change to Slide with [Num]  | [CTRL + Q] Quit "
	TerminalFooterEscapeState = "Enter the Slide Number to Switch To : "

	// Colors
	TableTitleColor                = tcell.ColorLightCyan
	BorderColor                    = tcell.ColorLightGray
	FooterTextColor                = tcell.ColorGray
	InfoTextColor                  = tcell.ColorLightSlateGray
	ErrorTextColor                 = tcell.ColorRed
	PromptTextColor                = tcell.ColorLightGreen
	LoggerTextColor                = tcell.ColorGreen
	TerminalFooterTextColor        = tcell.ColorGreen
	TerminalFooterEscapeStateColor = tcell.ColorDarkGreen
)
