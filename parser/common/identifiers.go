// Copyright (c) 2019 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package common

// Keywords is a map mapping identifier strings to the symbols to use
// for keyword tokens.
type Keywords map[string]*Symbol

// Copy produces a new copy of a Keywords object.
func (k Keywords) Copy() Keywords {
	// Construct a new map
	new := Keywords{}
	for text, sym := range k {
		new[text] = sym
	}

	return new
}

// Add adds a new keyword.
func (k Keywords) Add(sym *Symbol) {
	// Be idempotent
	if _, ok := k[sym.Name]; ok {
		return
	}

	k[sym.Name] = sym
}

// Remove removes a keyword.
func (k Keywords) Remove(sym *Symbol) {
	// Be idempotent
	if _, ok := k[sym.Name]; ok {
		delete(k, sym.Name)
	}
}
