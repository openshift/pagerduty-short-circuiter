package pdcli

import (
	"sort"
	"strings"
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

type OncallLayer struct {
	LayerId string
	Users   []OncallUser
}

// TeamSREOnCall fetches the current roles and names of on-call users.
func TeamSREOnCall(c client.PagerDutyClient) ([]OncallLayer, error) {
	var callOpts pagerduty.ListOnCallOptions
	var oncallLayers []OncallLayer

	callOpts.ScheduleIDs = []string{
		constants.PrimaryScheduleID,
		constants.SecondaryScheduleID,
		constants.OncallManager,
		constants.OncallIDWeekend,
		constants.InvestigatorID,
	}
	since := time.Now().Add(time.Hour * -11)
	until := time.Now().Add(time.Hour * 13)
	callOpts.Since = since.String()
	callOpts.Until = until.String()
	callOpts.Limit = 100
	// Fetch the oncall data from pagerduty API
	oncallListing, err := c.ListOnCalls(callOpts)

	if err != nil {
		return nil, err
	}

	startTime, _ := utils.FormatTimestamp(oncallListing.OnCalls[0].Start)
	var temp []OncallUser
	var mgmtUsers []OncallUser

	// OnCalls array contains all information about the API object
	for ind, y := range oncallListing.OnCalls {

		timeConversionStart, err := utils.FormatTimestamp(y.Start)

		if err != nil {
			return nil, err
		}

		timeConversionEnd, err := utils.FormatTimestamp(y.End)

		if err != nil {
			return nil, err
		}

		// Parse the oncall data to OncallUser object
		if strings.Contains(y.Schedule.Summary, "Management") {
			tempUser := OncallUser{}
			tempUser.EscalationPolicy = y.EscalationPolicy.Summary
			tempUser.OncallRole = y.Schedule.Summary
			tempUser.Name = y.User.Summary
			tempUser.Start = timeConversionStart
			tempUser.End = timeConversionEnd
			mgmtUsers = append(mgmtUsers, tempUser)
			startTime, _ = utils.FormatTimestamp(oncallListing.OnCalls[ind+1].Start)
			continue
		}

		if timeConversionStart != startTime {
			temp = append(temp, mgmtUsers...)
			oncallStartTime := startTime[11:16]
			var layerId string
			switch oncallStartTime {
			case "22:30":
				layerId = "Layer 1 [ APAC-E ]"
			case "03:30":
				layerId = "Layer 2 [ APAC-W ]"
			case "08:30":
				layerId = "Layer 3 [ EMEA ]"
			case "13:30":
				layerId = "Layer 4 [ NASA-E ]"
			case "18:00":
				layerId = "Layer 5 [ NASA-W ]"
			}
			oncallLayer := &OncallLayer{
				LayerId: layerId,
				Users:   temp,
			}
			oncallLayers = append(oncallLayers, *oncallLayer)
			temp = nil
			startTime = timeConversionStart
		}

		tempUser := OncallUser{}
		tempUser.EscalationPolicy = y.EscalationPolicy.Summary
		tempUser.OncallRole = y.Schedule.Summary
		tempUser.Name = y.User.Summary
		tempUser.Start = timeConversionStart
		tempUser.End = timeConversionEnd
		temp = append(temp, tempUser)
	}
	return oncallLayers, err
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
func UserNextOncallSchedule(c client.PagerDutyClient, userID string) ([]OncallUser, error) {
	var callOpts pagerduty.ListOnCallOptions
	var nextOncallData []OncallUser

	callOpts.Until = time.Now().AddDate(0, 3, 0).String()

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
