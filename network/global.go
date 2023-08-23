package network

import (
	"crypto/tls"
	"net/http"

	"go.farcloser.world/core/log"
)

var network *Network //nolint:gochecknoglobals

// Init should be called when the app starts, from config objects.
func Init(clientConf *Config, serverConf *Config) {
	log.Debug().Msg("Initializing network core with config")

	network = &Network{
		clientConfig: clientConf,
		serverConfig: serverConf,
	}

	http.DefaultTransport = network.Transport()
}

func GetTLSConfig() *tls.Config {
	return network.TLSConfig()
}

func GetTransport() *Transport {
	return network.Transport()
}
