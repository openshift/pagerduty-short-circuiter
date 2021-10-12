package tests

import (
	"bytes"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/teams"
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

	When("pdcli teams is run", func() {
		It("prompts the user to select a team and retuns the team ID", func() {

			userResponse := &pdApi.User{
				Name: "my-user",
				Teams: []pdApi.Team{
					{Name: "my-team-a", APIObject: pdApi.APIObject{ID: "ABCD123"}},
					{Name: "my-team-b", APIObject: pdApi.APIObject{ID: "EFGH456"}},
					{Name: "my-team-c", APIObject: pdApi.APIObject{ID: "IJKL789"}},
				},
			}

			expectedResult := "EFGH456"

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)

			var stdin bytes.Buffer

			stdin.Write([]byte("2\n"))

			result, err := teams.SelectTeam(mockClient, &stdin)

			Expect(err).ToNot(HaveOccurred())

			Expect(result).To(Equal(expectedResult))
		})
	})

	When("a user enters an option which is not valid (or) is out of bounds", func() {
		It("throws an error", func() {

			userResponse := &pdApi.User{
				Name: "my-user",
				Teams: []pdApi.Team{
					{Name: "my-team-a", APIObject: pdApi.APIObject{ID: "ABCD123"}},
					{Name: "my-team-b", APIObject: pdApi.APIObject{ID: "EFGH456"}},
					{Name: "my-team-c", APIObject: pdApi.APIObject{ID: "IJKL789"}},
				},
			}

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)

			var stdin bytes.Buffer

			stdin.Write([]byte("X\n"))

			result, err := teams.SelectTeam(mockClient, &stdin)

			Expect(err).To(HaveOccurred())

			Expect(result).To(BeEmpty())
		})
	})
})
