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
	"errors"
	"net"

	"go.farcloser.world/core/log"
)

/*
	FlagUp           Flags = 1 << iota // interface is administratively up
	FlagBroadcast                      // interface supports broadcast access capability
	FlagLoopback                       // interface is a loopback interface
	FlagPointToPoint                   // interface belongs to a point-to-point link
	FlagMulticast                      // interface supports multicast access capability
	FlagRunning                        // interface is in running state
*/

type (
	// Interface represents a mapping between network interface name.
	Interface = net.Interface
	// Address represents a network address, which can be an IP address or a Unix socket address.
	Address = net.Addr
)

// Interfaces is a struct that provides methods to retrieve network interfaces and their addresses.
type Interfaces struct{}

// GetAddresses retrieves all network addresses of the system, optionally filtering by IPv4.
//
//revive:disable:flag-parameter
func (*Interfaces) GetAddresses(onlyIPv4 bool, onlyName string) ([]Address, error) {
	list, err := net.Interfaces()
	if err != nil {
		return nil, errors.Join(ErrInterfacesRetrievalFailed, err)
	}

	var addresses []net.Addr

	for _, iface := range list {
		// If we want a specific name, make sure we have that
		if onlyName != "" && onlyName != iface.Name {
			continue
		}

		// Ignore interfaces that are down
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}

		// Ignore loopback interfaces
		if (iface.Flags & net.FlagLoopback) > 0 {
			continue
		}

		// Ignore ptp
		if (iface.Flags & net.FlagPointToPoint) > 0 {
			continue
		}

		// Get only multicast and broadcast enabled
		if (iface.Flags & net.FlagMulticast) == 0 {
			continue
		}

		if (iface.Flags & net.FlagBroadcast) == 0 {
			continue
		}

		var addrs []Address

		addrs, err = iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if !onlyIPv4 || addr.(*net.IPNet).IP.To4() != nil { //nolint:forcetypeassert
				log.Info().Str("iface name", iface.Name).
					Str("addr", addr.(*net.IPNet).String()). //nolint:forcetypeassert
					Msg("Found eligible interface")

				addresses = append(addresses, addr)
			}
		}
	}

	return addresses, nil
}
