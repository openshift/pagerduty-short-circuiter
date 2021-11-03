package tests

import (
	"bytes"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/pagerduty-short-circuiter/cmd/kite/teams"
	mockpd "github.com/openshift/pagerduty-short-circuiter/pkg/client/mock"
)

var _ = Describe("select team", func() {
	var (
		mockCtrl   *gomock.Controller
		mockClient *mockpd.MockPagerDutyClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mockpd.NewMockPagerDutyClient(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	When("kite teams is run", func() {
		It("prompts the user to select a team and retuns the team ID", func() {

			userResponse := &pdApi.User{
				Name: "my-user",
				Teams: []pdApi.Team{
					{APIObject: pdApi.APIObject{ID: "ABCD123", Summary: "my-team-a"}},
					{APIObject: pdApi.APIObject{ID: "EFGH456", Summary: "my-team-b"}},
					{APIObject: pdApi.APIObject{ID: "IJKL789", Summary: "my-team-c"}},
				},
			}

			expectedTeamID := "EFGH456"
			expectedTeamName := "my-team-b"

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)

			var stdin bytes.Buffer

			stdin.Write([]byte("2\n"))

			teamID, teamName, err := teams.SelectTeam(mockClient, &stdin)

			Expect(err).ToNot(HaveOccurred())

			Expect(teamID).To(Equal(expectedTeamID))

			Expect(teamName).To(Equal(expectedTeamName))
		})
	})

	When("a user enters an option which is not valid (or) is out of bounds", func() {
		It("throws an error", func() {

			userResponse := &pdApi.User{
				Name: "my-user",
				Teams: []pdApi.Team{
					{APIObject: pdApi.APIObject{ID: "ABCD123", Summary: "my-team-a"}},
					{APIObject: pdApi.APIObject{ID: "EFGH456", Summary: "my-team-b"}},
					{APIObject: pdApi.APIObject{ID: "IJKL789", Summary: "my-team-c"}},
				},
			}

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)

			var stdin bytes.Buffer

			stdin.Write([]byte("X\n"))

			teamID, teamName, err := teams.SelectTeam(mockClient, &stdin)

			Expect(err).To(HaveOccurred())

			Expect(teamID).To(BeEmpty())

			Expect(teamName).To(BeEmpty())
		})
	})
})
