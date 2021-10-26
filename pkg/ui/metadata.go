package ui

// SetAlertsSecondaryData initializes the text data for the secondary textview component of the UI.
func (tui *TUI) SetAlertsSecondaryData() {
	tui.secondaryText = "Logged in user: " + tui.Username + "\n" +
		"Viewing alerts assigned to: " + tui.AssginedTo + "\n" +
		"Number of alerts fetched: " + tui.FetchedAlerts
}
