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
