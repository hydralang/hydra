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
	"bytes"
	"strings"
	"unicode"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/utils"
)

// buffer is a buffer for recogString.
type buffer interface {
	// putC writes a rune to a buffer.
	putC(r rune) error

	// get gets the fully built, buffered object.
	get() interface{}

	// sym returns the token symbol to use when pushing the token.
	sym() *common.Symbol
}

// bufString is a buffer for strings, composed of sequences of runes.
type bufString struct {
	strings.Builder
}

// putC writes a rune to a bufString.
func (b *bufString) putC(r rune) error {
	// Handle bad character
	if r > unicode.MaxRune {
		return utils.ErrBadStrChar
	}

	b.WriteRune(r)
	return nil
}

// get gets the fully built, buffered object.
func (b *bufString) get() interface{} {
	return b.String()
}

// sym returns the token symbol to use when pushing the token.
func (b bufString) sym() *common.Symbol {
	return common.TokString
}

// bufBytes is a buffer for bytes, composed of sequences of bytes.
type bufBytes struct {
	bytes.Buffer
}

// putC writes a rune to a bufBytes.
func (b *bufBytes) putC(r rune) error {
	// Handle bad character
	if r > 0xff {
		return utils.ErrBadStrChar
	}

	b.WriteByte(byte(r))
	return nil
}

// get gets the fully built, buffered object.
func (b *bufBytes) get() interface{} {
	return b.Bytes()
}

// sym returns the token symbol to use when pushing the token.
func (b bufBytes) sym() *common.Symbol {
	return common.TokBytes
}

// recognizeString is a recognizer for strings, consisting of groups
// of characters enclosed in, e.g., double quote ('"'), though the
// exact quote characters are dynamic.  It should be called when the
// character is a quote character.
type recognizeString struct {
	l      *lexer         // The lexer
	flags  uint8          // Recognized string flags
	loc    utils.Location // Location of first char
	buf    buffer         // Buffer for the string value
	q      rune           // The quote character to look for
	qcnt   int            // Number of quote characters spotted so far
	runLoc utils.Location // Beginning of last run of quote characters
}

// recogString constructs a recognizer for strings.
func recogString(l *lexer) Recognizer {
	return &recognizeString{
		l: l,
	}
}

// setFlag sets a string flag on the recognizer.  It returns the
// recognizer, or returns nil if the character is not a string flag or
// quote.
func (r *recognizeString) setFlag(ch common.AugChar) *recognizeString {
	// Save the character location
	if r.flags == 0 {
		r.loc = ch.Loc
	}

	// Set the flag
	switch {
	case ch.Class&common.CharStrFlag != 0:
		r.flags |= r.l.opts.Prof.StrFlags[ch.C]

	case ch.Class&common.CharQuote != 0:
		r.flags |= r.l.opts.Prof.Quotes[ch.C]

	default:
		return nil
	}

	return r
}

// escape handles an escape character encountered while processing a
// string.  Returns an error and a location if an error is
// encountered.
func (r *recognizeString) escape(ch common.AugChar) (utils.Location, error) {
	loc := ch.Loc

	// Handle raw string escapes
	if r.flags&common.StrRaw != 0 {
		// Write the backslash
		if err := r.buf.putC(ch.C); err != nil {
			return ch.Loc, err
		}

		// And the character following it
		ch = r.l.s.Next()
		if ch.C == common.Err {
			return ch.Loc, ch.Val.(error)
		} else if ch.C == common.EOF {
			return loc.ThruEnd(ch.Loc), utils.ErrUnclosedStr
		} else if err := r.buf.putC(ch.C); err != nil {
			return loc.ThruEnd(ch.Loc), err
		}

		return utils.Location{}, nil
	}

	// Check if we're escaping a quote
	ch = r.l.s.Next()
	if ch.Class&common.CharQuote != 0 {
		if err := r.buf.putC(ch.C); err != nil {
			return loc.ThruEnd(ch.Loc), err
		}
		return utils.Location{}, nil
	}

	// Make sure it's a valid escape
	if esc, ok := r.l.opts.Prof.Escapes[ch.C]; ok {
		// Process the escape
		c, eLoc, err := esc(ch, r.l.s, r.flags)
		if err != nil {
			return loc.ThruEnd(eLoc), err
		}

		// Write the escaped character; escapes return EOF to
		// indicate no character should be written, e.g.,
		// escaped newline
		if c != common.EOF {
			if err = r.buf.putC(c); err != nil {
				return loc.ThruEnd(eLoc), err
			}
		}

		return utils.Location{}, nil
	}

	// Not a valid escape
	return loc.ThruEnd(ch.Loc), utils.ErrBadEscape
}

// Recognize is called to recognize a string.  Will be called with the
// first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeString) Recognize(ch common.AugChar) {
	// Begin by applying flags for the quote character
	r.setFlag(ch)

	// Construct the buffer
	if r.flags&common.StrBytes == 0 {
		r.buf = &bufString{}
	} else {
		r.buf = &bufBytes{}
	}

	// Save the quote character
	r.q = ch.C

	// Look to see if we have triple quotes
	ch = r.l.s.Next()
	if r.flags&common.StrTriple != 0 {
		if ch.C == r.q {
			ch = r.l.s.Next()
			if ch.C == r.q {
				// Triple quote; mark for multi-line
				r.flags |= common.StrMulti

				// Get next character
				ch = r.l.s.Next()
			} else {
				// Empty string; push the character
				// back
				r.l.s.Push(ch)

				// Push a token
				r.l.pushTok(
					r.buf.sym(),
					r.loc.Thru(ch.Loc),
					r.buf.get(),
				)

				return
			}
		} else {
			// Not triple-quoted
			r.flags &^= common.StrTriple
		}
	}

	// Accumulate characters until the close
	for ; ; ch = r.l.s.Next() {
		// Handle quotes
		if ch.C == r.q {
			// Save the start of the run
			if r.qcnt == 0 {
				r.runLoc = ch.Loc
			}
			r.qcnt++

			// Does it close the string?
			if r.flags&common.StrTriple == 0 || r.qcnt >= 3 {
				break
			}
			continue
		}

		// Found quotes but didn't close the string?
		for ; r.qcnt > 0; r.qcnt-- {
			// Add quotes we skipped over before
			if err := r.buf.putC(r.q); err != nil {
				r.l.pushErr(r.runLoc, err)
				return
			}
		}

		// Handle the character
		switch ch.C {
		case common.Err: // Error occurred
			r.l.pushErr(ch.Loc, ch.Val.(error))
			return

		case common.EOF: // EOF in a string?
			r.l.pushErr(ch.Loc, utils.ErrUnclosedStr)
			return

		case '\\': // Introduces an escape
			if loc, err := r.escape(ch); err != nil {
				r.l.pushErr(loc, err)
				return
			}

		case '\n': // Newline, possible unclosed string
			if r.flags&common.StrMulti == 0 {
				r.l.pushErr(ch.Loc, utils.ErrUnclosedStr)
				return
			}
			fallthrough
		default: // Regular character
			if err := r.buf.putC(ch.C); err != nil {
				r.l.pushErr(ch.Loc, err)
				return
			}
		}
	}

	// Push a token
	r.l.pushTok(
		r.buf.sym(),
		r.loc.ThruEnd(ch.Loc),
		r.buf.get(),
	)
}
