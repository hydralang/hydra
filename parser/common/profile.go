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

import "golang.org/x/text/runes"

// Profile describes a profile for the parser.  A profile is simply
// the version-specific rules, with desired options applied, and
// covers such things as the sets of identifier characters, etc.
type Profile struct {
	IDStart  runes.Set          // Set of valid identifier start chars
	IDCont   runes.Set          // Set of valid identifier continue chars
	StrFlags map[rune]uint8     // Valid string flags
	Quotes   map[rune]uint8     // Valid quote characters
	Escapes  map[rune]StrEscape // String escapes
}

// Copy generates a copy of a profile.  An Options structure always
// contains a profile copy, to enable it to be mutated by options
// without accidentally changing the master profile.
func (p *Profile) Copy() *Profile {
	return &Profile{
		IDStart:  p.IDStart,
		IDCont:   p.IDCont,
		StrFlags: p.StrFlags,
		Quotes:   p.Quotes,
		Escapes:  p.Escapes,
	}
}
