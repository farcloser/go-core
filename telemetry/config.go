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

// traceEndpoint := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")
// PROMETHEUS ExporterType = "prometheus"
// OTLP       ExporterType = "otlp"

type ExporterType string

type Config struct {
	ServiceName string       `json:"serviceName"`
	Disabled    bool         `json:"disabled"`
	Type        ExporterType `json:"type"`

	// Only for jaegger it seems
	Endpoint string `json:"endpoint"`
}
