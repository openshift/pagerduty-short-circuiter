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

	mockpd "github.com/openshift/pagerduty-short-circuiter/pkg/client/mock"
	"github.com/openshift/pagerduty-short-circuiter/pkg/pdcli"
)

func incident(incidentID string) pdApi.Incident {
	return pdApi.Incident{
		Id: incidentID,
	}
}

func alert(incidentID string, name string, clusterName string, clusterID string, severity string, status string) pdApi.IncidentAlert {
	return pdApi.IncidentAlert{
		APIObject: pdApi.APIObject{
			Summary: name,
		},

		Incident: pdApi.APIReference{
			ID: incidentID,
		},

		Service: pdApi.APIObject{
			Summary: clusterName,
		},

		Body: map[string]interface{}{
			"details": map[string]interface{}{
				"cluster_id": clusterID,
			},
		},

		Severity: severity,
		Status:   status,
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
		It("returns incidents", func() {

			incidentsResponse := &pdApi.ListIncidentsResponse{
				Incidents: []pdApi.Incident{
					incident("incident-id-1"),
					incident("incident-id-2"),
				},
			}

			mockClient.EXPECT().ListIncidents(gomock.Any()).Return(incidentsResponse, nil).Times(1)

			result, err := pdcli.GetIncidents(mockClient, &pdApi.ListIncidentsOptions{})

			Expect(err).ShouldNot(HaveOccurred())

			expectedIncidents := []pdApi.Incident{
				{Id: "incident-id-1"},
				{Id: "incident-id-2"},
			}

			Expect(result).Should(Equal(expectedIncidents))
		})
	})

	When("the alerts command is run", func() {
		It("returns alerts for an incident", func() {

			alertResponse := &pdApi.ListAlertsResponse{
				APIListObject: pdApi.APIListObject{},
				Alerts: []pdApi.IncidentAlert{
					alert(
						"incident-id-1",
						"alert-name",
						"my-cluster",
						"cluster-id",
						"critical",
						"triggered",
					),
				},
			}

			serviceResponse := &pdApi.Service{
				Description: "my-cluster",
			}

			expectedAlert := []pdcli.Alert{
				{
					IncidentID:  "incident-id-1",
					AlertID:     "",
					ClusterID:   "cluster-id",
					ClusterName: "my-cluster",
					Name:        "alert-name",
					Console:     "<nil>",
					Labels:      "<nil>",
					LastCheckIn: "",
					Severity:    "critical",
					Status:      "triggered",
					Sop:         "<nil>",
					Token:       "",
					Tags:        "",
				},
			}

			mockClient.EXPECT().GetService("", gomock.Any()).Return(serviceResponse, nil).Times(1)
			mockClient.EXPECT().ListIncidentAlerts(gomock.Any()).Return(alertResponse, nil).Times(1)

			result, err := pdcli.GetIncidentAlerts(mockClient, "incident-id-1")

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(expectedAlert))

		})
	})

})
