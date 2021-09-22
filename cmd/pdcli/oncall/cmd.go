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
	OncallRole string
	Name       string
	Start      string
	End        string
}

//Oncall implements the fetching of current roles and names of users
func OnCall(cmd *cobra.Command, args []string) error {

	var call pagerduty.ListOnCallOptions

	call.ScheduleIDs = []string{constants.PrimaryScheduleID, constants.SecondaryScheduleID}

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

	//OnCalls array contains all information about the API object
	for _, y := range oncallListing.OnCalls {

		//Storing information about Primary Role in oncallMap
		if y.Schedule.Summary == "0-SREP: Weekday Primary" {
			oncallMap["Primary"] = map[string]string{}

			//Converting Start and End timestamps to date and time
			timeConversionStart := timeConversion(y.Start)
			timeConversionEnd := timeConversion(y.End)

			oncallMap["Primary"]["Oncall Role"] = y.Schedule.Summary
			oncallMap["Primary"]["Name"] = y.User.Summary
			oncallMap["Primary"]["Start"] = timeConversionStart
			oncallMap["Primary"]["End"] = timeConversionEnd
		}

		//Storing information about Secondary Role in oncallMap
		if y.Schedule.Summary == "0-SREP: Weekday Secondary" {
			oncallMap["Secondary"] = map[string]string{}

			//Converting Start and End timestamps to date and time
			timeConversionStart := timeConversion(y.Start)
			timeConversionEnd := timeConversion(y.End)

			oncallMap["Secondary"]["Oncall Role"] = y.Schedule.Summary
			oncallMap["Secondary"]["Name"] = y.User.Summary
			oncallMap["Secondary"]["Start"] = timeConversionStart
			oncallMap["Secondary"]["End"] = timeConversionEnd
		}

	}
	data := storeData(oncallMap)
	printOncalls(data)

	return nil

}

func storeData(oncallMap map[string]map[string]string) []User {
	var oncallData []User
	tempCallObjectP := User{}
	tempCallObjectS := User{}

	for x, y := range oncallMap["Primary"] {
		if x == "Name" {
			tempCallObjectP.Name = y
		}
		if x == "Oncall Role" {
			tempCallObjectP.OncallRole = y
		}
		if x == "Start" {
			tempCallObjectP.Start = y
		}
		if x == "End" {
			tempCallObjectP.End = y
		}
	}
	oncallData = append(oncallData, tempCallObjectP)

	for x, y := range oncallMap["Secondary"] {
		if x == "Name" {
			tempCallObjectS.Name = y
		}
		if x == "Oncall Role" {
			tempCallObjectS.OncallRole = y
		}
		if x == "Start" {
			tempCallObjectS.Start = y
		}
		if x == "End" {
			tempCallObjectS.End = y
		}
	}
	oncallData = append(oncallData, tempCallObjectS)

	return oncallData

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
	headers := []string{"Oncall Role", "Name", "From", "To"}
	table.SetHeaders(headers)
	for _, v := range oncallData {
		printData = []string{v.OncallRole, v.Name, v.Start, v.End}
		table.AddRow(printData)
	}
	table.SetData()
	table.Print()
}
