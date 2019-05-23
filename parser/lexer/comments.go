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
	"strings"

	"github.com/hydralang/hydra/parser/common"
)

// recognizeComment is a recognizer for comments.  It should be called
// when the character is the comment character, '#'.
type recognizeComment struct {
	l   *lexer           // The lexer
	loc common.Location  // Location of first char
	buf *strings.Builder // Buffer to accumulate doc comment
}

// recogComment constructs a recognizer for comments.
func recogComment(l *lexer) Recognizer {
	return &recognizeComment{
		l: l,
	}
}

// Recognize is called to recognize a comment.  Will be called with
// the first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeComment) Recognize(ch common.AugChar) {
	// Begin by saving the start location
	r.loc = ch.Loc

	// See if we have a doc comment
	next := r.l.s.Next()
	if next.C == ch.C {
		// Initialize the buffer
		r.buf = &strings.Builder{}
	} else {
		// Put it back for reprocessing
		r.l.s.Push(next)
	}

	// Skip through the comment
	for ch = r.l.s.Next(); ; ch = r.l.s.Next() {
		// Handle errors
		if ch.C == common.Err {
			r.l.pushErr(ch.Loc, ch.Val.(error))
			return
		}

		// Process up to newline or EOF
		if ch.C == common.EOF || ch.Class&common.CharNL != 0 {
			break
		}

		// Accumulate characters only if it's a doc comment
		if r.buf != nil {
			r.buf.WriteRune(ch.C)
		}
	}

	// Put the character back
	r.l.s.Push(ch)

	// Generate a doc comment token
	if r.buf != nil {
		r.l.pushTok(common.TokDocComment, r.loc.Thru(ch.Loc), r.buf.String())
	}
}
