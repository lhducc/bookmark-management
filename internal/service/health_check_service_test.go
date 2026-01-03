package service

import (
	"context"
	"errors"
	"github.com/lhducc/bookmark-management/internal/repository"
	"github.com/lhducc/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var testConnectError = errors.New("can't connect to redis")

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
		setupMock        func(t *testing.T) repository.HealthCheck

		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
		expectedError       error
	}{
		{
			name: "normal case",

			inputServiceName: "bookmark-manager",
			inputInstanceID:  "2025",
			setupMock: func(t *testing.T) repository.HealthCheck {
				repoMock := mocks.NewHealthCheck(t)
				repoMock.On("Ping", mock.Anything).Return(nil)
				return repoMock
			},

			expectedMessage:     "OK",
			expectedServiceName: "bookmark-manager",
			expectedInstanceID:  "2025",
			expectedError:       nil,
		},
		{
			name: "unhealthy case",

			inputServiceName: "bookmark-manager",
			inputInstanceID:  "2025",
			setupMock: func(t *testing.T) repository.HealthCheck {
				repoMock := mocks.NewHealthCheck(t)
				repoMock.On("Ping", mock.Anything).Return(testConnectError)
				return repoMock
			},

			expectedMessage:     "NOT OK",
			expectedServiceName: "bookmark-manager",
			expectedInstanceID:  "2025",
			expectedError:       testConnectError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			healthCheckRepoMock := tc.setupMock(t)

			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID, healthCheckRepoMock)

			message, serviceName, instanceID, err := testSvc.Check(ctx)

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
