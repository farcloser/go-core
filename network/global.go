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

package network

import (
	"crypto/tls"
	"net/http"

	"go.farcloser.world/core/log"
)

var network *Network //nolint:gochecknoglobals

// Init should be called when the app starts, from config objects.
func Init(clientConf, serverConf *Config) {
	log.Debug().Msg("Initializing network core with config")

	network = &Network{
		clientConfig: clientConf,
		serverConfig: serverConf,
	}

	http.DefaultTransport = network.Transport()
}

// GetTLSConfig returns the TLS configuration for the network.
func GetTLSConfig() *tls.Config {
	return network.TLSConfig()
}

// GetTransport returns the HTTP transport for the network.
func GetTransport() *Transport {
	return network.Transport()
}
