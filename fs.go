package main

import (
	"flag"
	"log"
	"time"

	"github.com/horseinthesky/fs/device"
	"github.com/horseinthesky/fs/internal"
	"github.com/horseinthesky/fs/parser"
)

func main() {
	// Flags
	configFileName := flag.String("c", "config.yml", "Tool config filename")
	fileName := flag.String("f", "", "Flowdata config filename. If not provided Web version is used")
	token := flag.String("t", "", "Netbox token. Must be provided if Web version is used")
	logFileName := flag.String("l", "fs.log", "Log file name")
	logLevelName := flag.String("d", "info", "Log level. Available values: trace/debug/info/warn/error")
	flag.Parse()

	if *fileName == "" && *token == "" {
		log.Fatalf("Error: data filename or Netbox token must be prodived")
		return
	}

	// Setup logger
	fileLogger := internal.SetupLogger(*&logFileName, *logLevelName)

	// Get tool config
	cfg, err := internal.ParseConfig(*configFileName)
	if err != nil {
		log.Fatalf("Error reading config file: %+v\n", err)
		return
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
			fileLogger.Error().Err(err).Msg("failed to get flow spec routes data")
			time.Sleep(time.Second * time.Duration(cfg.Cooldown))
			continue
		}

		// Build XML config
		deviceConfig, err := parser.BuildConfig(data)
		if err != nil {
			fileLogger.Error().Err(err).Msg("failed to build XML config")
			time.Sleep(time.Second * time.Duration(cfg.Cooldown))
			continue
		}
		fileLogger.Trace().Str("parser", "config").Msg(deviceConfig)
		// fmt.Println(deviceConfig)

		// Get inventory
		devices, err := device.GetDeviceList(cfg.Inventory)
		if err != nil {
			fileLogger.Error().Err(err).Msg("failed to get devices")
		}

		// Configure devices
		device.Configure(devices, deviceConfig, cfg.Creds, fileLogger)

		time.Sleep(time.Second * time.Duration(cfg.Cooldown))
	}
}
