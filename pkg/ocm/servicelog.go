/*
Copyright © 2023 Red Hat, Inc

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

package ocm

import (
	"fmt"

	"github.com/openshift-online/ocm-cli/pkg/ocm"
	v1 "github.com/openshift-online/ocm-sdk-go/servicelogs/v1"
)

// Local constants
const (
	slFormatSpecifier = "Service Name: %s\nSeverity: %s\nCluster ID: %s\nCluster UUID: %s\nSummary: %s\nDescription: %s\nCreated At: %s\nInternal Only: %s\n"
	serviceName       = "SREManualAction"
	page              = 1
	listSize          = 50
)

// GetClusterServiceLogs retrieves the service log items for a given cluster UUID
func GetClusterServiceLogs(clusterID string) (*v1.LogEntryList, error) {
	// Create OCM client
	connection, err := ocm.NewConnection().Build()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OCM client: %v", err)
	}
	defer connection.Close()

	// Send a GET request to the service log API endpoint
	// Swagger codegen can be looked up on
	// https://api.openshift.com/?urls.primaryName=Service%20logs#/default/get_api_service_logs_v1_cluster_logs
	response, err := connection.ServiceLogs().V1().Clusters().
		Cluster(clusterID).
		ClusterLogs().
		List().
		Search(fmt.Sprintf("service_name = '%s'", serviceName)).
		Page(page).
		Size(listSize).
		Send()

	if err != nil {
		return nil, err
	}

	return response.Items(), nil
}

// ParseServiceLogItems parses the servicelog body into a human-readable format
func ParseServiceLogItems(items *v1.LogEntryList) string {
	var parsedServicelogs, internalOnly string

	for _, item := range items.Slice() {
		internalOnly = "Fasle"
		if item.InternalOnly() {
			internalOnly = "True"
		}

		parsedServicelogs = parsedServicelogs + fmt.Sprintf(slFormatSpecifier,
			item.ServiceName(),
			item.Severity(),
			item.ClusterID(),
			item.ClusterUUID(),
			item.Summary(),
			item.Description(),
			item.Timestamp(),
			internalOnly,
		) + "---------------\n"
	}

	return parsedServicelogs
}
