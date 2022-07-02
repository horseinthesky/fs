package main

import (
	"flag"
	"log"
	"time"

	"fs/device"
	"fs/internal"
	"fs/parser"
)

func main() {
	// Flags
	configFileName := flag.String("c", "config.yml", "Tool config filename")
	fileName := flag.String("f", "", "Flowdata filename. If not provided Web version is used")
	token := flag.String("t", "", "Netbox token. Must be provided if Web version is used")
	logFilePath := flag.String("l", "fs.log", "Log file name")
	logLevelName := flag.String("d", "info", "Log level. Available values: trace/debug/info/warn/error")
	flag.Parse()

	if *fileName == "" && *token == "" {
		log.Fatal("Error: data filename or Netbox token must be prodived")
	}

	// Setup logger
	logger, err := internal.SetupLogger(*&logFilePath, *logLevelName)
	if err != nil {
		log.Fatal(err)
	}

	// Get tool config
	cfg, err := internal.ParseConfig(*configFileName)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Get flow spec data
		var data *parser.FlowSpecData
		var err error

		if *fileName != "" {
			data, err = parser.GetDataFromFile(*fileName)
		} else {
			data, err = parser.GetDataFromWeb(cfg.Source, *token)
		}
		if err != nil {
			logger.Error().Err(err).Msg("failed to get flow spec routes data")
			time.Sleep(time.Second * time.Duration(cfg.Interval))
			continue
		}

		// Build XML config
		deviceConfig, err := parser.BuildConfig(data)
		if err != nil {
			logger.Error().Err(err).Msg("failed to build XML config")
			time.Sleep(time.Second * time.Duration(cfg.Interval))
			continue
		}
		logger.Trace().Str("parser", "config").Msg(deviceConfig)

		// Get inventory
		devices, err := device.GetDeviceList(cfg.Inventory)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get devices")
		}

		// Configure devices
		device.Configure(devices, deviceConfig, cfg.Creds, logger)

		time.Sleep(time.Second * time.Duration(cfg.Interval))
	}
}
