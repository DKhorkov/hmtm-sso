package repositories

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	sq "github.com/Masterminds/squirrel"

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
	DESC                            = "DESC"
	ASC                             = "ASC"
)

type UsersRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

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

func (repo *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
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

func (repo *UsersRepository) GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	builder := sq.
		Select(selectAllColumns).
		From(usersTableName).
		OrderBy(fmt.Sprintf("%s %s", idColumnName, DESC)).
		PlaceholderFormat(sq.Dollar)

	if pagination != nil && pagination.Limit != nil {
		builder = builder.Limit(*pagination.Limit)
	}

	if pagination != nil && pagination.Offset != nil {
		builder = builder.Offset(*pagination.Offset)
	}

	stmt, params, err := builder.ToSql()
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

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	builder := sq.
		Update(usersTableName).
		Where(sq.Eq{idColumnName: userProfileData.UserID}).
		Set(userPhoneColumnName, userProfileData.Phone).       // Update every time, because field is nullable
		Set(userTelegramColumnName, userProfileData.Telegram). // Update every time, because field is nullable
		PlaceholderFormat(sq.Dollar)                           // pq postgres driver works only with $ placeholders

	if userProfileData.DisplayName != nil {
		builder = builder.Set(userDisplayNameColumnName, userProfileData.DisplayName)
	}

	if userProfileData.Avatar != nil {
		builder = builder.Set(userAvatarColumnName, userProfileData.Avatar)
	}

	// If user deletes phone - we should update phone-confirmed field:
	if userProfileData.Phone == nil {
		builder = builder.Set(userPhoneConfirmedColumnName, false)
	}

	// If user deletes telegram - we should update telegram-confirmed field:
	if userProfileData.Telegram == nil {
		builder = builder.Set(userTelegramConfirmedColumnName, false)
	}

	stmt, params, err := builder.ToSql()
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
