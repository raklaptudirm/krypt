// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmdutil

import "fmt"

// Version represents the version of a krypt executable.
type Version struct {
	version   string // krypt version
	buildDate string // executable build date
}

// String converts a Version to a version string.
func (v *Version) String() string {
	if v.buildDate == "" {
		return v.version
	}

	// krypt version-date
	return fmt.Sprintf("krypt %v-%v\n", v.version, v.buildDate)
}

// NewVersion returns a *Version from the provided version and build
// date.
func NewVersion(v, d string) *Version {
	return &Version{
		version:   v,
		buildDate: d,
	}
}
