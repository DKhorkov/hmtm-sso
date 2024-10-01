package mocks

import (
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
	"time"
)

type MockedSsoRepository struct {
	UsersStorage map[int]*entities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Credentials.Email {
			return 0, &customerrors.UserAlreadyExistsError{}
		}
	}

	user := entities.User{
		Email:     userData.Credentials.Email,
		Password:  userData.Credentials.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user.ID = len(repo.UsersStorage) + 1
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
