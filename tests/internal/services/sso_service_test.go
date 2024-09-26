package services__test

import (
	"sort"
	"testing"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	mocks "github.com/DKhorkov/hmtm-sso/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/services"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	const testUserID = 1

	testCases := []struct {
		name     string
		input    entities.RegisterUserDTO
		expected int
		message  string
	}{
		{
			name: "should register a new user",
			input: entities.RegisterUserDTO{
				Credentials: entities.LoginUserDTO{
					Email:    "tests@example.com",
					Password: "password",
				},
			},
			expected: testUserID,
			message:  "should return a new user id",
		},
	}

	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*entities.User{}}
	ssoService := &services.CommonAuthService{AuthRepository: ssoRepository}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ssoService.RegisterUser(tc.input)
			require.NoError(
				t,
				err,
				"%s - error: %v", tc.message, err)
			assert.Equal(
				t,
				tc.expected,
				actual,
				"\n%s - actual: %v, expected: %v", tc.message, actual, tc.expected)
		})
	}
}

func TestGetAllUsersWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*entities.User{}}
	ssoService := &services.CommonUsersService{UsersRepository: ssoRepository}

	users, err := ssoService.GetAllUsers()
	require.NoError(t, err, "Should return no error")
	assert.Empty(t, users, "Should return an empty list")
}

func TestGetAllUsersWithExistingUsers(t *testing.T) {
	testUsers := [3]entities.RegisterUserDTO{
		{
			Credentials: entities.LoginUserDTO{
				Email:    "test1@example.com",
				Password: "password1",
			},
		},
		{
			Credentials: entities.LoginUserDTO{
				Email:    "test2@example.com",
				Password: "password2",
			},
		},
		{
			Credentials: entities.LoginUserDTO{
				Email:    "test3@example.com",
				Password: "password3",
			},
		},
	}

	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*entities.User{}}
	ssoService := &services.CommonUsersService{UsersRepository: ssoRepository}
	authService := services.CommonAuthService{AuthRepository: ssoRepository}
	for index, userData := range testUsers {
		registeredUserID, err := authService.RegisterUser(userData)
		require.NoError(t, err, "Should create user without error")
		assert.Equal(t, registeredUserID, index+1, "Should return correct ID for registered user")
	}

	users, err := ssoService.GetAllUsers()
	require.NoError(t, err, "Should return no error")
	assert.Len(t, users, len(testUsers), "Should return correct number of users")

	// Sorting slice of users to avoid IDs and Emails mismatch errors due to slice structure:
	sort.Slice(
		users,
		func(i, j int) bool {
			return users[i].ID < users[j].ID
		},
	)

	for index, user := range users {
		assert.Equal(
			t,
			user.Email,
			testUsers[index].Credentials.Email,
			"Should return correct email for user")
		assert.Equal(
			t,
			user.ID,
			index+1,
			"Should return correct ID for user")
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected *entities.User
	}{
		{
			name:     "should find user by ID",
			input:    4,
			expected: &entities.User{ID: 4, Email: "test@example4.com"},
		},
	}

	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: map[int]*entities.User{
			1: {ID: 1, Email: "test@example.com"},
			2: {ID: 2, Email: "test@example2.com"},
			3: {ID: 3, Email: "test@example3.com"},
			4: {ID: 4, Email: "test@example4.com"},
		},
	}

	ssoService := &services.CommonUsersService{UsersRepository: ssoRepository}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ssoService.GetUserByID(tc.input)
			require.NoError(t, err, "Should return no error")
			assert.Equal(
				t,
				tc.expected,
				actual,
				"\n%s - actual: %v, expected: %v", tc.name, actual, tc.expected)
		})
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	const testUserID = 1

	ssoService := &services.CommonUsersService{
		UsersRepository: &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{},
		},
	}

	userID, err := ssoService.GetUserByID(testUserID)
	assert.Nil(
		t,
		userID,
		"should return nil for user with ID=%d", testUserID)
	assert.IsType(t, &customerrors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestLoginUser(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: map[int]*entities.User{
			1: {
				Email:    "test@example.com",
				Password: "password",
			},
		},
	}

	ssoService := &services.CommonAuthService{
		UsersRepository: ssoRepository,
		AuthRepository:  ssoRepository,
	}

	testCases := []struct {
		name          string
		input         entities.LoginUserDTO
		expected      string
		expectedError error
	}{
		{
			name:          "should return token",
			input:         entities.LoginUserDTO{Email: "test@example.com", Password: "password"},
			expected:      "someToken",
			expectedError: nil,
		},
		{
			name:          "should return error if user not found",
			input:         entities.LoginUserDTO{Email: "nonexistent@example.com", Password: "password"},
			expected:      "",
			expectedError: &customerrors.UserNotFoundError{},
		},
		{
			name:          "should return error if password is incorrect",
			input:         entities.LoginUserDTO{Email: "test@example.com", Password: "wrongPassword"},
			expected:      "",
			expectedError: &customerrors.InvalidPasswordError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ssoService.LoginUser(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err, "Should return an error")
				assert.IsType(t, tc.expectedError, err)
			} else {
				require.NoError(t, err, "Should return no error")
			}

			assert.Equal(
				t,
				tc.expected,
				actual,
				"\n%s - actual: %v, expected: %v", tc.name, actual, tc.expected)
		})
	}
}
