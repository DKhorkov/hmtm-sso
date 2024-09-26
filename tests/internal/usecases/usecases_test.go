package usecases_test

import (
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/services"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mocks "github.com/DKhorkov/hmtm-sso/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	const testUserID = 1
	const testUserEmail = "test@example.com"

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*entities.User),
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{AuthService: ssoService}

		userData := entities.RegisterUserDTO{
			Credentials: entities.LoginUserDTO{
				Email:    testUserEmail,
				Password: "password",
			},
		}

		userID, err := useCases.RegisterUser(userData)
		require.NoError(t, err)
		assert.Equal(t, testUserID, userID)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{
				testUserID: {Email: testUserEmail},
			},
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{AuthService: ssoService}

		userData := &entities.RegisterUserDTO{
			Credentials: entities.LoginUserDTO{
				Email:    testUserEmail,
				Password: "password",
			},
		}

		userID, err := useCases.AuthService.RegisterUser(*userData)
		require.Error(t, err)
		assert.Equal(t, 0, userID)
		assert.IsType(t, &customerrors.UserAlreadyExistsError{}, err)
	})
}

func TestLoginUser(t *testing.T) {
	const testUserID = 1
	const testUserEmail = "test@example.com"
	const testUserPassword = "password"

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{
				testUserID: {
					Email:    testUserEmail,
					Password: testUserPassword,
				},
			},
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{AuthService: ssoService}

		userData := entities.LoginUserDTO{
			Email:    testUserEmail,
			Password: testUserPassword,
		}

		token, err := useCases.LoginUser(userData)
		require.NoError(t, err)
		assert.Equal(t, "someToken", token)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*entities.User),
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{AuthService: ssoService}

		userData := entities.LoginUserDTO{
			Email:    testUserEmail,
			Password: testUserPassword,
		}

		token, err := useCases.LoginUser(userData)
		require.Error(t, err)
		assert.Equal(t, "", token)
		assert.IsType(t, &customerrors.UserNotFoundError{}, err)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{
				testUserID: {
					Email:    testUserEmail,
					Password: testUserPassword,
				},
			},
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{AuthService: ssoService}

		userData := entities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "wrongPassword",
		}

		token, err := useCases.LoginUser(userData)
		require.Error(t, err)
		assert.Equal(t, "", token)
		assert.IsType(t, &customerrors.InvalidPasswordError{}, err)
	})
}

func TestGetUserByID(t *testing.T) {
	const testUserID = 1
	const testUserEmail = "test@example.com"

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{
				testUserID: {
					ID:        testUserID,
					Email:     testUserEmail,
					Password:  "password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{
			AuthService:  ssoService,
			UsersService: ssoRepository,
		}

		user, err := useCases.GetUserByID(testUserID)
		require.NoError(t, err)
		assert.Equal(t, testUserEmail, user.Email)
		assert.Equal(t, testUserID, user.ID)
	})
}

func TestGetUserByIDNotFound(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: map[int]*entities.User{
			1: {
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			2: {
				ID:        2,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	ssoService := &services.CommonAuthService{
		UsersRepository: ssoRepository,
		AuthRepository:  ssoRepository,
	}

	useCases := &usecases.CommonUseCases{
		AuthService:  ssoService,
		UsersService: ssoRepository,
	}

	userID := 3
	user, err := useCases.GetUserByID(userID)
	assert.IsType(t, &customerrors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
}

func TestGetAllUsersWithExistingUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*entities.User{
				1: {
					ID:        1,
					Email:     "test@example.com",
					Password:  "password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				2: {
					ID:        2,
					Email:     "test2@example.com",
					Password:  "password2",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonAuthService{
			UsersRepository: ssoRepository,
			AuthRepository:  ssoRepository,
		}

		useCases := &usecases.CommonUseCases{
			AuthService:  ssoService,
			UsersService: ssoRepository,
		}

		users, err := useCases.GetAllUsers()
		require.NoError(t, err)
		assert.Len(
			t,
			users,
			len(ssoRepository.UsersStorage),
			"expected to get %d users, got %d", len(ssoRepository.UsersStorage), len(users))
	})
}

func TestGetAllUsersWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: make(map[int]*entities.User),
	}

	ssoService := &services.CommonAuthService{
		UsersRepository: ssoRepository,
		AuthRepository:  ssoRepository,
	}

	useCases := &usecases.CommonUseCases{
		AuthService:  ssoService,
		UsersService: ssoRepository,
	}

	users, err := useCases.GetAllUsers()
	require.NoError(t, err)
	assert.Empty(t, users)
}
