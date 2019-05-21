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

import (
	"github.com/hydralang/hydra/utils"
)

// Defined string flags
const (
	StrRaw    uint8 = 1 << iota // Raw strings, ignores escapes
	StrBytes                    // Byte strings
	StrMulti                    // Multi-line (triple-quoted) string
	StrTriple                   // Quote allows triples
)

// StrFlags is a mapping of string flags to names.
var StrFlags = utils.FlagSet8{
	StrRaw:    "raw",
	StrBytes:  "bytes",
	StrMulti:  "multi-line",
	StrTriple: "triple quote",
}

// StrEscape is a function type for handling string escapes.  It is
// called with the character, the scanner, and the string flags, and
// should return a rune to add to the buffer.  If an error is
// returned, the error location should also be returned.
type StrEscape func(ch AugChar, s Scanner, flags uint8) (rune, Location, error)

// SimpleEscape sets up a StrEscape that returns a specified
// character.
func SimpleEscape(r rune) StrEscape {
	return func(ch AugChar, s Scanner, flags uint8) (rune, Location, error) {
		return r, Location{}, nil
	}
}

// HexEscape sets up a StrEscape that consumes the specified number of
// hexadecimal digits and returns the specified rune.
func HexEscape(cnt int) StrEscape {
	return func(ch AugChar, s Scanner, flags uint8) (rune, Location, error) {
		var r rune

		// Count off the specified number of characters
		for cnt--; cnt >= 0; cnt-- {
			ch = s.Next()
			if ch.C == Err {
				return 0, ch.Loc, ch.Val.(error)
			} else if ch.Class&CharHexDigit == 0 {
				return 0, ch.Loc, ErrBadEscape
			}

			// Accumulate the digit
			r |= rune(ch.Val.(int)) << (4 * uint(cnt))
		}

		// Return the rune
		return r, Location{}, nil
	}
}

// OctEscape is a StrEscape that consumes octal digits and returns the
// specified rune.
func OctEscape(ch AugChar, s Scanner, flags uint8) (rune, Location, error) {
	r := rune(ch.Val.(int))
	el := 1

	// 0x1f << 3 is still 8 bits
	for el < 3 && r <= 0x1f {
		ch = s.Next()
		if ch.C == Err {
			return 0, ch.Loc, ch.Val.(error)
		} else if ch.Class&CharOctDigit != 0 {
			// Another component of the code
			el++
			r = (r << 3) | rune(ch.Val.(int))
		} else {
			// Push it back
			s.Push(ch)
			break
		}
	}

	return r, Location{}, nil
}
