package daemon

import (
	"dcs/config"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var timeFormat = time.RFC3339

func configureLogger() func() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if DEBUGMODE {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = timeFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logDir := path.Join(config.GetConfigHome(), "logs")
	os.MkdirAll(logDir, 0755)
	logPath := path.Join(logDir, fmt.Sprintf("server-%s.log", time.Now().Format(timeFormat)))
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return func() {
		logFile.Close()
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_address", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Access")
		next.ServeHTTP(w, r)
	})
}

func getJobLogger(job DownloadJob) zerolog.Logger {
	jobLogger := log.With().Str("job_id", job.ID).Logger()
	return jobLogger
}
