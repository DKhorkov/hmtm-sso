package usecases

import (
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/libs/security"
)

type CommonUseCases struct {
	authService  interfaces.AuthService
	usersService interfaces.UsersService
	hashCost     int
	jwtConfig    security.JWTConfig
}

func (useCases *CommonUseCases) RegisterUser(userData entities.RegisterUserDTO) (uint64, error) {
	hashedPassword, err := security.Hash(userData.Credentials.Password, useCases.hashCost)
	if err != nil {
		return 0, err
	}

	userData.Credentials.Password = hashedPassword
	return useCases.authService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData entities.LoginUserDTO) (*entities.TokensDTO, error) {
	// Check if user with provided email exists and password is valid:
	user, err := useCases.usersService.GetUserByEmail(userData.Email)
	if err != nil {
		return nil, err
	}

	if !security.ValidateHash(userData.Password, user.Password) {
		return nil, &errors.InvalidPasswordError{}
	}

	if dbRefreshToken, err := useCases.authService.GetRefreshTokenByUserID(user.ID); err == nil {
		if err = useCases.authService.ExpireRefreshToken(dbRefreshToken.Value); err != nil {
			return nil, err
		}
	}

	// Create tokens:
	accessToken, err := security.GenerateJWT(
		user.ID,
		useCases.jwtConfig.SecretKey,
		useCases.jwtConfig.AccessTokenTTL,
		useCases.jwtConfig.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := security.GenerateJWT(
		accessToken,
		useCases.jwtConfig.SecretKey,
		useCases.jwtConfig.RefreshTokenTTL,
		useCases.jwtConfig.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	// Save token to Database:
	if _, err = useCases.authService.CreateRefreshToken(
		user.ID,
		refreshToken,
		useCases.jwtConfig.RefreshTokenTTL,
	); err != nil {
		return nil, err
	}

	// Encoding refresh token for secure usage via internet:
	encodedRefreshToken := security.Encode([]byte(refreshToken))
	return &entities.TokensDTO{
			AccessToken:  accessToken,
			RefreshToken: encodedRefreshToken,
		},
		nil
}

func (useCases *CommonUseCases) GetUserByID(id uint64) (*entities.User, error) {
	return useCases.usersService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]entities.User, error) {
	return useCases.usersService.GetAllUsers()
}

func (useCases *CommonUseCases) GetMe(accessToken string) (*entities.User, error) {
	accessTokenPayload, err := security.ParseJWT(accessToken, useCases.jwtConfig.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	userID := uint64(accessTokenPayload.(float64))
	return useCases.usersService.GetUserByID(userID)
}

func (useCases *CommonUseCases) RefreshTokens(refreshToken string) (*entities.TokensDTO, error) {
	// Decoding refresh token to get original JWT and compare its value with value in Database:
	oldRefreshTokenBytes, err := security.Decode(refreshToken)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving refresh token payload to get access token from refresh token:
	oldRefreshToken := string(oldRefreshTokenBytes)
	refreshTokenPayload, err := security.ParseJWT(oldRefreshToken, useCases.jwtConfig.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	oldAccessToken, ok := refreshTokenPayload.(string)
	if !ok {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving access token payload to get user ID:
	accessTokenPayload, err := security.ParseJWT(oldAccessToken, useCases.jwtConfig.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Selecting refresh token model from Database, if refresh token has not expired yet:
	userID := uint64(accessTokenPayload.(float64))
	dbRefreshToken, err := useCases.authService.GetRefreshTokenByUserID(userID)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Checking if access token belongs to refresh token:
	if oldRefreshToken != dbRefreshToken.Value {
		return nil, &errors.AccessTokenDoesNotBelongToRefreshTokenError{}
	}

	// Expiring old refresh token in Database to have only one valid refresh token instance:
	if err = useCases.authService.ExpireRefreshToken(dbRefreshToken.Value); err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Create tokens:
	newAccessToken, err := security.GenerateJWT(
		userID,
		useCases.jwtConfig.SecretKey,
		useCases.jwtConfig.AccessTokenTTL,
		useCases.jwtConfig.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := security.GenerateJWT(
		newAccessToken,
		useCases.jwtConfig.SecretKey,
		useCases.jwtConfig.RefreshTokenTTL,
		useCases.jwtConfig.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	// Save token to Database:
	if _, err = useCases.authService.CreateRefreshToken(
		userID,
		newRefreshToken,
		useCases.jwtConfig.RefreshTokenTTL,
	); err != nil {
		return nil, err
	}

	// Encoding refresh token for secure usage via internet:
	encodedRefreshToken := security.Encode([]byte(newRefreshToken))
	return &entities.TokensDTO{
			AccessToken:  newAccessToken,
			RefreshToken: encodedRefreshToken,
		},
		nil
}

func NewCommonUseCases(
	authService interfaces.AuthService,
	usersService interfaces.UsersService,
	hashCost int,
	jwtConfig security.JWTConfig,
) *CommonUseCases {
	return &CommonUseCases{
		authService:  authService,
		usersService: usersService,
		hashCost:     hashCost,
		jwtConfig:    jwtConfig,
	}
}
