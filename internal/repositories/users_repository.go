package repositories

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

func NewUsersRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *UsersRepository {
	return &UsersRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type UsersRepository struct {
	dbConnector   db.Connector
	logger        *slog.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *UsersRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM users AS u
			WHERE u.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM users AS u
			WHERE u.email = $1
		`,
		email,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UsersRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM users
		`,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			logging.LogErrorContext(
				ctx,
				repo.logger,
				"error during closing SQL rows",
				err,
			)
		}
	}()

	var users []entities.User
	for rows.Next() {
		user := entities.User{}
		columns := db.GetEntityColumns(&user) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UsersRepository) UpdateUserProfile(
	ctx context.Context,
	userProfileData entities.UpdateUserProfileDTO,
) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	// No fields to update:
	if userProfileData.DisplayName == "" &&
		userProfileData.Phone == "" &&
		userProfileData.Telegram == "" &&
		userProfileData.Avatar == "" {
		return nil
	}

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	builder := sq.Update("users").Where(sq.Eq{"id": userProfileData.UserID})
	if userProfileData.DisplayName != "" {
		builder = builder.Set("display_name", userProfileData.DisplayName)
	}

	if userProfileData.Phone != "" {
		builder = builder.Set("phone", userProfileData.Phone)
	}

	if userProfileData.Telegram != "" {
		builder = builder.Set("telegram", userProfileData.Telegram)
	}

	if userProfileData.Avatar != "" {
		builder = builder.Set("avatar", userProfileData.Avatar)
	}

	// pq postgres driver works only with $ placeholders:
	stmt, params, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
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

func (repo *UsersRepository) Close() error {
	return nil
}
