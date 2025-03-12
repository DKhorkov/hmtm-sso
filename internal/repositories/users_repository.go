package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

const (
	selectAllColumns                = "*"
	usersTableName                  = "users"
	idColumnName                    = "id"
	userDisplayNameColumnName       = "display_name"
	userEmailColumnName             = "email"
	userEmailConfirmedColumnName    = "email_confirmed"
	userPasswordColumnName          = "password"
	userPhoneColumnName             = "phone"
	userPhoneConfirmedColumnName    = "phone_confirmed"
	userTelegramColumnName          = "telegram"
	userTelegramConfirmedColumnName = "telegram_confirmed"
	userAvatarColumnName            = "avatar"
)

func NewUsersRepository(
	dbConnector db.Connector,
	logger logging.Logger,
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
	logger        logging.Logger
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

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(usersTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
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

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(usersTableName).
		Where(sq.Eq{userEmailColumnName: email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
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

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(usersTableName).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := connection.QueryContext(
		ctx,
		stmt,
		params...,
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

	builder := sq.Update(usersTableName).Where(sq.Eq{idColumnName: userProfileData.UserID})
	if userProfileData.DisplayName != "" {
		builder = builder.Set(userDisplayNameColumnName, userProfileData.DisplayName)
	}

	if userProfileData.Phone != "" {
		builder = builder.Set(userPhoneColumnName, userProfileData.Phone)
	}

	if userProfileData.Telegram != "" {
		builder = builder.Set(userTelegramColumnName, userProfileData.Telegram)
	}

	if userProfileData.Avatar != "" {
		builder = builder.Set(userAvatarColumnName, userProfileData.Avatar)
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
