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

package pass

import (
	"fmt"
	"os"
	"strings"

	"github.com/raklaptudirm/krypt/pkg/dir"
)

func GetS(ident string, key []byte) (name string, pass *Password, err error) {
	passDir, err := dir.Pass()
	if err != nil {
		return
	}

	files, err := os.ReadDir(passDir)
	if err != nil {
		return
	}

	for _, file := range files {
		name = file.Name()
		if strings.HasPrefix(name, ident) {
			pass, err = get(ident, key)
			return
		}

		pass, err = get(name, key)
		if err == nil && pass.Name == ident {
			return
		}
	}

	err = fmt.Errorf("no password matched %v", ident)
	return
}
