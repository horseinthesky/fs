package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var levelMap = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
}

func SetupLogger(filePath *string, levelName string) (zerolog.Logger, error) {
	level, found := levelMap[levelName]
	if !found {
		return zerolog.Logger{}, fmt.Errorf("Error: level %s does not exist", levelName)
	}

	// Create parent dirs
	if strings.Contains(*filePath, "/") {
		logFileDirParts := strings.Split(*filePath, "/")
		logFileDir := strings.Join(logFileDirParts[:len(logFileDirParts)-1], "/")

		if _, err := os.Stat(*filePath); os.IsNotExist(err) {

			err = os.MkdirAll(logFileDir, 775)
			if err != nil {
				return zerolog.Logger{}, fmt.Errorf("Error: failed to create log file dir, error: %+w", err)
			}
		}
	}

	// Create log file
	logFile, err := os.OpenFile(*filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("Error: failed to open log file, error: %+w", err)
	}

	// Setup logger
	fileLogger := zerolog.New(logFile).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(level)

	return fileLogger, nil
}
