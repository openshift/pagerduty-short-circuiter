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

	//"github.com/openshift/pagerduty-short-circuiter/pkg/oncall"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "Oncall to the PagerDuty CLI",
	Long:  "Running the pdcli oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	//RunE:  OnCall(),
}

//func init() {

//Cmd.Flags().StringVar(&OncallArgs.apiKey, "key", "", "Access API key/token generated from "+APIKeyURL+"\nUse this option to overwrite the existing API key.")
//}
//type oncall struct{
	//string
//}
func OnCall(cmd*cobra.Command, args [] string) error {

	var call pagerduty.ListOnCallOptions
	
	connection, err := pdcli.NewConnection().Build()
	if err != nil {
		fmt.Println(err)
	}
	etc, err := connection.ListOnCalls(call)
	if err != nil {
		fmt.Println(err)
	}

	for _, y := range etc.OnCalls {

		fmt.Println(y.User)
	}
	for _, y := range etc.OnCalls {

		fmt.Println(y.Schedule)
		

	}

	//fmt.Printf("User: %v\n", User)

	return nil

}


