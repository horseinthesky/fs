package internal

import (
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var levelMap = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info": zerolog.InfoLevel,
	"warn": zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
}

func SetupLogger(fileName *string, levelName string) zerolog.Logger {
	level, found := levelMap[levelName]
	if !found {
		log.Fatalf("Level %s does not exist", levelName)
	}

	// Create parent dirs
	if strings.Contains(*fileName, "/") {
		logFileDirParts := strings.Split(*fileName, "/")
		logFileDir := strings.Join(logFileDirParts[:len(logFileDirParts)-1], "/")

		if _, err := os.Stat(*fileName); os.IsNotExist(err) {

			err = os.MkdirAll(logFileDir, 775)
			if err != nil {
				log.Fatalf("Failed to create log file dir, error: %+v\n", err)
			}
		}
	}

	// Create log file
	logFile, err := os.OpenFile(*fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file, error: %+v\n", err)
	}

	// Setup logger
	fileLogger := zerolog.New(logFile).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(level)

	return fileLogger
}
