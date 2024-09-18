package mocks

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/entities"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
)

type MockedSsoRepository struct {
	UsersStorage map[int]*entities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	var user entities.User
	user.Email = userData.Credentials.Email
	user.ID = len(repo.UsersStorage) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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
