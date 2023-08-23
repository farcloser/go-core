package config

import (
	"crypto/tls"
	"time"

	"go.farcloser.world/core/log"
)

const (
	defaultLogLevel            = log.InfoLevel
	defaultTLSClientMinVersion = tls.VersionTLS12
	defaultTLSServerMinVersion = tls.VersionTLS13
	defaultDialerKeepAlive     = 30 * time.Second
	defaultDialerTimeout       = 30 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
	defaultCertPath            = "x509.crt"
	defaultKeyPath             = "x509.key"
)
