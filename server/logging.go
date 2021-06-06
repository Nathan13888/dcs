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

const timeFormat = time.RFC3339

var logDir string = path.Join(config.GetConfigHome(), "logs")
var logFile *os.File

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
	logFile = getLogFile(logName)

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

type loggingResponseWriter struct {
	http.ResponseWriter
	ResponseStatus int
	ResponseSize   int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseSize += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseStatus = statusCode
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		processedRequests++
		start := time.Now()
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			ResponseStatus: 200,
			ResponseSize:   0,
		}

		next.ServeHTTP(lw, r)
		log.Info().
			Str("method", r.Method).
			Int("status", lw.ResponseStatus).
			Int("size", lw.ResponseSize).
			Dur("duration", time.Since(start)).
			Str("path", r.URL.Path).
			Str("remote_address", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Access")
	})
}

// func getDBLogger() logger.Writer {
// 	zl := log.Logger.With().Str("db", config.DSN()).Logger()
// 	writer := logger.Writer{}
// 	return writer
// }

func getJobLogger(id string) zerolog.Logger {
	log.Info().Msgf("Created logger for JOB %s", id)
	jobLogFile := getLogFile(getJobLogName(id))
	// consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	// multi := zerolog.MultiLevelWriter(consoleWriter, jobLogFile)
	jobLogger := zerolog.New(jobLogFile).With().Timestamp().
		Str("job_id", id).
		Logger()
	return jobLogger
}

func getJobLogName(id string) string {
	return fmt.Sprintf("job-%s.log", id)
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
