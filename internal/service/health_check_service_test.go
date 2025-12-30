package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestHealthCheckService_Check tests the Check method of the HealthCheckService.
// It tests the normal case where the Check method returns a tuple of (message, serviceName, instanceID) without any error.
// The test case checks if the returned values match the expected values.
// The test case also checks if the returned instanceID matches the expected instanceID.
// The test case is run in parallel mode to ensure that the test runs quickly.
// It uses the testify package to assert that the returned values match the expected values.
func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name: "normal case",

			inputServiceName: "bookmark-manager",
			inputInstanceID:  "2025",

			expectedMessage:     "OK",
			expectedServiceName: "bookmark-manager",
			expectedInstanceID:  "2025",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID)

			message, serviceName, instanceID := testSvc.Check()

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}

}
