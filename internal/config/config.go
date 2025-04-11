package config

import (
	"fmt"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"github.com/DKhorkov/libs/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func New() Config {
	return Config{
		Environment: loadenv.GetEnv("ENVIRONMENT", "local"),
		Version:     loadenv.GetEnv("VERSION", "latest"),
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8070),
		},
		Security: security.Config{
			HashCost: loadenv.GetEnvAsInt("HASH_COST", 8), // Auth speed sensitive if large
			JWT: security.JWTConfig{
				RefreshTokenTTL: time.Hour * time.Duration(
					loadenv.GetEnvAsInt("REFRESH_TOKEN_JWT_TTL", 168),
				),
				AccessTokenTTL: time.Minute * time.Duration(
					loadenv.GetEnvAsInt("ACCESS_TOKEN_JWT_TTL", 15),
				),
				Algorithm: loadenv.GetEnv("JWT_ALGORITHM", "HS256"),
				SecretKey: loadenv.GetEnv("JWT_SECRET", "defaultSecret"),
			},
		},
		Database: db.Config{
			Host:         loadenv.GetEnv("POSTGRES_HOST", "0.0.0.0"),
			Port:         loadenv.GetEnvAsInt("POSTGRES_PORT", 5432),
			User:         loadenv.GetEnv("POSTGRES_USER", "postgres"),
			Password:     loadenv.GetEnv("POSTGRES_PASSWORD", "postgres"),
			DatabaseName: loadenv.GetEnv("POSTGRES_DB", "postgres"),
			SSLMode:      loadenv.GetEnv("POSTGRES_SSL_MODE", "disable"),
			Driver:       loadenv.GetEnv("POSTGRES_DRIVER", "postgres"),
			Pool: db.PoolConfig{
				MaxIdleConnections: loadenv.GetEnvAsInt("MAX_IDLE_CONNECTIONS", 1),
				MaxOpenConnections: loadenv.GetEnvAsInt("MAX_OPEN_CONNECTIONS", 1),
				MaxConnectionLifetime: time.Second * time.Duration(
					loadenv.GetEnvAsInt("MAX_CONNECTION_LIFETIME", 20),
				),
				MaxConnectionIdleTime: time.Second * time.Duration(
					loadenv.GetEnvAsInt("MAX_CONNECTION_IDLE_TIME", 10),
				),
			},
		},
		Logging: logging.Config{
			Level:       logging.Levels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().UTC().Format("02-01-2006")),
		},
		Validation: ValidationConfig{
			EmailRegExp: loadenv.GetEnv(
				"EMAIL_REGEXP",
				"^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$",
			),
			PasswordRegExps: loadenv.GetEnvAsSlice(
				"PASSWORD_REGEXPS",
				[]string{
					".{8,}",
					"[a-z]",
					"[A-Z]",
					"[0-9]",
					"[^\\d\\w]",
				},
				";",
			),
			DisplayNameRegExps: loadenv.GetEnvAsSlice(
				"DISPLAY_NAME_REGEXPS",
				[]string{
					"^.{4,70}$",
				},
				";",
			),
			PhoneRegExps: loadenv.GetEnvAsSlice(
				"PHONE_REGEXPS",
				[]string{
					"^(\\+7|8) ?(\\(495\\)|\\(499\\)|\\(812\\)|\\d{3}[ \\-]?) ?\\d{3}[ \\-]?\\d{2}[ \\-]?\\d{2}$",
				},
				";",
			),
			TelegramRegExps: loadenv.GetEnvAsSlice(
				"TELEGRAM_REGEXPS",
				[]string{
					"^@([a-zA-Z0-9_]){5,32}$",
				},
				";",
			),
		},
		Tracing: TracingConfig{
			Server: tracing.Config{
				ServiceName:    loadenv.GetEnv("TRACING_SERVICE_NAME", "hmtm-sso"),
				ServiceVersion: loadenv.GetEnv("VERSION", "latest"),
				JaegerURL: fmt.Sprintf(
					"http://%s:%d/api/traces",
					loadenv.GetEnv("TRACING_JAEGER_HOST", "0.0.0.0"),
					loadenv.GetEnvAsInt("TRACING_API_TRACES_PORT", 14268),
				),
			},
			Spans: SpansConfig{
				Root: tracing.SpanConfig{
					Opts: []trace.SpanStartOption{
						trace.WithAttributes(
							attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
						),
					},
					Events: tracing.SpanEventsConfig{
						Start: tracing.SpanEventConfig{
							Name: "Calling handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String(
										"Environment",
										loadenv.GetEnv("ENVIRONMENT", "local"),
									),
								),
							},
						},
						End: tracing.SpanEventConfig{
							Name: "Received response from handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String(
										"Environment",
										loadenv.GetEnv("ENVIRONMENT", "local"),
									),
								),
							},
						},
					},
				},
				Repositories: SpanRepositories{
					Auth: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
					Users: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
				},
			},
		},
		NATS: NATSConfig{
			ClientURL: fmt.Sprintf(
				"nats://%s:%d",
				loadenv.GetEnv("NATS_HOST", "0.0.0.0"),
				loadenv.GetEnvAsInt("NATS_CLIENT_PORT", 4222),
			),
			Subjects: NATSSubjects{
				VerifyEmail:    loadenv.GetEnv("NATS_VERIFY_EMAIL_SUBJECT", "verify-email"),
				ForgetPassword: loadenv.GetEnv("NATS_FORGET_PASSWORD_SUBJECT", "forget-password"),
			},
			Publisher: NATSPublisher{
				Name: loadenv.GetEnv("NATS_PUBLISHER_NAME", "hmtm-sso-publisher"),
			},
		},
	}
}

type HTTPConfig struct {
	Host string
	Port int
}

type ValidationConfig struct {
	EmailRegExp        string
	PasswordRegExps    []string // since Go's regex doesn't support backtracking.
	DisplayNameRegExps []string
	PhoneRegExps       []string
	TelegramRegExps    []string
}

type TracingConfig struct {
	Server tracing.Config
	Spans  SpansConfig
}

type SpansConfig struct {
	Root         tracing.SpanConfig
	Repositories SpanRepositories
}

type SpanRepositories struct {
	Auth  tracing.SpanConfig
	Users tracing.SpanConfig
}

type NATSConfig struct {
	ClientURL string
	Subjects  NATSSubjects
	Publisher NATSPublisher
}

type NATSSubjects struct {
	VerifyEmail    string
	ForgetPassword string
}

type NATSPublisher struct {
	Name string
}

type Config struct {
	HTTP        HTTPConfig
	Security    security.Config
	Database    db.Config
	Logging     logging.Config
	Validation  ValidationConfig
	Tracing     TracingConfig
	Environment string
	Version     string
	NATS        NATSConfig
}
