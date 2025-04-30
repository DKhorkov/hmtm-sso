package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"github.com/DKhorkov/libs/validation"
	"github.com/golang-jwt/jwt/v5"

	notifications "github.com/DKhorkov/hmtm-notifications/dto"
	customnats "github.com/DKhorkov/libs/nats"

	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

func New(
	authService interfaces.AuthService,
	usersService interfaces.UsersService,
	securityConfig security.Config,
	validationConfig config.ValidationConfig,
	natsPublisher customnats.Publisher,
	natsConfig config.NATSConfig,
	logger logging.Logger,
) *UseCases {
	return &UseCases{
		authService:      authService,
		usersService:     usersService,
		securityConfig:   securityConfig,
		validationConfig: validationConfig,
		natsPublisher:    natsPublisher,
		natsConfig:       natsConfig,
		logger:           logger,
	}
}

type UseCases struct {
	authService      interfaces.AuthService
	usersService     interfaces.UsersService
	securityConfig   security.Config
	validationConfig config.ValidationConfig
	natsPublisher    customnats.Publisher
	natsConfig       config.NATSConfig
	logger           logging.Logger
}

func (useCases *UseCases) RegisterUser(
	ctx context.Context,
	userData entities.RegisterUserDTO,
) (uint64, error) {
	if !validation.ValidateValueByRule(userData.Email, useCases.validationConfig.EmailRegExp) {
		return 0, &customerrors.InvalidEmailError{}
	}

	if !validation.ValidateValueByRules(userData.Password, useCases.validationConfig.PasswordRegExps) {
		return 0, &customerrors.InvalidPasswordError{}
	}

	if !validation.ValidateValueByRules(userData.DisplayName, useCases.validationConfig.DisplayNameRegExps) {
		return 0, &customerrors.InvalidDisplayNameError{}
	}

	hashedPassword, err := security.Hash(userData.Password, useCases.securityConfig.HashCost)
	if err != nil {
		return 0, err
	}

	userData.Password = hashedPassword

	userID, err := useCases.authService.RegisterUser(ctx, userData)
	if err != nil {
		return 0, err
	}

	verifyEmailDTO := &notifications.VerifyEmailDTO{
		UserID: userID,
	}

	content, err := json.Marshal(verifyEmailDTO)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying to encode data for email verification for User with ID=%d",
				userID,
			),
			err,
		)
	}

	if err = useCases.natsPublisher.Publish(useCases.natsConfig.Subjects.VerifyEmail, content); err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying send Verfiy Email message to User with ID=%d",
				userID,
			),
			err,
		)
	}

	return userID, nil
}

func (useCases *UseCases) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	// Check if user with provided email exists and password is valid:
	user, err := useCases.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, err
	}

	if !user.EmailConfirmed {
		return nil, &customerrors.EmailIsNotConfirmedError{}
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
	encodedRefreshToken := security.RawEncode([]byte(refreshToken))

	return &entities.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

func (useCases *UseCases) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	return useCases.usersService.GetUserByID(ctx, id)
}

func (useCases *UseCases) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	return useCases.usersService.GetUserByEmail(ctx, email)
}

func (useCases *UseCases) GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error) {
	return useCases.usersService.GetUsers(ctx, pagination)
}

func (useCases *UseCases) UpdateUserProfile(
	ctx context.Context,
	rawUserProfileData entities.RawUpdateUserProfileDTO,
) error {
	if rawUserProfileData.DisplayName != nil &&
		!validation.ValidateValueByRules(
			*rawUserProfileData.DisplayName,
			useCases.validationConfig.DisplayNameRegExps,
		) {
		return &customerrors.InvalidDisplayNameError{}
	}

	if rawUserProfileData.Phone != nil &&
		!validation.ValidateValueByRules(
			*rawUserProfileData.Phone,
			useCases.validationConfig.PhoneRegExps,
		) {
		return &customerrors.InvalidPhoneError{}
	}

	if rawUserProfileData.Telegram != nil &&
		!validation.ValidateValueByRules(
			*rawUserProfileData.Telegram,
			useCases.validationConfig.TelegramRegExps,
		) {
		return &customerrors.InvalidTelegramError{}
	}

	accessTokenPayload, err := security.ParseJWT(
		rawUserProfileData.AccessToken,
		useCases.securityConfig.JWT.SecretKey,
	)
	if err != nil {
		return &security.InvalidJWTError{}
	}

	floatUserID, ok := accessTokenPayload.(float64)
	if !ok {
		return &security.InvalidJWTError{}
	}

	userID := uint64(floatUserID)

	user, err := useCases.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	userProfileData := entities.UpdateUserProfileDTO{
		UserID:      user.ID,
		DisplayName: rawUserProfileData.DisplayName,
		Phone:       rawUserProfileData.Phone,
		Telegram:    rawUserProfileData.Telegram,
		Avatar:      rawUserProfileData.Avatar,
	}

	return useCases.usersService.UpdateUserProfile(ctx, userProfileData)
}

func (useCases *UseCases) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
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

func (useCases *UseCases) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*entities.TokensDTO, error) {
	// Decoding refresh token to get original JWT and compare its value with value in Database:
	oldRefreshTokenBytes, err := security.RawDecode(refreshToken)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving refresh token payload to get access token from refresh token:
	oldRefreshToken := string(oldRefreshTokenBytes)

	refreshTokenPayload, err := security.ParseJWT(
		oldRefreshToken,
		useCases.securityConfig.JWT.SecretKey,
	)
	if err != nil {
		return nil, &security.InvalidJWTError{}
	}

	oldAccessToken, ok := refreshTokenPayload.(string)
	if !ok {
		return nil, &security.InvalidJWTError{}
	}

	// Retrieving access token payload to get user ID:
	accessTokenPayload, err := security.ParseJWT(
		oldAccessToken,
		useCases.securityConfig.JWT.SecretKey,
		jwt.WithoutClaimsValidation(), // not validating claims due to expiration of JWT TTL
	)
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
	encodedRefreshToken := security.RawEncode([]byte(newRefreshToken))

	return &entities.TokensDTO{
		AccessToken:  newAccessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

func (useCases *UseCases) LogoutUser(ctx context.Context, accessToken string) error {
	accessTokenPayload, err := security.ParseJWT(accessToken, useCases.securityConfig.JWT.SecretKey)
	if err != nil {
		return &security.InvalidJWTError{}
	}

	floatUserID, ok := accessTokenPayload.(float64)
	if !ok {
		return &security.InvalidJWTError{}
	}

	userID := uint64(floatUserID)

	refreshToken, _ := useCases.authService.GetRefreshTokenByUserID(ctx, userID)
	if refreshToken == nil {
		return nil
	}

	return useCases.authService.ExpireRefreshToken(ctx, refreshToken.Value)
}

func (useCases *UseCases) VerifyUserEmail(ctx context.Context, verifyEmailToken string) error {
	strUserID, err := security.RawDecode(verifyEmailToken)
	if err != nil {
		return err
	}

	intUserID, err := strconv.Atoi(string(strUserID))
	if err != nil {
		return err
	}

	user, err := useCases.GetUserByID(ctx, uint64(intUserID))
	if err != nil {
		return err
	}

	if user.EmailConfirmed {
		return &customerrors.EmailAlreadyConfirmedError{}
	}

	return useCases.authService.VerifyUserEmail(ctx, user.ID)
}

func (useCases *UseCases) ForgetPassword(ctx context.Context, forgetPasswordToken, newPassword string) error {
	if !validation.ValidateValueByRules(newPassword, useCases.validationConfig.PasswordRegExps) {
		return &customerrors.InvalidPasswordError{}
	}

	strUserID, err := security.RawDecode(forgetPasswordToken)
	if err != nil {
		return err
	}

	intUserID, err := strconv.Atoi(string(strUserID))
	if err != nil {
		return err
	}

	user, err := useCases.GetUserByID(ctx, uint64(intUserID))
	if err != nil {
		return err
	}

	if security.ValidateHash(newPassword, user.Password) {
		return &customerrors.InvalidPasswordError{
			Message: "New password can not be equal to old password",
		}
	}

	hashedPassword, err := security.Hash(newPassword, useCases.securityConfig.HashCost)
	if err != nil {
		return err
	}

	return useCases.authService.ForgetPassword(ctx, user.ID, hashedPassword)
}

func (useCases *UseCases) ChangePassword(
	ctx context.Context,
	accessToken string,
	oldPassword string,
	newPassword string,
) error {
	if oldPassword == newPassword {
		return &customerrors.InvalidPasswordError{
			Message: "New password can not be equal to old password",
		}
	}

	if !validation.ValidateValueByRules(newPassword, useCases.validationConfig.PasswordRegExps) {
		return &customerrors.InvalidPasswordError{}
	}

	accessTokenPayload, err := security.ParseJWT(accessToken, useCases.securityConfig.JWT.SecretKey)
	if err != nil {
		return &security.InvalidJWTError{}
	}

	floatUserID, ok := accessTokenPayload.(float64)
	if !ok {
		return &security.InvalidJWTError{}
	}

	userID := uint64(floatUserID)

	user, err := useCases.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !security.ValidateHash(oldPassword, user.Password) {
		return &customerrors.WrongPasswordError{}
	}

	hashedPassword, err := security.Hash(newPassword, useCases.securityConfig.HashCost)
	if err != nil {
		return err
	}

	return useCases.authService.ChangePassword(ctx, userID, hashedPassword)
}

func (useCases *UseCases) SendVerifyEmailMessage(ctx context.Context, email string) error {
	user, err := useCases.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.EmailConfirmed {
		return &customerrors.EmailAlreadyConfirmedError{}
	}

	verifyEmailDTO := &notifications.VerifyEmailDTO{
		UserID: user.ID,
	}

	content, err := json.Marshal(verifyEmailDTO)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying to encode data for email verification for User with ID=%d",
				user.ID,
			),
			err,
		)

		return err
	}

	if err = useCases.natsPublisher.Publish(useCases.natsConfig.Subjects.VerifyEmail, content); err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying send verify-email message to User with ID=%d",
				user.ID,
			),
			err,
		)

		return err
	}

	return nil
}

func (useCases *UseCases) SendForgetPasswordMessage(ctx context.Context, email string) error {
	user, err := useCases.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if !user.EmailConfirmed {
		return &customerrors.EmailIsNotConfirmedError{}
	}

	forgetPasswordDTO := &notifications.ForgetPasswordDTO{
		UserID: user.ID,
	}

	content, err := json.Marshal(forgetPasswordDTO)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying to encode data for forget-password message for User with ID=%d",
				user.ID,
			),
			err,
		)

		return err
	}

	if err = useCases.natsPublisher.Publish(useCases.natsConfig.Subjects.ForgetPassword, content); err != nil {
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			fmt.Sprintf(
				"Error occurred while trying send forget-password message to User with ID=%d",
				user.ID,
			),
			err,
		)

		return err
	}

	return nil
}
