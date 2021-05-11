package server

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
var logDir string = path.Join(config.GetConfigHome(), "logs")

func configureLogger() func() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if DEBUGMODE {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = timeFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		logError(err)
	}

	logName := fmt.Sprintf("server-%s.log", getTime())
	logFile := getLogFile(logName)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return func() {
		logFile.Close()
	}
}

func logError(err error) {
	log.Error().Stack().Err(err).Msg("")
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

func getJobLogger(job *DownloadJob) zerolog.Logger {
	log.Info().Msgf("Created logger for JOB %s", job.ID)
	jobLogFile := getLogFile(fmt.Sprintf("job-%s.log", job.ID))
	// consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	// multi := zerolog.MultiLevelWriter(consoleWriter, jobLogFile)
	jobLogger := zerolog.New(jobLogFile).With().Timestamp().
		Str("job_id", job.ID).
		Logger()
	return jobLogger
}

func getLogFile(name string) *os.File {
	logPath := path.Join(logDir, name)
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logError(err)
	}
	return logFile
}

func getTime() string {
	return time.Now().Format(timeFormat)
}
