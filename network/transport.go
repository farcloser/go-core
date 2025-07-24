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
	"fmt"
	"net/http"
	"strings"
)

// Transport implements http.Transport with a RoundTrip that has baked-in defaults, notably for GitHub
// It is not meant to be instantiated directly, but rather obtained through Get().Transport().
type Transport struct {
	http.Transport
	TokenValue string
	TokenType  string
}

// RoundTrip implements the http.RoundTripper interface.
func (adt *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.TokenValue != "" {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", adt.TokenType, adt.TokenValue))
	}

	if strings.HasSuffix(req.Host, "github.com") {
		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	resp, err := adt.Transport.RoundTrip(req)
	if err != nil {
		err = errors.Join(ErrRoundTrip, err)
	}

	return resp, err
}
