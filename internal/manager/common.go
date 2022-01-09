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

package manager

import (
	"os"
	"path/filepath"
	"runtime"
)

// Pass in the internal password manager provided by krypt.
var Pass *pass

// Auth in the internal authentication manager provided by krypt.
var Auth *auth

func init() {
	dir := dataDir()
	os.Mkdir(dir, 0755)

	Pass = &pass{Dir: dir}
	Auth = &auth{Dir: dir}
}

// Data path precedence
// 1. XDG_DATA_HOME
// 2. LocalAppData (windows only)
// 3. HOME
func dataDir() string {
	a := os.Getenv("XDG_DATA_HOME")
	if a != "" {
		return filepath.Join(a, "krypt")
	}

	a = os.Getenv("LOCAL_APP_DATA")
	if runtime.GOOS == "windows" && a != "" {
		return filepath.Join(a, "Krypt")
	}

	a, _ = os.UserHomeDir()
	return filepath.Join(a, ".krypt")
}
