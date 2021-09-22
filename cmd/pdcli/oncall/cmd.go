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
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/olekukonko/tablewriter"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
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

//Oncall implements the fetching of current roles and names of users
func OnCall(cmd *cobra.Command, args []string) error {

	var call pagerduty.ListOnCallOptions

	call.ScheduleIDs = []string{constants.PrimaryScheduleID, constants.SecondaryScheduleID}

	// Establish a secure connection with the PagerDuty API
	connection, err := pdcli.NewConnection().Build()
	if err != nil {
		return err
	}

	oncallListing, err := connection.ListOnCalls(call)

	if err != nil {
		return err
	}

	//oncallMapis a hash map to avoid repetition of data in the data array
	//oncallMap := make(map[string]string)

	//oncallData maintains two key

	var oncallData [][]string

	for _, y := range oncallListing.OnCalls {

		if y.Schedule.Summary == "0-SREP: Weekday Primary" {
			//oncallMap["Primary"] = y.User.Summary,y.Schedule.Summary
			oncallData = append(oncallData, []string{y.Schedule.Summary, y.User.Summary})
		}

		if y.Schedule.Summary == "0-SREP: Weekday Secondary" {
			//oncallMap["Secondary"] = y.User.Summary
			oncallData = append(oncallData, []string{y.Schedule.Summary, y.User.Summary})
		}

	}

	//prints all oncall role and name to the console in a tabular format.
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Oncall Role", "Name", "From", "To"})
	table.AppendBulk(oncallData)
	table.Render()

	return nil
}
