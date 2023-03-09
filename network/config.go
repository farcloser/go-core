package network

import "time"

// Config defines configuration to be applied to network communication, allowing to globally specify TLS certificates
// and minimum TLS version, timeouts, and other network properties.
// This should typically be marshalled from a local config file, and fed to network.Init.
type Config struct {
	// Common
	CertPath            string        `json:"certPath,omitempty"`
	KeyPath             string        `json:"keyPath,omitempty"`
	TLSMin              uint16        `json:"tlsMin,omitempty"`
	TLSHandshakeTimeout time.Duration `json:"tlsHandshakeTimeout,omitempty"`
	// Client only
	DialerTimeout      time.Duration `json:"dialerTimeout,omitempty"`
	DialerKeepAlive    time.Duration `json:"dialerKeepAlive,omitempty"`
	RootCAs            []string      `json:"rootCa,omitempty"`
	DisallowSystemRoot bool          `json:"disallowSystemRoot,omitempty"`
	// Server only
	ClientCA          string `json:"clientCa,omitempty"`
	ClientCertRequire bool   `json:"clientCertRequire,omitempty"`
	Port              uint16 `json:"port,omitempty"`

	Resolve func(pth ...string) string `json:"-"`
}
