package usecases

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/libs/security"
)

func NewCommonUseCases(
	authService interfaces.AuthService,
	usersService interfaces.UsersService,
	securityConfig security.Config,
	validationConfig config.ValidationConfig,
) *CommonUseCases {
	return &CommonUseCases{
		authService:      authService,
		usersService:     usersService,
		securityConfig:   securityConfig,
		validationConfig: validationConfig,
	}
}

type CommonUseCases struct {
	authService      interfaces.AuthService
	usersService     interfaces.UsersService
	securityConfig   security.Config
	validationConfig config.ValidationConfig
}

func (useCases *CommonUseCases) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	if !validateEmail(userData.Credentials.Email, useCases.validationConfig.EmailRegExp) {
		return 0, &customerrors.InvalidEmailError{}
	}

	if !validatePassword(userData.Credentials.Password, useCases.validationConfig.PasswordRegExps) {
		return 0, &customerrors.InvalidPasswordError{}
	}

	hashedPassword, err := security.Hash(userData.Credentials.Password, useCases.securityConfig.HashCost)
	if err != nil {
		return 0, err
	}

	userData.Credentials.Password = hashedPassword
	return useCases.authService.RegisterUser(ctx, userData)
}

func (useCases *CommonUseCases) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	// Check if user with provided email exists and password is valid:
	user, err := useCases.usersService.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, err
	}

	if !security.ValidateHash(userData.Password, user.Password) {
		return nil, &customerrors.WrongPasswordError{}
	}

	if dbRefreshToken, err := useCases.authService.GetRefreshTokenByUserID(ctx, user.ID); err == nil {
		if err = useCases.authService.ExpireRefreshToken(ctx, dbRefreshToken.Value); err != nil {
			return nil, err
		}
	}

	// Create tokens:
	accessToken, err := security.GenerateJWT(
		user.ID,
		useCases.securityConfig.JWT.SecretKey,
		useCases.securityConfig.JWT.AccessTokenTTL,
		useCases.securityConfig.JWT.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := security.GenerateJWT(
		accessToken,
		useCases.securityConfig.JWT.SecretKey,
		useCases.securityConfig.JWT.RefreshTokenTTL,
		useCases.securityConfig.JWT.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	// Save token to Database:
	if _, err = useCases.authService.CreateRefreshToken(
		ctx,
		user.ID,
		refreshToken,
		useCases.securityConfig.JWT.RefreshTokenTTL,
	); err != nil {
		return nil, err
	}

	// Encoding refresh token for secure usage via internet:
	encodedRefreshToken := security.Encode([]byte(refreshToken))
	return &entities.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

func (useCases *CommonUseCases) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	return useCases.usersService.GetUserByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return useCases.usersService.GetAllUsers(ctx)
}

func (useCases *CommonUseCases) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	accessTokenPayload, err := security.ParseJWT(accessToken, useCases.securityConfig.JWT.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	floatUserID, ok := accessTokenPayload.(float64)
	if !ok {
		return nil, &security.InvalidJWTError{}
	}

	userID := uint64(floatUserID)
	return useCases.usersService.GetUserByID(ctx, userID)
}

func (useCases *CommonUseCases) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	// Decoding refresh token to get original JWT and compare its value with value in Database:
	oldRefreshTokenBytes, err := security.Decode(refreshToken)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving refresh token payload to get access token from refresh token:
	oldRefreshToken := string(oldRefreshTokenBytes)
	refreshTokenPayload, err := security.ParseJWT(oldRefreshToken, useCases.securityConfig.JWT.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	oldAccessToken, ok := refreshTokenPayload.(string)
	if !ok {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving access token payload to get user ID:
	accessTokenPayload, err := security.ParseJWT(oldAccessToken, useCases.securityConfig.JWT.SecretKey)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Selecting refresh token model from Database, if refresh token has not expired yet:
	floatUserID, ok := accessTokenPayload.(float64)
	if !ok {
		return nil, &security.InvalidJWTError{}
	}

	userID := uint64(floatUserID)
	dbRefreshToken, err := useCases.authService.GetRefreshTokenByUserID(ctx, userID)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Checking if access token belongs to refresh token:
	if oldRefreshToken != dbRefreshToken.Value {
		return nil, &customerrors.AccessTokenDoesNotBelongToRefreshTokenError{}
	}

	// Expiring old refresh token in Database to have only one valid refresh token instance:
	if err = useCases.authService.ExpireRefreshToken(ctx, dbRefreshToken.Value); err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Create tokens:
	newAccessToken, err := security.GenerateJWT(
		userID,
		useCases.securityConfig.JWT.SecretKey,
		useCases.securityConfig.JWT.AccessTokenTTL,
		useCases.securityConfig.JWT.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := security.GenerateJWT(
		newAccessToken,
		useCases.securityConfig.JWT.SecretKey,
		useCases.securityConfig.JWT.RefreshTokenTTL,
		useCases.securityConfig.JWT.Algorithm,
	)

	if err != nil {
		return nil, err
	}

	// Save token to Database:
	if _, err = useCases.authService.CreateRefreshToken(
		ctx,
		userID,
		newRefreshToken,
		useCases.securityConfig.JWT.RefreshTokenTTL,
	); err != nil {
		return nil, err
	}

	// Encoding refresh token for secure usage via internet:
	encodedRefreshToken := security.Encode([]byte(newRefreshToken))
	return &entities.TokensDTO{
		AccessToken:  newAccessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}
