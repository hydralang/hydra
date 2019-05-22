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

	"github.com/hydralang/hydra/parser/common"
)

// recognizeIdentifier is a recognizer for identifiers.  It should be
// called when the character is a legal identifier start character.
type recognizeIdentifier struct {
	l   *lexer           // The lexer
	s   *recognizeString // String recognizer
	loc common.Location  // Location of first char
	buf *bytes.Buffer    // Buffer to accumulate identifier
}

// recogIdentifier constructs a recognizer for identifiers.
func recogIdentifier(l *lexer) Recognizer {
	return &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
}

// Recognize is called to recognize a identifier.  Will be called with
// the first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeIdentifier) Recognize(ch common.AugChar) {
	// Begin by saving the start location and initializing the
	// buffer
	r.loc = ch.Loc
	r.buf = &bytes.Buffer{}
	r.buf.WriteRune(ch.C)

	// Process for string flags
	r.s = r.s.setFlag(ch)

	// Process remaining identifier characters
	for ch = r.l.s.Next(); ; ch = r.l.s.Next() {
		// Handle errors
		if ch.C == common.Err {
			r.l.pushErr(ch.Loc, ch.Val.(error))
			return
		}

		// Process for string flags or quotes
		if r.s != nil {
			if ch.Class&common.CharQuote != 0 {
				r.s.Recognize(ch)
				return
			}

			r.s = r.s.setFlag(ch)
		}

		// Is it the end of the identifier?
		if ch.Class == 0 || ch.Class&common.CharWS != 0 {
			break
		} else if ch.Class&common.CharIDCont == 0 {
			// Bad character
			r.l.pushErr(ch.Loc, common.ErrBadIdent)
			return
		}

		// Add character to the buffer
		r.buf.WriteRune(ch.C)
	}

	// Push back last character retrieved
	r.l.s.Push(ch)

	// Get the identifier string
	ident := string(r.l.opts.Prof.Norm.Bytes(r.buf.Bytes()))

	// See if it's a keyword
	if sym, ok := r.l.opts.Prof.Keywords[ident]; ok {
		r.l.pushTok(sym, r.loc.Thru(ch.Loc), ident)
	} else {
		// Push the identifier
		r.l.pushTok(common.TokIdent, r.loc.Thru(ch.Loc), ident)
	}
}
