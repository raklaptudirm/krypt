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

// Manager represents a password manager.
type Manager interface {
	// Password fetches the password data with the given checksum.
	Password([]byte) ([]byte, error)
	// Passwords fetches all the password data from the manager.
	Passwords() ([][]byte, error)
	// Write writes password data to the manager.
	Write([]byte) error
	// Delete deletes data with the given checksum from the manager.
	Delete([]byte) error
}
