package service

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var urlSafeRegex = regexp.MustCompile(`^[A-Za-z0-9]+$`)

// TestPasswordService_GeneratePassword tests the GeneratePassword method of the passwordService.
// It tests that the method returns a password of the correct length and that it does not return an error.
// It also tests that the generated password matches the url safe regex.
func TestPasswordService_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLen int
		expectErr   error
	}{
		{
			name: "normal case",

			expectedLen: 10,
			expectErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewPassword()

			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLen, len(pass))
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(pass), true)
		})
	}
}
