package pdcli

import (
	"sort"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
)

type OncallUser struct {
	EscalationPolicy string
	OncallRole       string
	Name             string
	Start            string
	End              string
}

//TeamSREOnCall fetches the current roles and names of on-call users.
func TeamSREOnCall(c client.PagerDutyClient) ([]OncallUser, error) {
	var callOpts pagerduty.ListOnCallOptions
	var oncallData []OncallUser

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
		return nil, err
	}

	// OnCalls array contains all information about the API object
	for _, y := range oncallListing.OnCalls {

		timeConversionStart, err := utils.FormatTimestamp(y.Start)

		if err != nil {
			return nil, err
		}

		timeConversionEnd, err := utils.FormatTimestamp(y.End)

		if err != nil {
			return nil, err
		}

		// Parse the oncall data to OncallUser object
		temp := OncallUser{}
		temp.EscalationPolicy = y.EscalationPolicy.Summary
		temp.OncallRole = y.Schedule.Summary
		temp.Name = y.User.Summary
		temp.Start = timeConversionStart
		temp.End = timeConversionEnd
		oncallData = append(oncallData, temp)
	}

	return oncallData, err
}

// AllTeamsOncall displays the oncall data of all Red Hat PagerDuty teams.
func AllTeamsOncall(c client.PagerDutyClient) ([]OncallUser, error) {
	var callOpts pagerduty.ListOnCallOptions
	var oncallData []OncallUser

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
			return nil, err
		}

		// Parse oncall data
		for _, y := range oncallListing.OnCalls {
			temp := OncallUser{}
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

	return oncallData, nil
}

// UserNextOncallSchedule displays the current user's
// next oncall schedule.
func UserNextOncallSchedule(c client.PagerDutyClient) ([]OncallUser, error) {
	var callOpts pagerduty.ListOnCallOptions
	var nextOncallData []OncallUser

	callOpts.Until = time.Now().AddDate(0, 3, 0).String()

	userID, err := GetCurrentUserID(c)

	if err != nil {
		return nil, err
	}

	callOpts.UserIDs = append(callOpts.UserIDs, userID)

	// Fetch the oncall data from pagerduty API
	onCallOncallUser, err := c.ListOnCalls(callOpts)

	if err != nil {
		return nil, err
	}

	for _, y := range onCallOncallUser.OnCalls {

		start, err := utils.FormatTimestamp(y.Start)

		if err != nil {
			return nil, err
		}

		end, err := utils.FormatTimestamp(y.End)

		if err != nil {
			return nil, err
		}

		temp := OncallUser{}
		temp.EscalationPolicy = y.EscalationPolicy.Summary
		temp.OncallRole = y.Schedule.Summary
		temp.Name = y.User.Summary
		temp.Start = start
		temp.End = end
		nextOncallData = append(nextOncallData, temp)
	}

	return nextOncallData, nil
}
