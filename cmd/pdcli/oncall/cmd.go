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
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
	"github.com/spf13/cobra"
)

var options struct {
	allTeams   bool
	nextOncall bool
}

var Cmd = &cobra.Command{
	Use:   "oncall",
	Short: "oncall to the PagerDuty CLI",
	Long:  "Running the pdcli oncall command will display the current primary and secondary oncall SRE",
	Args:  cobra.NoArgs,
	RunE:  oncallHandler,
}

func init() {

	// Shows who is on-call in all teams
	Cmd.Flags().BoolVarP(
		&options.allTeams,
		"all",
		"a",
		false,
		"Show who is on-call in all teams",
	)

	// Next oncall
	Cmd.Flags().BoolVar(
		&options.nextOncall,
		"next-oncall",
		false,
		"Show the current user's next oncall schedule",
	)
}

// oncallHandler is the main handler for pdcli oncall.
func oncallHandler(cmd *cobra.Command, args []string) (err error) {

	// Establish a secure connection with the PagerDuty API
	client, err := client.NewClient().Connect()

	if err != nil {
		return err
	}

	switch {
	case options.allTeams:
		// Fetch oncall data from all teams
		err = pdcli.AllTeamsOncall(client)

		if err != nil {
			return err
		}

	case options.nextOncall:
		err = pdcli.UserNextOncallSchedule(client)

		if err != nil {
			return err
		}

	default:
		// Fetch oncall data from Platform-SRE team
		err = pdcli.TeamSREOnCall(client)

		if err != nil {
			return err
		}
	}

	return nil
}
