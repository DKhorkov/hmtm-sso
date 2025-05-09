package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Общая функция для тестирования ошибок
func testError(t *testing.T, err interface {
	Error() string
	Unwrap() error
}, defaultMessage, customMessage string, baseErr error) {
	t.Run("default message, no base error", func(t *testing.T) {
		e := err
		switch v := e.(type) {
		case *WrongPasswordError:
			*v = WrongPasswordError{}
		case *UserAlreadyExistsError:
			*v = UserAlreadyExistsError{}
		case *RefreshTokenAlreadyExistsError:
			*v = RefreshTokenAlreadyExistsError{}
		case *RefreshTokenNotFoundError:
			*v = RefreshTokenNotFoundError{}
		case *InvalidEmailError:
			*v = InvalidEmailError{}
		case *InvalidPasswordError:
			*v = InvalidPasswordError{}
		case *InvalidDisplayNameError:
			*v = InvalidDisplayNameError{}
		case *EmailAlreadyConfirmedError:
			*v = EmailAlreadyConfirmedError{}
		case *EmailIsNotConfirmedError:
			*v = EmailIsNotConfirmedError{}
		case *InvalidPhoneError:
			*v = InvalidPhoneError{}
		case *InvalidTelegramError:
			*v = InvalidTelegramError{}
		}
		require.Equal(t, defaultMessage, e.Error())
		require.Nil(t, e.Unwrap())
	})

	t.Run("custom message, no base error", func(t *testing.T) {
		e := err
		switch v := e.(type) {
		case *WrongPasswordError:
			*v = WrongPasswordError{Message: customMessage}
		case *UserAlreadyExistsError:
			*v = UserAlreadyExistsError{Message: customMessage}
		case *RefreshTokenAlreadyExistsError:
			*v = RefreshTokenAlreadyExistsError{Message: customMessage}
		case *RefreshTokenNotFoundError:
			*v = RefreshTokenNotFoundError{Message: customMessage}
		case *InvalidEmailError:
			*v = InvalidEmailError{Message: customMessage}
		case *InvalidPasswordError:
			*v = InvalidPasswordError{Message: customMessage}
		case *InvalidDisplayNameError:
			*v = InvalidDisplayNameError{Message: customMessage}
		case *EmailAlreadyConfirmedError:
			*v = EmailAlreadyConfirmedError{Message: customMessage}
		case *EmailIsNotConfirmedError:
			*v = EmailIsNotConfirmedError{Message: customMessage}
		case *InvalidPhoneError:
			*v = InvalidPhoneError{Message: customMessage}
		case *InvalidTelegramError:
			*v = InvalidTelegramError{Message: customMessage}
		}
		require.Equal(t, customMessage, e.Error())
		require.Nil(t, e.Unwrap())
	})

	t.Run("default message, with base error", func(t *testing.T) {
		e := err
		switch v := e.(type) {
		case *WrongPasswordError:
			*v = WrongPasswordError{BaseErr: baseErr}
		case *UserAlreadyExistsError:
			*v = UserAlreadyExistsError{BaseErr: baseErr}
		case *RefreshTokenAlreadyExistsError:
			*v = RefreshTokenAlreadyExistsError{BaseErr: baseErr}
		case *RefreshTokenNotFoundError:
			*v = RefreshTokenNotFoundError{BaseErr: baseErr}
		case *InvalidEmailError:
			*v = InvalidEmailError{BaseErr: baseErr}
		case *InvalidPasswordError:
			*v = InvalidPasswordError{BaseErr: baseErr}
		case *InvalidDisplayNameError:
			*v = InvalidDisplayNameError{BaseErr: baseErr}
		case *EmailAlreadyConfirmedError:
			*v = EmailAlreadyConfirmedError{BaseErr: baseErr}
		case *EmailIsNotConfirmedError:
			*v = EmailIsNotConfirmedError{BaseErr: baseErr}
		case *InvalidPhoneError:
			*v = InvalidPhoneError{BaseErr: baseErr}
		case *InvalidTelegramError:
			*v = InvalidTelegramError{BaseErr: baseErr}
		}
		expected := defaultMessage + ". Base error: " + baseErr.Error()
		require.Equal(t, expected, e.Error())
		require.Equal(t, baseErr, e.Unwrap())
	})
}

func TestErrors(t *testing.T) {
	baseErr := errors.New("underlying error")

	tests := []struct {
		name string
		err  interface {
			Error() string
			Unwrap() error
		}
		defaultMessage string
		customMessage  string
	}{
		{
			name:           "WrongPasswordError",
			err:            &WrongPasswordError{},
			defaultMessage: "wrong password",
			customMessage:  "incorrect password provided",
		},
		{
			name:           "UserAlreadyExistsError",
			err:            &UserAlreadyExistsError{},
			defaultMessage: "user with provided email already exists",
			customMessage:  "email test@example.com is taken",
		},
		{
			name:           "RefreshTokenAlreadyExistsError",
			err:            &RefreshTokenAlreadyExistsError{},
			defaultMessage: "refresh token already exists",
			customMessage:  "token for user ID=1 already issued",
		},
		{
			name:           "RefreshTokenNotFoundError",
			err:            &RefreshTokenNotFoundError{},
			defaultMessage: "refresh token not found",
			customMessage:  "token for user ID=1 not found",
		},
		{
			name:           "InvalidEmailError",
			err:            &InvalidEmailError{},
			defaultMessage: "email does not meet the requirements",
			customMessage:  "email format invalid",
		},
		{
			name:           "InvalidPasswordError",
			err:            &InvalidPasswordError{},
			defaultMessage: "password does not meet the requirements",
			customMessage:  "password too weak",
		},
		{
			name:           "InvalidDisplayNameError",
			err:            &InvalidDisplayNameError{},
			defaultMessage: "display name not meet the requirements",
			customMessage:  "display name too short",
		},
		{
			name:           "EmailAlreadyConfirmedError",
			err:            &EmailAlreadyConfirmedError{},
			defaultMessage: "provided email has been already confirmed",
			customMessage:  "email test@example.com already verified",
		},
		{
			name:           "EmailIsNotConfirmedError",
			err:            &EmailIsNotConfirmedError{},
			defaultMessage: "provided email is not confirmed",
			customMessage:  "email test@example.com not verified",
		},
		{
			name:           "InvalidPhoneError",
			err:            &InvalidPhoneError{},
			defaultMessage: "phone not meet the requirements",
			customMessage:  "phone number format invalid",
		},
		{
			name:           "InvalidTelegramError",
			err:            &InvalidTelegramError{},
			defaultMessage: "telegram not meet the requirements",
			customMessage:  "telegram handle invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testError(t, tc.err, tc.defaultMessage, tc.customMessage, baseErr)

			// Проверка, что ошибка реализует интерфейс error
			_, ok := tc.err.(error)
			require.True(t, ok, tc.name+" should implement error interface")
		})
	}
}
