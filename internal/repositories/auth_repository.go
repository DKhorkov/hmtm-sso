package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

const (
	refreshTokensTableName      = "refresh_tokens"
	refreshTokenValueColumnName = "value"
	refreshTokenTTLColumnName   = "ttl"
	createdAtColumnName         = "created_at"
	updatedAtColumnName         = "updated_at"
	returningIDSuffix           = "RETURNING id"
	userIDColumnName            = "user_id"
)

func NewAuthRepository(
	dbConnector db.Connector,
	logger logging.Logger,
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
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *AuthRepository) RegisterUser(
	ctx context.Context,
	userData entities.RegisterUserDTO,
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

	stmt, params, err := sq.
		Insert(usersTableName).
		Columns(
			userDisplayNameColumnName,
			userEmailColumnName,
			userPasswordColumnName,
		).
		Values(
			userData.DisplayName,
			userData.Email,
			userData.Password,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var userID uint64
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&userID); err != nil {
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

	stmt, params, err := sq.
		Insert(refreshTokensTableName).
		Columns(
			userIDColumnName,
			refreshTokenValueColumnName,
			refreshTokenTTLColumnName,
		).
		Values(
			userID,
			refreshToken,
			time.Now().UTC().Add(ttl),
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var refreshTokenID uint64
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&refreshTokenID); err != nil {
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

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(refreshTokensTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		Where(
			sq.Expr(
				refreshTokenTTLColumnName + " > CURRENT_TIMESTAMP",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	refreshToken := &entities.RefreshToken{}

	columns := db.GetEntityColumns(refreshToken)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
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

	stmt, params, err := sq.
		Update(refreshTokensTableName).
		Where(sq.Eq{refreshTokenValueColumnName: refreshToken}).
		Set(
			refreshTokenTTLColumnName,
			time.Now().UTC().Add(time.Hour*time.Duration(-24)),
		).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(
		ctx,
		stmt,
		params...,
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

	stmt, params, err := sq.
		Update(usersTableName).
		Where(sq.Eq{idColumnName: userID}).
		Set(userEmailConfirmedColumnName, true).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(
		ctx,
		stmt,
		params...,
	)

	return err
}

func (repo *AuthRepository) ForgetPassword(
	ctx context.Context,
	userID uint64,
	newPassword string,
) error {
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

	stmt, params, err := sq.
		Update(usersTableName).
		Where(sq.Eq{idColumnName: userID}).
		Set(userPasswordColumnName, newPassword).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
		return err
	}

	stmt, params, err = sq.
		Select(selectAllColumns).
		From(refreshTokensTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		Where(
			sq.Expr(
				refreshTokenTTLColumnName + " > CURRENT_TIMESTAMP",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	// Getting refresh token for expiring:
	refreshToken := &entities.RefreshToken{}
	columns := db.GetEntityColumns(refreshToken)

	err = transaction.QueryRowContext(ctx, stmt, params...).Scan(columns...)

	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return err
	case err == nil: // if active refresh token is found - expire it
		stmt, params, err = sq.
			Update(refreshTokensTableName).
			Where(sq.Eq{refreshTokenValueColumnName: refreshToken.Value}).
			Set(
				refreshTokenTTLColumnName,
				time.Now().UTC().Add(time.Hour*time.Duration(-24)),
			).
			PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
			ToSql()
		if err != nil {
			return err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func (repo *AuthRepository) ChangePassword(
	ctx context.Context,
	userID uint64,
	newPassword string,
) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	stmt, params, err := sq.
		Update(usersTableName).
		Where(sq.Eq{idColumnName: userID}).
		Set(userPasswordColumnName, newPassword).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(
		ctx,
		stmt,
		params...,
	)

	return err
}
