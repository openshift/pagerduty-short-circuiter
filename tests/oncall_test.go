package tests

import (
	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mockpd "github.com/openshift/pagerduty-short-circuiter/pkg/client/mock"
	pdcli "github.com/openshift/pagerduty-short-circuiter/pkg/pdcli/oncall"
)

var _ = Describe("kite oncall", func() {
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

	When("Configuration file doesn't exist", func() {
		It("throws an error", func() {
			cmd := NewCommand().Args("oncall").Run()

			Expect(cmd.ExitCode()).NotTo(BeZero())
			Expect(cmd.ErrString()).ToNot(BeEmpty())
		})
	})

	When("kite oncall is run", func() {
		It("shows users currently oncall", func() {

			listOnCallsResponse := &pdApi.ListOnCallsResponse{
				OnCalls: []pdApi.OnCall{
					{
						Schedule: pdApi.Schedule{
							APIObject: pdApi.APIObject{
								Summary: "0-SREP Weekday Primary",
							},
						},
						EscalationPolicy: pdApi.EscalationPolicy{
							APIObject: pdApi.APIObject{
								Summary: "Openshift Escalation",
							},
						},
						Start: "2021-10-25T03:30:00Z",
						End:   "2021-10-25T08:30:00Z",
						User: pdApi.User{
							Description: "Red Hat SRE",
						},
					},
				},
			}

			mockClient.EXPECT().ListOnCalls(gomock.Any()).Return(listOnCallsResponse, nil).Times(1)
			// TODO: Fix unit tests
			// expectedResponse := []pdcli.OncallLayer{
			// 	{
			// 		LayerId: "Layer 2 [APAC-W]",
			// 		Users: []pdcli.OncallUser{
			// 			{EscalationPolicy: "Openshift Escalation",
			// 				OncallRole: "0-SREP Weekday Primary",
			// 				Name:       "Red Hat SRE",
			// 				Start:      "10-25-2021 03:30 UTC",
			// 				End:        "10-25-2021 08:30 UTC",
			// 			},
			// 		},
			// 	},
			// }

			_, err := pdcli.TeamSREOnCall(mockClient)
			Expect(err).ToNot(HaveOccurred())
			// TODO: Fix unit tests
			// Expect(result).To(Equal(expectedResponse))
		})
	})
})
