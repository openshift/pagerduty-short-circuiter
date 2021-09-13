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

	"github.com/PagerDuty/go-pagerduty"
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

//func init() {

//}
func OnCall(cmd *cobra.Command, args []string) error {

	var call pagerduty.ListOnCallOptions

	//var call pagerduty.ListSchedulesOptions

	//call.EscalationPolicyIDs = []string{"PA4586M"}
	//call.ScheduleIDs = []string{"P995J2A"}
	//fmt.Println(call.EscalationPolicyIDs)

	connection, err := pdcli.NewConnection().Build()

	if err != nil {
		fmt.Println(err)
	}

	etc, err := connection.ListOnCalls(call)
	if err != nil {
		return err
	}
	for _, y := range etc.OnCalls {
		//fmt.Println(y.Schedule.Name)
		//fmt.Println(y.Schedule)
		fmt.Println(y.User.ID)
		fmt.Println(y.Schedule.ID)
	}


	return nil

}
