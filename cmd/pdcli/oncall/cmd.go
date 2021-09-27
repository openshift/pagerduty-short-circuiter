/*
Copyright © 2021 Red Hat, Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oncall

import (
	"fmt"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
	"github.com/openshift/pagerduty-short-circuiter/pkg/output"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "oncall to the PagerDuty CLI",
	Long:  "Running the pdcli oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	RunE:  OnCall,
}

type User struct {
	EscalationPolicy string
	OncallRole       string
	Name             string
	Start            string
	End              string
}

//Oncall implements the fetching of current roles and names of users
func OnCall(cmd *cobra.Command, args []string) error {

	var call pagerduty.ListOnCallOptions
	call.ScheduleIDs = []string{constants.PrimaryScheduleID, constants.SecondaryScheduleID, constants.OncallManager}
	// Establish a secure connection with the PagerDuty API
	client, err := pdcli.NewConnection().Build()
	if err != nil {
		return nil
	}

	oncallListing, err := client.ListOnCalls(call)

	if err != nil {
		return nil

	}

	//oncallMap is used to store Primary and Secondary oncall information
	oncallMap := map[string]map[string]string{}
	var oncallData []User

	//OnCalls array contains all information about the API object
	for _, y := range oncallListing.OnCalls {

		s := y.EscalationPolicy.Summary + y.Schedule.Summary

		timeConversionStart := timeConversion(y.Start)
		timeConversionEnd := timeConversion(y.End)
		if _, ok := oncallMap[s]; ok {
			continue
		} else {
			oncallMap[s] = map[string]string{}
			oncallMap[s]["Escalation Policy"] = removespace(y.EscalationPolicy.Summary)
			oncallMap[s]["Oncall Role"] = removespace(y.Schedule.Summary)
			oncallMap[s]["Name"] = removespace(y.User.Summary)
			oncallMap[s]["Start"] = removespace(timeConversionStart)
			oncallMap[s]["End"] = removespace(timeConversionEnd)

		}

		for q, r := range oncallMap[s] {

			temp := User{}

			if q == "Escalation Policy" {
				temp.EscalationPolicy = removespace(r)

			}
			if q == "Oncall Role" {
				temp.OncallRole = removespace(r)

			}
			if q == "Name" {
				temp.Name = removespace(r)

			}
			if q == "Start" {
				temp.Start = removespace(r)

			}
			if q == "End" {
				temp.End = removespace(r)

			}

			oncallData = append(oncallData, temp)

		}
	}
	fmt.Println(oncallData)

	printOncalls(oncallData)

	return nil

}

//timeConversion converts timestamp into time and date
func timeConversion(s string) string {

	timeString := s
	timeConverted, err := time.Parse(time.RFC3339, timeString)

	if err != nil {
		fmt.Println(err)
	}
	finalTimeString := timeConverted.String()
	finalTimeString = strings.ReplaceAll(finalTimeString, " +0000 UTC", " UTC")

	return finalTimeString

}

//printOncalls prints data in a tabular form
func printOncalls(oncallData []User) {
	var printData []string
	table := output.NewTable()
	headers := []string{"Escalation Policy", "Oncall Role", "Name", "Start", "End"}
	table.SetHeaders(headers)
	for _, v := range oncallData {
		printData = []string{v.EscalationPolicy, v.OncallRole, v.Name, v.Start, v.End}
		table.AddRow(printData)
	}

	table.SetData()
	table.Print()

}
func removespace(s string) string {
	space := strings.TrimSpace(s)
	return space
}
