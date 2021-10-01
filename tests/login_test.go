package tests

import (
	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/login"
	mockpd "github.com/openshift/pagerduty-short-circuiter/pkg/client/mock"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
)

var _ = Describe("login", func() {

	var (
		mockCtrl     *gomock.Controller
		mockClient   *mockpd.MockPagerDutyClient
		jsonTemplate string
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mockpd.NewMockPagerDutyClient(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	BeforeEach(func() {
		jsonTemplate = `{ 
			"api_key" : "` + constants.SampleKey + `"
			}`
	})

	When("the login command is run", func() {
		It("creates a configuration file", func() {
			result := NewCommand().
				Args(
					"login",
				).
				GetStdIn(constants.SampleKey + "\n").
				Run()

			Expect(result.ConfigFile()).ToNot(BeEmpty())
			Expect(result.ConfigString()).To(MatchJSON(jsonTemplate))
		})
	})

	When("the login command is run with the 'api-key' option", func() {
		It("creates a configuration file", func() {

			result := NewCommand().
				Args(
					"login",
					"--api-key",
					constants.SampleKey,
				).
				Run()

			Expect(result.ConfigFile()).ToNot(BeEmpty())
			Expect(result.ConfigString()).To(MatchJSON(jsonTemplate))
		})
	})

	When("a user tries to login with an invalid API key", func() {
		It("it doesn't login and throws an error", func() {

			invalidKey := "ABCDEFGHIJ12345"

			result := NewCommand().
				Args(
					"login",
					"--api-key",
					invalidKey,
				).
				Run()

			Expect(result.ErrString()).ToNot(BeEmpty())
			Expect(result.ConfigString()).To(BeEmpty())
		})
	})

	When("a user tries to login with a valid API key", func() {
		It("it sucessfully logs into pdcli and returns the username", func() {

			loginResponse := &pdApi.User{
				Name: "my-user",
			}

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(loginResponse, nil).Times(1)

			user, err := login.Login(constants.SampleKey, mockClient)

			Expect(err).ToNot(HaveOccurred())

			Expect(user).To(Equal(loginResponse.Name))
		})
	})
})
