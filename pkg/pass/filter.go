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
	"math"
	"os"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/raklaptudirm/krypt/pkg/dir"
)

type FilterData int

const (
	FilterName FilterData = iota
	FilterUsername
	FilterPassStrength
)

type Filter struct {
	Type FilterData
	Data string
}

func Get(key []byte, filters ...Filter) (pass map[string]Password, err error) {
	pass = make(map[string]Password)

	passDir, err := dir.Pass()
	if err != nil {
		return
	}

	files, err := os.ReadDir(passDir)
	if err != nil {
		return
	}

	for _, file := range files {
		name := file.Name()

		password, err := get(name, key)
		if err != nil {
			return pass, err
		}

		if matchAll(*password, filters...) {
			pass[name] = *password
		}
	}

	if len(pass) == 0 {
		err = fmt.Errorf("no matching passwords")
	}

	return
}

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

func matchAll(pass Password, filters ...Filter) bool {
	for _, filter := range filters {
		if match(pass, filter) {
			return true
		}
	}

	return false
}

func match(pass Password, filter Filter) bool {
	minDist := 2

	switch filter.Type {
	case FilterName:
		filterName := strings.ToLower(filter.Data)
		passName := strings.ToLower(pass.Name)

		dist := levenshtein.ComputeDistance(filterName, passName)
		dist = int(math.Abs(float64(dist)))

		if dist <= minDist {
			return true
		}

		return false
	case FilterUsername:
		filterUser := strings.ToLower(filter.Data)
		passUserID := strings.ToLower(pass.UserID)

		dist := levenshtein.ComputeDistance(filterUser, passUserID)
		dist = int(math.Abs(float64(dist)))

		if dist <= minDist {
			return true
		}

		return false
	case FilterPassStrength:
		return false
	}

	// unreachable
	return false
}
