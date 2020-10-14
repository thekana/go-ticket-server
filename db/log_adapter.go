package db

import (
	"context"

	"github.com/jackc/pgx/v4"

	"ticket-reservation/log"
)

type PostgresLogger struct {
	Logger log.Logger
}

func NewDatabaseLogger(logger *log.Logger) *PostgresLogger {
	return &PostgresLogger{Logger: *logger}
}

func (pglog *PostgresLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	// idea from https://github.com/jackc/pgx/blob/master/log/logrusadapter/adapter.go
	var logger = pglog.Logger
	if data != nil {
		logger = logger.WithFields(data)
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.WithFields(createFields("PGX_LOG_LEVEL", level)).Debugf(msg)
	case pgx.LogLevelDebug:
		logger.Debugf(msg)
	case pgx.LogLevelInfo:
		logger.Infof(msg)
	case pgx.LogLevelWarn:
		logger.Warnf(msg)
	case pgx.LogLevelError:
		logger.Errorf(msg)
	default:
		logger.WithFields(createFields("INVALID_PGX_LOG_LEVEL", level)).Errorf(msg)
	}
}

func createFields(key string, value interface{}) log.Fields {
	var fieldMap = make(map[string]interface{})
	fieldMap[key] = value
	return fieldMap
}
