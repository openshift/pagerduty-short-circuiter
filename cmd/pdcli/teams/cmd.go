package teams

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/openshift/pagerduty-short-circuiter/pkg/client"
	"github.com/openshift/pagerduty-short-circuiter/pkg/config"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "teams",
	Short: "This command will list all the teams associated with your user account.",
	Args:  cobra.NoArgs,
	RunE:  teamsHandler,
}

func teamsHandler(cmd *cobra.Command, args []string) error {

	// PagerDuty API client
	pdClient, err := client.NewClient().Connect()

	if err != nil {
		return err
	}

	// Load the configuration file
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	// Fetch the user selected team ID
	teamID, err := SelectTeam(pdClient, os.Stdin)

	if err != nil {
		return err
	}

	cfg.TeamID = teamID

	// Save the modified configuration
	err = config.Save(cfg)

	if err != nil {
		return err
	}

	fmt.Println("PagerDuty team successfully selected.")

	return nil
}

// SelectTeam prompts the user to select a team and returns the selected team ID.
func SelectTeam(c client.PagerDutyClient, stdin io.Reader) (string, error) {
	var selectedTeam string
	var userOptions pdApi.GetCurrentUserOptions

	userTeams := make(map[string]string)

	// Fetch the currently logged in user details
	user, err := c.GetCurrentUser(userOptions)

	if err != nil {
		return "", err
	}

	// Check if the user belongs to any team
	if len(user.Teams) == 0 {
		fmt.Println("Your user account is currently not a part of any team")
		os.Exit(0)
	}

	for i, team := range user.Teams {
		index := strconv.Itoa(i + 1)
		userTeams[index] = team.ID
		fmt.Printf("%s. %s\n", index, team.Summary)
	}

	fmt.Print("Select Team: ")

	reader := bufio.NewReader(stdin)

	input, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	if val, ok := userTeams[input]; ok {
		selectedTeam = val
	} else {
		return "", fmt.Errorf("please select a valid option")
	}

	return selectedTeam, nil
}
