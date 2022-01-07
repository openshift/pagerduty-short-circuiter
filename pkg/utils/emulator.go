package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
)

// Installed terminal emulator
var Emulator string

// InitTerminalEmulator tries to set a terminal emulator by trying some known terminal emulators.
func InitTerminalEmulator() {
	emulators := []string{
		"gnome-terminal",
		"x-terminal-emulator",
		"mate-terminal",
		"terminator",
		"xfce4-terminal",
		"urxvt",
		"rxvt",
		"termit",
		"Eterm",
		"aterm",
		"uxterm",
		"xterm",
		"roxterm",
		"termite",
		"kitty",
		"hyper",
	}

	for _, t := range emulators {
		cmd := exec.Command("command", "-v", t)
		output, _ := cmd.CombinedOutput()
		cmd.ProcessState.Exited()
		term := string(output)
		term = strings.TrimSpace(term)

		if term != "" {
			Emulator = term
		}
	}
}

// ClusterLoginShell spawns an instance of ocm-container in the same shell.
func ClusterLoginShell(clusterID string) *exec.Cmd {
	// Check if ocm-container is installed locally
	ocmContainer, err := exec.LookPath("ocm-container")

	if err != nil {
		fmt.Println("ocm-container is not found.\nPlease install it via:", constants.OcmContainerURL)
	}

	cmd := exec.Command(ocmContainer, clusterID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// ClusterLoginEmulator spawns an instance of ocm-container in a new terminal.
func ClusterLoginEmulator(clusterID string) error {
	var cmd *exec.Cmd

	// Check if ocm-container is installed locally
	ocmContainer, err := exec.LookPath("ocm-container")

	if err != nil {
		return errors.New("ocm-container is not found.\nPlease install it via: " + constants.OcmContainerURL)
	}

	// OCM container command to be executed for cluster login
	ocmCommand := ocmContainer + " " + clusterID

	cmd = exec.Command(Emulator, "-e", ocmCommand)

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
