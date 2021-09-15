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

package tests

import (
	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/alerts"
	mockpd "github.com/openshift/pagerduty-short-circuiter/pkg/client/mock"
)

func incident(id string, name string) pdApi.Incident {
	return pdApi.Incident{
		Id:    id,
		Title: name,
	}
}

func alert(name string, clusterID string) pdApi.IncidentAlert {
	return pdApi.IncidentAlert{
		APIObject: pdApi.APIObject{
			Summary: name,
		},

		Body: map[string]interface{}{
			"details": map[string]interface{}{
				"cluster_id": clusterID,
			},
		},
	}
}

var _ = Describe("view alerts", func() {
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

	When("configuration file doesn't exist", func() {
		It("throws an error", func() {
			cmd := NewCommand().Args("alerts").Run()

			Expect(cmd.ExitCode()).NotTo(BeZero())
			Expect(cmd.ErrString()).ToNot(BeEmpty())
		})
	})

	When("the alerts command is run", func() {
		It("returns incident IDs", func() {
			incidentsResponse := &pdApi.ListIncidentsResponse{
				Incidents: []pdApi.Incident{
					incident("T1234Z", "test-incident-a"),
					incident("R1234H", "test-incident-b"),
				},
			}

			userResponse := &pdApi.User{}

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)
			mockClient.EXPECT().ListIncidents(gomock.Any()).Return(incidentsResponse, nil).Times(1)

			result, err := alerts.GetIncidents(mockClient)

			Expect(err).ShouldNot(HaveOccurred())

			expectedIncidentIDs := []string{
				"T1234Z",
				"R1234H",
			}

			Expect(result).Should(Equal(expectedIncidentIDs))

		})
	})

	When("the alerts command is run", func() {
		It("returns alerts for an incident", func() {

			alertResponse := &pdApi.ListAlertsResponse{
				APIListObject: pdApi.APIListObject{},
				Alerts: []pdApi.IncidentAlert{
					alert("alert-title", "test-cluster-id"),
				},
			}

			expectedAlert := []alerts.Alert{
				{
					IncidentID: "T1234Z",
					Name:       "alert-title",
					ClusterID:  "test-cluster-id",
					Severity:   "",
					Status:     "",
				},
			}

			mockClient.EXPECT().ListIncidentAlerts(gomock.Any()).Return(alertResponse, nil).Times(1)

			result, err := alerts.GetIncidentAlerts(mockClient, "T1234Z")

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(expectedAlert))

		})
	})

})