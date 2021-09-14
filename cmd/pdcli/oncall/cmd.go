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
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/olekukonko/tablewriter"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
	//"github.com/olekukonko/tablewriter"
)

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "oncall to the PagerDuty CLI",
	Long:  "Running the pdcli oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	RunE:  OnCall,
}
//function for getting current primary and secondary oncalls
func OnCall(cmd *cobra.Command, args []string) error {

	var call pagerduty.ListOnCallOptions

	
	

    call.ScheduleIDs = []string{"P995J2A","P4TU2IT"}
	

    connection, err := pdcli.NewConnection().Build()
	if err != nil {
		fmt.Println(err)
	}
	etc,err:=connection.ListOnCalls(call)
	
		if err!=nil{
			return err
		}
	//for getting secondary/primary as per schedule and name
	count := 0
	for _, y  :=  range etc.OnCalls{
		if count==0 || count ==2{

		
		data:=[][]string{
			[]string{y.Schedule.Summary,y.User.Summary},
		} 
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Oncall Role","Name"})
		table.AppendBulk(data)
		
		table.Render()

		
		
	}
	count+=1
}


	return nil

}

