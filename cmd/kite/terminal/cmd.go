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
package terminal

import (
	"fmt"
	"os"

	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/openshift/pagerduty-short-circuiter/pkg/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "terminal",
	Short: "This command lets the user choose their preferred terminal emulator .",
	Args:  cobra.NoArgs,
	RunE:  selectTerminal,
}

func selectTerminal(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()
	if err != nil {
		err = fmt.Errorf("configuration file not found, run the 'kite login' command")
		return err

	}

	// List all the available terminal emulators and prompt the user to select one.
	selectedTerminal := utils.InitTerminalEmulator()

	if err != nil {
		fmt.Println(err)
		return err
	}
	if selectedTerminal == "" {
		os.Exit(1)
	}
	// Update the terminal configuration value in the config file.
	cfg.Terminal = selectedTerminal

	err = config.Save(cfg)
	if err != nil {
		return err
	}
	fmt.Println("Emulator selected successfully")
	return nil

}
