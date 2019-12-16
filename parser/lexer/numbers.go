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

package lexer

import (
	"math/big"
	"strings"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/utils"
)

// Flags that define tracking data needed by the number recognizer.
const (
	NumInt   uint8 = 1 << iota // Number may be an integer
	NumFloat                   // Number may be a float
	NumWhole                   // Collecting the whole part of a float/int
	NumFract                   // Collecting the fraction
	NumExp                     // Collecting the exponent
	NumSign                    // Sign allowed next

	NumType  = NumInt | NumFloat            // Number type
	NumState = NumWhole | NumFract | NumExp // Number state
)

// NumFlags provides a mapping between number flags and the string
// describing them.
var NumFlags = utils.FlagSet8{
	NumInt:   "integer",
	NumFloat: "float",
	NumWhole: "whole state",
	NumFract: "fraction state",
	NumExp:   "exponent state",
	NumSign:  "sign allowed",
}

// baseMap maps the current base to the proper character
// classification flag.
var baseMap = map[int]uint16{
	2:  common.CharBinDigit,
	8:  common.CharOctDigit,
	10: common.CharDecDigit,
	16: common.CharHexDigit,
}

// baseFlagMap maps a flag character to the appropriate base
var baseFlagMap = map[rune]int{
	'b': 2,
	'B': 2,
	'o': 8,
	'O': 8,
	'x': 16,
	'X': 16,
}

// recognizeNumber is a recognizer for numbers.  It should be called
// when the character is a decimal digit, or when it is '.' followed
// by a decimal digit.
type recognizeNumber struct {
	l     *lexer           // The lexer
	buf   *strings.Builder // Buffer for numeric characters
	loc   utils.Location   // Location of 1st char
	flags uint8            // State tracking flags
	base  int              // Base to use interpreting number
}

// recogNumber constructs a recognizer for numbers.
func recogNumber(l *lexer) Recognizer {
	return &recognizeNumber{
		l:     l,
		buf:   &strings.Builder{},
		flags: NumInt | NumFloat | NumWhole,
	}
}

// Recognize is called to recognize a number.  Will be called with the
// first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeNumber) Recognize(ch common.AugChar) {
	// Initialize the state
	r.loc = ch.Loc

	// Interpret the first character
	r.buf.WriteRune(ch.C)
	if ch.C == '.' {
		// Has to be a float
		r.flags = NumFloat | NumFract
		r.base = 10
	} else if ch.C != '0' {
		// Has to be decimal
		r.base = 10
	}

	// Step through characters
	for ch = r.l.s.Next(); ; ch = r.l.s.Next() {
		// Check for flag character
		if r.base == 0 {
			// Is it a flag char?
			if base, ok := baseFlagMap[ch.C]; ok {
				r.flags &= NumInt | NumState
				r.base = base
				r.buf.Reset()
				continue
			}

			// Must be decimal
			r.base = 10
		}

		// The _ allows grouping digits; ignore it
		if ch.C == '_' {
			continue
		}

		// Handle float-specific characters
		if r.flags&NumFloat != 0 {
			if r.flags&NumWhole != 0 && ch.C == '.' {
				// Float, now collecting fraction
				r.flags = NumFloat | NumFract
				r.buf.WriteRune(ch.C)
				continue
			} else if r.flags&(NumWhole|NumFract) != 0 && (ch.C == 'e' || ch.C == 'E') {
				// Float, now collecting exponent
				r.flags = NumFloat | NumExp | NumSign
				r.buf.WriteRune(ch.C)
				continue
			} else if r.flags&NumSign != 0 {
				// No more signs
				r.flags &^= NumSign

				// Save it
				if ch.C == '+' || ch.C == '-' {
					r.buf.WriteRune(ch.C)
					continue
				}
			}
		}

		// Make sure it's a digit
		if ch.Class&baseMap[r.base] == 0 {
			break
		}

		r.buf.WriteRune(ch.C)
	}

	// Only terminated by operators and whitespace
	if ch.Class != 0 && ch.Class&common.CharWS == 0 {
		r.l.pushErr(ch.Loc, utils.ErrBadNumber)
		return
	}

	// Push the character back
	r.l.s.Push(ch)

	// Convert the buffer, preferring integer
	if r.flags&NumInt != 0 {
		// Convert number
		value := &big.Int{}
		value.SetString(r.buf.String(), r.base) // can't error

		r.l.pushTok(common.TokInt, r.loc.Thru(ch.Loc), value)
	} else {
		// Convert number
		value := &big.Float{}
		value.SetString(r.buf.String()) // can't error

		r.l.pushTok(common.TokFloat, r.loc.Thru(ch.Loc), value)
	}
}
