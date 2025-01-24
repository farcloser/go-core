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
	"crypto/x509"
	"net"
	"net/http"

	"go.farcloser.world/core/log"
)

// Network holds network configuration for both client and server operations and provides helpers methods
// to retrieve TLSConfig and Transport objects.
type Network struct {
	clientConfig *Config
	serverConfig *Config
}

// TLSConfig returns a new tls.Config object populated against the configuration.
func (network *Network) TLSConfig() *tls.Config {
	cCA := x509.NewCertPool()
	if network.serverConfig.ClientCA != "" {
		ok := cCA.AppendCertsFromPEM([]byte(network.serverConfig.ClientCA))
		if !ok {
			log.Error().Msg("Invalid client CA in your config... Not loaded.")
		}
	}

	tlsMin := network.serverConfig.TLSMin
	if tlsMin < tls.VersionTLS12 {
		tlsMin = tls.VersionTLS13
	}

	tlsConfig := &tls.Config{ //nolint:gosec
		ClientCAs:  cCA,
		ClientAuth: tls.VerifyClientCertIfGiven,
		MinVersion: tlsMin,
		// XXX missing bits
		// VerifyPeerCertificate:
	}
	if network.serverConfig.ClientCertRequire {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConfig
}

// Transport returns a new Transport object populated against the configuration.
func (network *Network) Transport() *Transport {
	dialer := &net.Dialer{
		Timeout:   network.clientConfig.DialerTimeout,
		KeepAlive: network.clientConfig.DialerKeepAlive,
	}

	return &Transport{
		Transport: http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DialContext:         dialer.DialContext,
			TLSHandshakeTimeout: network.clientConfig.TLSHandshakeTimeout,
			TLSClientConfig:     network.getClientTLSConfig(),
		},
	}
}

func (network *Network) getClientTLSConfig() *tls.Config {
	var rootCAs *x509.CertPool
	if network.clientConfig.DisallowSystemRoot {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, _ = x509.SystemCertPool()
	}

	if network.clientConfig.RootCAs != nil {
		for _, v := range network.clientConfig.RootCAs {
			ok := rootCAs.AppendCertsFromPEM([]byte(v))
			if !ok {
				log.Error().Msg("Invalid root CA in your config... Not loaded.")
			}
		}
	}

	tlsMin := network.clientConfig.TLSMin
	if tlsMin < tls.VersionTLS12 {
		tlsMin = tls.VersionTLS13
	}

	tlsConfig := &tls.Config{ //nolint:gosec
		RootCAs:    rootCAs,
		MinVersion: tlsMin,
		// XXX missing bits
		// VerifyPeerCertificate:
	}

	return tlsConfig
}
