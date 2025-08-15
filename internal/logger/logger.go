package logger

import (
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"
)

type PMLogger struct {
}

func (l *PMLogger) Info(msg string, fields ...process.LogFields) {
	log.Info().Fields(convertFields(fields)).Msg(msg)
}

func (l *PMLogger) Error(msg string, err error, fields ...process.LogFields) {
	log.Error().Err(err).Fields(convertFields(fields)).Msg(msg)
}

func convertFields(fields []process.LogFields) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	return fields[0]
}
