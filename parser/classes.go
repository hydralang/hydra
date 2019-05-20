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

package parser

import (
	"unicode"

	"github.com/hydralang/hydra/utils"
)

// Special character constants.
const (
	EOF rune = -(iota + 1) // End of file
	Err                    // An error occurred
)

// Defined character classes.
const (
	CharWS       uint16 = 1 << iota // Whitespace characters
	CharNL                          // Newline character
	CharBinDigit                    // Binary digit
	CharOctDigit                    // Octal digit
	CharDecDigit                    // Decimal digit
	CharHexDigit                    // Hexadecimal digit
	CharIDStart                     // Valid character for ID start
	CharIDCont                      // Valid character for ID continue
	CharStrFlag                     // String flag character
	CharQuote                       // String quote character
	CharComment                     // Comment character
)

// CharClasses is a mapping of character class flags to names.
var CharClasses = utils.FlagSet16{
	CharWS:       "whitespace",
	CharNL:       "newline",
	CharBinDigit: "binary digit",
	CharOctDigit: "octal digit",
	CharDecDigit: "decimal digit",
	CharHexDigit: "hexadecimal digit",
	CharIDStart:  "ID start",
	CharIDCont:   "ID continue",
	CharStrFlag:  "string flag",
	CharQuote:    "quote",
	CharComment:  "comment",
}

// digitData is a structure containing data about a particular digit.
type digitData struct {
	Class uint16 // The digit character class
	Val   int    // The integer value
}

// digits is a mapping of digit characters to the digit classes
// appropriate for them.
var digits = map[rune]digitData{
	'0': {CharBinDigit | CharOctDigit | CharDecDigit | CharHexDigit, 0},
	'1': {CharBinDigit | CharOctDigit | CharDecDigit | CharHexDigit, 1},
	'2': {CharOctDigit | CharDecDigit | CharHexDigit, 2},
	'3': {CharOctDigit | CharDecDigit | CharHexDigit, 3},
	'4': {CharOctDigit | CharDecDigit | CharHexDigit, 4},
	'5': {CharOctDigit | CharDecDigit | CharHexDigit, 5},
	'6': {CharOctDigit | CharDecDigit | CharHexDigit, 6},
	'7': {CharOctDigit | CharDecDigit | CharHexDigit, 7},
	'8': {CharDecDigit | CharHexDigit, 8},
	'9': {CharDecDigit | CharHexDigit, 9},
	'a': {CharHexDigit, 10},
	'A': {CharHexDigit, 10},
	'b': {CharHexDigit, 11},
	'B': {CharHexDigit, 11},
	'c': {CharHexDigit, 12},
	'C': {CharHexDigit, 12},
	'd': {CharHexDigit, 13},
	'D': {CharHexDigit, 13},
	'e': {CharHexDigit, 14},
	'E': {CharHexDigit, 14},
	'f': {CharHexDigit, 15},
	'F': {CharHexDigit, 15},
}

// AugChar is a struct that packages together a character, its class,
// its location, and any numeric value it may have.  This is the type
// that the scanner returns.
type AugChar struct {
	C     rune        // The character
	Class uint16      // The character's class
	Loc   Location    // The character's location
	Val   interface{} // The "value"; an integer for digits
}

// Classify classifies a character and composes an AugChar describing
// the character.
func (opts *Options) Classify(ch rune, loc Location, err error) AugChar {
	var class uint16
	var val interface{}

	// Handle the special characters
	if ch == EOF || ch == Err {
		return AugChar{ch, 0, loc, err}
	}

	// Start off with whitespace and newline
	if unicode.IsSpace(ch) {
		class |= CharWS
		if ch == '\n' {
			class |= CharNL
		}

		// Space is exclusive with everything else
		return AugChar{ch, class, loc, val}
	}

	// See if it's a digit
	if digData, ok := digits[ch]; ok {
		class |= digData.Class
		val = digData.Val
	}

	// Check for identifiers
	if opts.IDStart.Contains(ch) {
		class |= CharIDStart
	}
	if opts.IDCont.Contains(ch) {
		class |= CharIDCont
	}

	// Check for string flags and quotes
	if _, ok := opts.StrFlags[ch]; ok {
		class |= CharStrFlag
	}
	if _, ok := opts.Quotes[ch]; ok {
		class |= CharQuote
	}

	// Check for the comment character
	if ch == '#' {
		class |= CharComment
	}

	return AugChar{ch, class, loc, val}
}

// Advance advances the location to account for the specified
// character.
func (opts *Options) Advance(ch rune, loc *Location) {
	switch ch {
	case EOF, Err: // End of file
		loc.Advance(FilePos{})

	case '\n': // New line
		loc.Advance(FilePos{L: 1})

	case '\t': // Hit a tab
		loc.AdvanceTab(opts.TabStop)

	case '\f': // Don't count form feeds at the beginning of lines
		if loc.B.C > 1 {
			loc.Advance(FilePos{C: 1})
		}

	default: // Everything else advances by one column
		loc.Advance(FilePos{C: 1})
	}
}
