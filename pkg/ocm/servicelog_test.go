package ocm

import (
	"testing"
)

// TODO: More test cases to come
func TestServiceLogs(t *testing.T) {
	t.Run("it returns an error when the cluster ID is empty", func(t *testing.T) {

		clusterID := ""
		_, err := GetClusterServiceLogs(clusterID)
		if err == nil {
			t.Errorf("expected to return an error but got: %v", err)
		}
	})
}
