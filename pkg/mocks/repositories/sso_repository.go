package mocks

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
)

type MockedSsoRepository struct {
	UsersStorage         map[int]*entities.User
	RefreshTokensStorage map[int]*entities.RefreshToken
}

func (repo *MockedSsoRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Credentials.Email {
			return 0, &customerrors.UserAlreadyExistsError{}
		}
	}

	user := entities.User{
		ID:        len(repo.UsersStorage) + 1,
		Email:     userData.Credentials.Email,
		Password:  userData.Credentials.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.UsersStorage[user.ID] = &user
	return user.ID, nil
}

func (repo *MockedSsoRepository) GetUserByID(id int) (*entities.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, &customerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) GetAllUsers() ([]*entities.User, error) {
	var users []*entities.User
	for _, user := range repo.UsersStorage {
		users = append(users, user)
	}

	return users, nil
}

func (repo *MockedSsoRepository) GetUserByEmail(email string) (*entities.User, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, &customerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) CreateRefreshToken(
	userID int,
	refreshToken string,
	ttl time.Duration,
) (int, error) {
	for _, token := range repo.RefreshTokensStorage {
		if token.Value == refreshToken && time.Now().Before(token.TTL) {
			return 0, &customerrors.RefreshTokenAlreadyExistsError{}
		}
	}

	token := entities.RefreshToken{
		ID:        len(repo.UsersStorage) + 1,
		Value:     refreshToken,
		TTL:       time.Now().Add(ttl),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	repo.RefreshTokensStorage[token.ID] = &token
	return token.ID, nil
}

func (repo *MockedSsoRepository) GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error) {
	for _, token := range repo.RefreshTokensStorage {
		if token.UserID == userID {
			return token, nil
		}
	}

	return nil, &customerrors.RefreshTokenNotFoundError{}
}

func (repo *MockedSsoRepository) ExpireRefreshToken(refreshToken string) error {
	for _, token := range repo.RefreshTokensStorage {
		if token.Value == refreshToken {
			token.TTL = time.Now().Add(time.Hour * time.Duration(-24))
			repo.RefreshTokensStorage[token.ID] = token
			return nil
		}
	}

	return &customerrors.RefreshTokenNotFoundError{}
}
