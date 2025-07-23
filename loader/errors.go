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

package loader

import "errors"

var (
	// ErrConfigLoadFail is returned when the configuration file cannot be loaded.
	ErrConfigLoadFail = errors.New("failed reading config file")
	// ErrConfigSaveFail is returned when the configuration file cannot be saved.
	ErrConfigSaveFail = errors.New("failed saving config file")
	// ErrConfigRemoveFail is returned when the configuration file cannot be removed.
	ErrConfigRemoveFail = errors.New("failed removing config file")
)
