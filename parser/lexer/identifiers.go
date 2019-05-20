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

import "github.com/hydralang/hydra/parser/common"

// recognizeIdentifier is a recognizer for identifiers.  It should be
// called when the character is a legal identifier start character.
type recognizeIdentifier struct {
	l *lexer // The lexer
}

// recogIdentifier constructs a recognizer for identifiers.
func recogIdentifier(l *lexer) Recognizer {
	return &recognizeIdentifier{
		l: l,
	}
}

// Recognize is called to recognize a identifier.  Will be called with
// the first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeIdentifier) Recognize(ch common.AugChar) {
	// XXX
}
