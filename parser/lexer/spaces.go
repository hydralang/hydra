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
	"container/list"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/utils"
)

// Flags that may be given to skipSpaces.
const (
	SkipLeadFF uint8 = 1 << iota // Skip leading form feeds
	SkipNL                       // Skip newlines as well
)

// SkipFlags is a mapping of skip flags to names.
var SkipFlags = utils.FlagSet8{
	SkipLeadFF: "skip leading form feeds",
	SkipNL:     "skip newlines",
}

// skipSpaces skips whitespace for the lexer.  It returns a boolean
// indicating whether the whitespace was all one type, or whether it
// was mixed (e.g., spaces and tabs).  Flags allow leading form feeds
// to be ignored for the mixed-space calculation, and also can allow
// newlines to be skipped.
func (l *lexer) skipSpaces(ch common.AugChar, flags uint8) (mixed bool) {
	// Initialize the mixed space algorithm
	lastChar := ch.C
	mixed = false

	// Step through the whitespace
	for ; ch.Class&common.CharWS != 0; ch = l.s.Next() {
		// Skipping leading FF?
		if flags&SkipLeadFF != 0 {
			// Preemtively skip the form feed
			if ch.C == '\f' {
				continue
			}

			// Found last FF; restart mixed detection here
			lastChar = ch.C
			flags &^= SkipLeadFF
		}

		// Skipping newlines?
		if ch.Class&common.CharNL != 0 && flags&SkipNL == 0 {
			break
		}

		// Detect mixed whitespace
		if ch.C != lastChar {
			mixed = true
		}

		// Save this as the last char
		lastChar = ch.C
	}

	// This character is not whitespace, so push it back
	l.s.Push(ch)

	return
}

// doIndent is the core routine that manages indentation tracking.  It
// will push TokIndent and TokDedent tokens onto the token stack as
// appropriate, depending on the indentation level.
func (l *lexer) doIndent(col int, loc common.Location) {
	// Handle the simple cases first
	curCol := l.indent.Back().Value.(int)
	if col == curCol {
		// Same column
		return
	} else if col > curCol {
		// Deeper indentation
		l.pushTok(common.TokIndent, loc, nil)
		l.indent.PushBack(col)
		return
	}

	// Shallower indentation; produce one or more dedents to get
	// back to that point
	var elem *list.Element
	for elem = l.indent.Back(); elem.Value.(int) > col; elem = l.indent.Back() {
		l.pushTok(common.TokDedent, loc, nil)
		l.indent.Remove(elem)
	}

	// Produce an error if there's inconsistent indentation
	if elem.Value.(int) != col {
		l.pushErr(loc, common.ErrBadIndent)
	}
}
