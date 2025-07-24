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

package telemetry

import "errors"

var (
	// ErrCloseError is returned when the telemetry provider fails to close properly.
	ErrCloseError = errors.New("close error")
	// ErrUnsupportedProviderType is returned when an unsupported provider type is used.
	ErrUnsupportedProviderType = errors.New("unsupported provider type")
	// ErrProviderCreationFailed is returned when the telemetry provider cannot be created.
	ErrProviderCreationFailed = errors.New("provider creation failed")
)
