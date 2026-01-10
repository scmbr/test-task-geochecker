package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Logger()
}
func Info(msg string, fields ...map[string]interface{}) {
	event := log.Info()
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Error(msg string, err error, fields ...map[string]interface{}) {
	event := log.Error().Err(err).Caller(1)
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}
