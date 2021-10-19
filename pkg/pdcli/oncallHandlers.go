package pdcli

import (
	"sort"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/output"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

type User struct {
	EscalationPolicy string
	OncallRole       string
	Name             string
	Start            string
	End              string
}

//TeamSREOnCall fetches the current roles and names of on-call users.
func TeamSREOnCall(c client.PagerDutyClient) error {
	var callOpts pagerduty.ListOnCallOptions
	var oncallData []User

	callOpts.ScheduleIDs = []string{
		constants.PrimaryScheduleID,
		constants.SecondaryScheduleID,
		constants.OncallManager,
		constants.OncallIDWeekend,
		constants.InvestigatorID,
	}

	// Fetch the oncall data from pagerduty API
	oncallListing, err := c.ListOnCalls(callOpts)

	if err != nil {
		return err
	}

	// OnCalls array contains all information about the API object
	for _, y := range oncallListing.OnCalls {

		timeConversionStart, err := utils.FormatTimestamp(y.Start)

		if err != nil {
			return err
		}

		timeConversionEnd, err := utils.FormatTimestamp(y.End)

		if err != nil {
			return err
		}

		// Parse the oncall data to user object
		temp := User{}
		temp.EscalationPolicy = y.EscalationPolicy.Summary
		temp.OncallRole = y.Schedule.Summary
		temp.Name = y.User.Summary
		temp.Start = timeConversionStart
		temp.End = timeConversionEnd
		oncallData = append(oncallData, temp)
	}

	// Print the oncall users to console
	printOncalls(oncallData)

	return nil
}

// AllTeamsOncall displays the oncall data of all Red Hat PagerDuty teams.
func AllTeamsOncall(c client.PagerDutyClient) error {
	var callOpts pagerduty.ListOnCallOptions
	var oncallData []User

	offset := []uint{0, 100, 200, 300, 400, 500, 600}

	callOpts.Limit = 100
	callOpts.Offset = 0

	for _, o := range offset {
		callOpts.Limit = 100
		callOpts.Offset = o
		callOpts.Earliest = true

		// Fetch the oncall data from pagerduty API
		oncallListing, err := c.ListOnCalls(callOpts)

		if err != nil {
			return err
		}

		// Parse oncall data
		for _, y := range oncallListing.OnCalls {
			temp := User{}
			temp.EscalationPolicy = y.EscalationPolicy.Summary
			temp.OncallRole = y.Schedule.Summary
			temp.Name = y.User.Summary
			temp.Start = y.Start
			temp.End = y.End
			oncallData = append(oncallData, temp)
		}
	}

	// Sort by escalation policy
	sort.SliceStable(oncallData, func(i, j int) bool {
		return oncallData[i].EscalationPolicy < oncallData[j].EscalationPolicy
	})

	// Print the oncall users to console
	printOncalls(oncallData)

	return nil
}

// UserNextOncallSchedule displays the current user's
// next oncall schedule.
func UserNextOncallSchedule(c client.PagerDutyClient) error {
	var callOpts pagerduty.ListOnCallOptions

	callOpts.Until = time.Now().AddDate(0, 3, 0).String()

	userID, err := GetCurrentUserID(c)

	if err != nil {
		return err
	}

	callOpts.UserIDs = append(callOpts.UserIDs, userID)

	// Fetch the oncall data from pagerduty API
	onCallUser, err := c.ListOnCalls(callOpts)

	if err != nil {
		return err
	}
	// Initialize table writer
	table := output.NewTable(true)

	for _, y := range onCallUser.OnCalls {
		var data []string

		start, err := utils.FormatTimestamp(y.Start)

		if err != nil {
			return err
		}

		end, err := utils.FormatTimestamp(y.End)

		if err != nil {
			return err
		}

		data = append(data, y.Schedule.Summary, start, end)

		table.AddRow(data)
	}

	headers := []string{"Oncall Role", "From", "To"}

	table.SetHeaders(headers)
	table.SetData()
	table.Print()

	return nil
}

//printOncalls prints data in a tabular form.
func printOncalls(oncallData []User) {

	// Initialize table writer
	table := output.NewTable(false)

	for _, v := range oncallData {

		var data []string

		if v.EscalationPolicy != "" {
			data = append(data, v.EscalationPolicy)
		} else {
			data = append(data, "N/A")
		}

		if v.Name != "" {
			data = append(data, v.Name)
		} else {
			data = append(data, "N/A")
		}

		if v.OncallRole != "" {
			data = append(data, v.OncallRole)
		} else {
			data = append(data, "N/A")
		}

		if v.Start != "" {
			data = append(data, v.Start)
		} else {
			data = append(data, "N/A")
		}

		if v.End != "" {
			data = append(data, v.End)
		} else {
			data = append(data, "N/A")
		}

		table.AddRow(data)
	}

	headers := []string{"Escalation Policy", "Oncall Role", "Name", "From", "To"}

	table.SetHeaders(headers)
	table.SetData()
	table.Print()
}
