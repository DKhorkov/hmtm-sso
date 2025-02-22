package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

func NewAuthRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *AuthRepository {
	return &AuthRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type AuthRepository struct {
	dbConnector   db.Connector
	logger        *slog.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *AuthRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	var userID uint64
	err = connection.QueryRowContext(
		ctx,
		`
			INSERT INTO users (display_name, email, password) 
			VALUES ($1, $2, $3)
			RETURNING users.id
		`,
		userData.DisplayName,
		userData.Email,
		userData.Password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (repo *AuthRepository) CreateRefreshToken(
	ctx context.Context,
	userID uint64,
	refreshToken string,
	ttl time.Duration,
) (uint64, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	var refreshTokenID uint64
	err = connection.QueryRowContext(
		ctx,
		`
			INSERT INTO refresh_tokens (user_id, value, ttl) 
			VALUES ($1, $2, $3)
			RETURNING refresh_tokens.id
		`,
		userID,
		refreshToken,
		time.Now().UTC().Add(ttl),
	).Scan(&refreshTokenID)

	if err != nil {
		return 0, err
	}

	return refreshTokenID, nil
}

func (repo *AuthRepository) GetRefreshTokenByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.RefreshToken, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	refreshToken := &entities.RefreshToken{}
	columns := db.GetEntityColumns(refreshToken)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM refresh_tokens AS rt
			WHERE rt.user_id = $1
			  AND rt.ttl > CURRENT_TIMESTAMP
		`,
		userID,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (repo *AuthRepository) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	_, err = connection.ExecContext(
		ctx,
		`
			UPDATE refresh_tokens
			SET ttl = $1
			WHERE value = $2
		`,
		time.Now().UTC().Add(time.Hour*time.Duration(-24)),
		refreshToken,
	)

	return err
}

func (repo *AuthRepository) VerifyUserEmail(ctx context.Context, userID uint64) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	_, err = connection.ExecContext(
		ctx,
		`
			UPDATE users
			SET email_confirmed = true
			WHERE id = $1
		`,
		userID,
	)

	return err
}

func (repo *AuthRepository) ForgetPassword(ctx context.Context, userID uint64, newPassword string) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	transaction, err := repo.dbConnector.Transaction(ctx)
	if err != nil {
		return err
	}

	// Rollback transaction according Go best practises https://go.dev/doc/database/execute-transactions.
	defer func() {
		if err = transaction.Rollback(); err != nil {
			logging.LogErrorContext(ctx, repo.logger, "failed to rollback db transaction", err)
		}
	}()

	_, err = transaction.ExecContext(
		ctx,
		`
			UPDATE users
			SET password = $1
			WHERE id = $2
		`,
		newPassword,
		userID,
	)

	if err != nil {
		return err
	}

	// Getting refresh token for expiring:
	refreshToken := &entities.RefreshToken{}
	columns := db.GetEntityColumns(refreshToken)
	err = transaction.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM refresh_tokens AS rt
			WHERE rt.user_id = $1
			  AND rt.ttl > CURRENT_TIMESTAMP
		`,
		userID,
	).Scan(columns...)

	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return err
	case err == nil: // if active refresh token is found - expire it
		_, err = transaction.ExecContext(
			ctx,
			`
			UPDATE refresh_tokens
			SET ttl = $1
			WHERE value = $2
		`,
			time.Now().UTC().Add(time.Hour*time.Duration(-24)),
			refreshToken.Value,
		)

		if err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func (repo *AuthRepository) ChangePassword(ctx context.Context, userID uint64, newPassword string) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	_, err = connection.ExecContext(
		ctx,
		`
			UPDATE users
			SET password = $1
			WHERE id = $2
		`,
		newPassword,
		userID,
	)

	return err
}

func (repo *AuthRepository) Close() error {
	return nil
}
