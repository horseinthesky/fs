package device

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"

	"fs/internal"
)

var wg sync.WaitGroup

func configureDevice(deviceName string, config string, creds internal.Creds, logger zerolog.Logger) {
	defer wg.Done()

	var authOption base.Option
	if creds.Key != "" {
		authOption = base.WithAuthPrivateKey(creds.Key)
	} else {
		authOption = base.WithAuthPassword(creds.Password)
	}

	d, _ := netconf.NewNetconfDriver(
		deviceName,
		base.WithPort(22),
		base.WithTimeoutTransport(5*time.Second),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername(creds.Username),
		authOption,
	)

	err := d.Open()
	if err != nil {
		logger.Error().Err(err).
			Str("device", deviceName).Str("result", "failure").
			Msg("failed to establish connection to " + deviceName)
		return
	}
	defer d.Close()

	r, err := d.EditConfig("candidate", string(config))
	if err != nil {
		logger.Error().Err(err).
			Str("device", deviceName).Str("result", "failure").
			Msg("failed to configure")
		return
	}

	logger.Info().
		Str("device", deviceName).Str("stage", "configure").Str("result", "success").
		Msg(r.Result)

	c, _ := d.Commit()
	if err != nil {
		logger.Error().Err(err).
			Str("device", deviceName).Str("result", "failure").
			Msg("failed to commit")
		return
	}

	logger.Info().
		Str("device", deviceName).Str("stage", "commit").Str("result", "success").
		Msg(c.Result)
}

func Configure(devices []*Device, config string, creds internal.Creds, logger zerolog.Logger) {
	for _, device := range devices {
		wg.Add(1)
		go configureDevice(device.Name, config, creds, logger)
	}

	wg.Wait()
}
