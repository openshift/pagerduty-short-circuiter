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
package main

import (
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/alerts"
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/login"
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/oncall"
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/teams"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "pdcli",
	Short:         "A CLI application called pdcli short for PagerDuty CLI.",
	Long:          `It can be used reduce the time taken, from the time, SRE receives a PD alert to the time where troubleshooting on the cluster actually begins. `,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(login.Cmd)
	rootCmd.AddCommand(alerts.Cmd)
	rootCmd.AddCommand(oncall.Cmd)
	rootCmd.AddCommand(teams.Cmd)

	//Do not provide the default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

}
