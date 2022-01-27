package ui

import "github.com/gdamore/tcell/v2"

const (

	// Text Format
	TitleFmt = " [lightcyan::b]%s "

	// Table Titles
	AlertsTableTitle         = "[ ALERTS ]"
	HighAlertsTableTitle     = "[ HIGH ALERTS ]"
	LowAlertsTableTitle      = "[ LOW ALERTS ]"
	AlertMetadataViewTitle   = "[ ALERT DATA ]"
	IncidentsTableTitle      = "[ TRIGERRED INCIDENTS ]"
	AckIncidentsTableTitle   = "[ ACKNOWLEDGED INCIDENTS ]"
	OncallTableTitle         = "[ ONCALL ]"
	NextOncallTableTitle     = "[ NEXT ONCALL ]"
	AllTeamsOncallTableTitle = "[ ALL TEAMS ONCALL ]"

	// Page Titles
	AlertsPageTitle         = "Alerts"
	AlertDataPageTitle      = "Metadata"
	HighAlertsPageTitle     = "High Alerts"
	LowAlertsPageTitle      = "Low Alerts"
	IncidentsPageTitle      = "Incidents"
	AckIncidentsPageTitle   = "AckIncidents"
	OncallPageTitle         = "Oncall"
	NextOncallPageTitle     = "Next Oncall"
	AllTeamsOncallPageTitle = "All Teams Oncall"

	// Footer
	FooterText            = "[Q] Quit | [Esc] Go Back"
	FooterTextAlertStatus = "[H] High Alerts | [L] Low Alerts"
	FooterTextAlerts      = FooterTextAlertStatus + "\n[1] Acknowledged Incidents | [2] Trigerred Incidents | [R] Refresh Alerts\n" + FooterText
	FooterTextIncidents   = "[ENTER] Select Incident  | [CTRL+A] Acknowledge Incidents\n" + FooterText
	FooterTextOncall      = "[N] Your Next Oncall Schedule | [A] All Teams Oncall\n" + FooterText

	// Colors
	TableTitleColor = tcell.ColorLightCyan
	BorderColor     = tcell.ColorLightGray
	FooterTextColor = tcell.ColorGray
	InfoTextColor   = tcell.ColorLightSlateGray
	ErrorTextColor  = tcell.ColorRed
	PromptTextColor = tcell.ColorLightGreen
	LoggerTextColor = tcell.ColorGreen
)
