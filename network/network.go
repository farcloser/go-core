package network

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"

	"github.com/codecomet-io/go-core/log"
)

// Network is a struct that holds the network configuration and provides methods to retrieve TLSConfig and Transport objects
type Network struct {
	clientConfig *Config
	serverConfig *Config
}

// TLSConfig returns a new tls.Config object populated against the configuration
func (network *Network) TLSConfig() *tls.Config {
	cCA := x509.NewCertPool()
	if network.serverConfig.ClientCA != "" {
		ok := cCA.AppendCertsFromPEM([]byte(network.serverConfig.ClientCA))
		if !ok {
			log.Error().Msg("Invalid client CA in your config... Not loaded.")
		}
	}
	/*
		if serverConfig.ClientCA != nil {
			for _, v := range serverConfig.ClientCA {
				ok := ClientCA.AppendCertsFromPEM([]byte(v))
				if !ok {
					log.Error().Msg("Invalid client CA in your config... Not loaded.")
				}
			}
		}
	*/

	tlsConfig := &tls.Config{
		ClientCAs:  cCA,
		ClientAuth: tls.VerifyClientCertIfGiven,
		MinVersion: network.serverConfig.TLSMin,
		// XXX missing bits
		// VerifyPeerCertificate:
	}
	if network.serverConfig.ClientCertRequire {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConfig
}

// Transport returns a new Transport object populated against the configuration
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

	tlsConfig := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: network.clientConfig.TLSMin,
		// XXX missing bits
		// VerifyPeerCertificate:
	}

	return tlsConfig
}
