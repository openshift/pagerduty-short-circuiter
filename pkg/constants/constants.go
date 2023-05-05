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

package constants

const (
	ConfigFilepath = "kite/config.json"

	APIKeyURL       = "https://support.pagerduty.com/docs/generating-api-keys#generating-a-personal-rest-api-key"
	AccessTokenURL  = "https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"
	OcmContainerURL = "https://github.com/openshift/ocm-container"
	OcmContainer    = "ocm-container"
	Shell           = "SHELL"

	// Regex
	APIKeyRegex     = "^[a-z|A-Z0-9+_-]{20}$"
	IncidentIdRegex = "^[A-Z0-9]{7,14}$"
	TeamIdRegex     = "^[A-Z0-9]{7}$"

	// Sample API key for testing
	SampleKey = "y_NbAkKc66ryYTWUXYEu"

	// Set limit to number of incidents fetched from pagerduty
	IncidentsLimit          = 10
	TrigerredIncidentsLimit = 25

	// PagerDuty IDs
	TeamID     = "PASPK4G"
	SilentTest = "P8QS6CC"
	NobodySREP = "P53J4TK"

	// Escalation Policy IDs
	SilentTestEscalationPolicyID         = "PCGXUDY"
	CADSilentTestEscalationPolicyID      = "PQXIBX3"
	CADSilentTestStageEscalationPolicyID = "PBWX63A"

	// PagerDuty Incident Statuses
	StatusTriggered    = "triggered"
	StatusAcknowledged = "acknowledged"
	StatusHigh         = "high"
	StatusLow          = "low"

	//ScheduleIDS for fetching oncalls as per pagerduty documentation (https://<host>.pagerduty.com/escalation_policies)
	PrimaryScheduleID   = "P995J2A"
	SecondaryScheduleID = "P4TU2IT"
	OncallIDWeekend     = "P7CC7UN"
	OncallManager       = "P1WFZIG"
	OncallID            = "PA4586M"
	InvestigatorID      = "PWQAANA"
)
