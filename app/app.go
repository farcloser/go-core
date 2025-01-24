/*
   Copyright Farcloser.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package app

// Use this for minimalistic apps that do not need configuration beyond Core,
// or take this as an example for your own app / config.
/*
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
		// Re-init logger with values
		log.Init(conf.Logger)
	}

	// Init network NOW before anything else - order matters!
	network.Init(conf.Client, conf.Server)

	// Init reporter
	if conf.Reporter != nil {
		reporter.Init(conf.Reporter)
	}

	// Init telemetry
	if conf.Telemetry != nil {
		telemetry.Init(conf.Telemetry)
	}

	return conf
}
*/
