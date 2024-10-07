package security__test

import (
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"

	"github.com/DKhorkov/hmtm-sso/internal/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityGenerateJWT(t *testing.T) {
	testCases := []struct {
		name          string
		secretKey     string
		algorithm     string
		ttl           time.Duration
		message       string
		errorExpected bool
	}{
		{
			name:          "should generate valid token",
			secretKey:     "testSecret",
			algorithm:     "HS256",
			ttl:           time.Hour,
			message:       "should return valid JWT token",
			errorExpected: false,
		},
	}

	user := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := security.GenerateJWT(user, tc.secretKey, tc.ttl, tc.algorithm)

			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.Equal(
					t,
					"",
					token,
					"\n%s - actual: '%v', expected: '%v'", tc.message, token, "")
			} else {
				require.NoError(t, err, tc.message)
				assert.NotEqual(
					t,
					"",
					token,
					"\n%s - actual: '%v', expected: '%v'", tc.message, token, "SomeJWTValue")
			}
		})
	}
}

func TestSecurityParseJWT(t *testing.T) {
	secretKey := "testSecret"
	testCases := []struct {
		name          string
		secretKey     string
		algorithm     string
		ttl           time.Duration
		message       string
		errorExpected bool
		errorType     error
		expected      int
	}{
		{
			name:          "correct JWT",
			secretKey:     secretKey,
			algorithm:     "HS256",
			ttl:           time.Hour,
			message:       "should return valid JWT token",
			errorExpected: false,
			errorType:     nil,
			expected:      1,
		},
		{
			name:          "invalid secret key",
			secretKey:     "invalidSecret",
			algorithm:     "HS256",
			ttl:           time.Hour,
			message:       "should raise an error due to invalid secret key",
			errorExpected: true,
			errorType:     &customerrors.InvalidJWTError{},
			expected:      0,
		},
		{
			name:          "expired JWT",
			secretKey:     secretKey,
			algorithm:     "HS256",
			ttl:           time.Duration(0),
			message:       "should raise an error due to expired JWT",
			errorExpected: true,
			errorType:     &customerrors.InvalidJWTError{},
			expected:      0,
		},
	}

	user := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := security.GenerateJWT(user, secretKey, tc.ttl, tc.algorithm)
			require.NoError(t, err, tc.message)

			userID, err := security.ParseJWT(token, tc.secretKey)
			assert.Equal(
				t,
				tc.expected,
				userID,
				"\n%s - actual: '%v', expected: '%v'", tc.message, userID, tc.expected)

			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.IsType(t, tc.errorType, err)
			}
		})
	}
}
