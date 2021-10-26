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

// incident retuns a pagerduty incident object with pre-configured data.
func incident(incidentID string) pdApi.Incident {
	return pdApi.Incident{
		Id: incidentID,
	}
}

// alert retuns a pagerduty alert object with pre-configured data.
func alert(incidentID string, serviceID string, name string, clusterID string, severity string, status string) pdApi.IncidentAlert {
	return pdApi.IncidentAlert{

		Incident: pdApi.APIReference{
			ID: incidentID,
		},

		Service: pdApi.APIObject{
			ID: serviceID,
		},

		APIObject: pdApi.APIObject{
			Summary: name,
		},

		Body: map[string]interface{}{
			"details": map[string]interface{}{
				"cluster_id": clusterID,
			},
		},

		Severity:  severity,
		Status:    status,
		CreatedAt: "2006-01-02T15:04:05Z",
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

			expectedIncidents := []pdApi.Incident{
				{Id: "incident-id-1"},
				{Id: "incident-id-2"},
			}

			mockClient.EXPECT().ListIncidents(gomock.Any()).Return(incidentsResponse, nil).Times(1)

			result, err := pdcli.GetIncidents(mockClient, &pdApi.ListIncidentsOptions{})

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(expectedIncidents))
		})
	})

	When("the alert data is fetched", func() {
		It("the cluster name is retrieved from the alert service", func() {

			// Set the mock cluster name
			serviceResponse := &pdApi.Service{
				Description: "my-cluster-name belongs to cluster.101.hive.apps.com",
			}

			mockClient.EXPECT().GetService("", gomock.Any()).Return(serviceResponse, nil).Times(1)

			expectedResult := "my-cluster-name"

			result, err := pdcli.GetClusterName("", mockClient)

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(expectedResult))
		})
	})

	When("the alerts command is run", func() {
		It("returns alerts for an incident", func() {

			alertResponse := &pdApi.ListAlertsResponse{
				APIListObject: pdApi.APIListObject{},
				Alerts: []pdApi.IncidentAlert{
					alert(
						"incident-id-1",
						"my-service-id",
						"alert-name",
						"cluster-id",
						"critical",
						"triggered",
					),
				},
			}

			// Set the mock cluster name
			serviceResponse := &pdApi.Service{
				Description: "my-cluster-name",
			}

			expectedAlert := []pdcli.Alert{
				{
					IncidentID:  "incident-id-1",
					ClusterID:   "cluster-id",
					ClusterName: "my-cluster-name",
					CreatedAt:   "01-02-2006 15:04 UTC",
					Name:        "alert-name",
					Console:     "<nil>",
					Labels:      "<nil>",
					Severity:    "critical",
					Status:      "triggered",
					Sop:         "<nil>",
				},
			}

			mockClient.EXPECT().GetService("my-service-id", gomock.Any()).Return(serviceResponse, nil).Times(1)

			mockClient.EXPECT().ListIncidentAlerts(gomock.Any()).Return(alertResponse, nil).Times(1)

			result, err := pdcli.GetIncidentAlerts(mockClient, "incident-id-5")

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(expectedAlert))
		})
	})

	When("the incident alerts are fetched", func() {
		It("parses the alert data to the struct", func() {

			var alertData pdcli.Alert

			alertResponse := &pdApi.ListAlertsResponse{
				APIListObject: pdApi.APIListObject{},
				Alerts: []pdApi.IncidentAlert{
					alert(
						"incident-id-1",
						"my-service-id",
						"alert-name",
						"cluster-id",
						"critical",
						"triggered",
					),
				},
			}

			// Set the mock cluster name
			serviceResponse := &pdApi.Service{
				Description: "my-cluster-name",
			}

			expectedAlertData := pdcli.Alert{
				IncidentID:  "incident-id-1",
				Name:        "alert-name",
				ClusterID:   "cluster-id",
				ClusterName: "my-cluster-name",
				CreatedAt:   "01-02-2006 15:04 UTC",
				Severity:    "critical",
				Status:      "triggered",
				Console:     "<nil>",
				Labels:      "<nil>",
				Sop:         "<nil>",
			}

			mockClient.EXPECT().GetService("my-service-id", gomock.Any()).Return(serviceResponse, nil).Times(1)

			err := alertData.ParseAlertData(mockClient, &alertResponse.Alerts[0])

			Expect(err).ShouldNot(HaveOccurred())

			Expect(alertData).To(Equal(expectedAlertData))

		})
	})

	When("a user acknowledges an incident(s)", func() {
		It("it changes the incident status to acknowledged and returns the incident(s)", func() {

			userResponse := &pdApi.User{
				APIObject: pdApi.APIObject{
					ID: "my-user-id",
				},
				Email: "example@redhat.com",
			}

			incidentResponse := &pdApi.ListIncidentsResponse{
				Incidents: []pdApi.Incident{
					{
						Id:     "ABC123",
						Status: "acknowledged",
						Acknowledgements: []pdApi.Acknowledgement{
							{
								Acknowledger: userResponse.APIObject,
							},
						},
					},
				},
			}

			expectedResponse := []pdApi.Incident{
				{
					Id:     "ABC123",
					Status: "acknowledged",
					Acknowledgements: []pdApi.Acknowledgement{
						{
							Acknowledger: userResponse.APIObject,
						},
					},
				},
			}

			mockClient.EXPECT().GetCurrentUser(gomock.Any()).Return(userResponse, nil).Times(1)

			mockClient.EXPECT().ManageIncidents(gomock.Any(), gomock.Any()).Return(incidentResponse, nil).Times(1)

			result, err := pdcli.AcknowledgeIncidents(mockClient, []string{"ABC123"})

			Expect(err).ToNot(HaveOccurred())

			Expect(result).To(Equal(expectedResponse))

		})
	})
})
