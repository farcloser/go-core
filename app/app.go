package app

import (
	"go.codecomet.dev/core/config"
	"go.codecomet.dev/core/log"
	"go.codecomet.dev/core/network"
	"go.codecomet.dev/core/reporter"
	"go.codecomet.dev/core/telemetry"
)

func New(appName string, location ...string) *config.Core {
	// Create a new config object
	conf := config.New(appName, location...)

	// Init logger now so we get some acceptable output if shit hit the fan on loading the conf
	log.Init(conf.Logger)

	// Load configuration now if it exists
	if config.Exist(conf) {
		err := config.Load(conf)
		if err != nil {
			log.Fatal().Err(err).Msgf("Configuration file is invalid and needs to be fixed or removed: %s",
				conf.Resolve(conf.GetLocation()...))
		}
	}

	// Re-init logger with values
	log.Init(conf.Logger)

	// Init network NOW before anything else - order matters!
	network.Init(conf.Client, conf.Server)

	// Init reporter
	reporter.Init(conf.Reporter)

	// Init telemetry
	telemetry.Init(conf.Telemetry)

	return conf
}
