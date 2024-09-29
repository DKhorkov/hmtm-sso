package security__test

import (
	"testing"

	"github.com/DKhorkov/hmtm-sso/internal/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		name          string
		hashCost      int
		password      string
		message       string
		errorExpected bool
	}{
		{
			name:          "password successfully hashed",
			hashCost:      14,
			password:      "password",
			message:       "should return hash for password",
			errorExpected: false,
		},
		{
			name:          "too long password > 72 bytes",
			hashCost:      14,
			password:      "tooLongPasswordThatCanNotBeLessThanSeventyTwoBytesForSureAndThereCouldAlsoBeSomeStory",
			message:       "should return error",
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := security.HashPassword(tc.password, tc.hashCost)

			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.Equal(
					t,
					"",
					hashedPassword,
					"\n%s - actual: '%v', expected: '%v'", tc.message, hashedPassword, "")
			} else {
				require.NoError(t, err, tc.message)
				assert.NotEqual(
					t,
					"",
					hashedPassword,
					"\n%s - actual: '%v', expected: '%v'", tc.message, hashedPassword, "SomeHashedValue")
			}
		})
	}
}

func TestValidateHashedPassword(t *testing.T) {
	passwordToHash := "password"
	testCases := []struct {
		name     string
		expected bool
		password string
		message  string
	}{
		{
			name:     "hashed password was created based on provided password",
			password: passwordToHash,
			expected: true,
			message:  "should return true",
		},
		{
			name:     "hash password was not created based on provided password\"",
			password: "IncorrectPassword",
			expected: false,
			message:  "should return false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, _ := security.HashPassword(passwordToHash, 0)
			passwordIsValid := security.ValidateHashedPassword(tc.password, hashedPassword)

			assert.Equal(
				t,
				tc.expected,
				passwordIsValid,
				"\n%s - actual: '%v', expected: '%v'", tc.message, passwordIsValid, tc.expected)
		})
	}
}
