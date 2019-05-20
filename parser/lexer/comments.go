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

// recognizeComment is a recognizer for comments.  It should be called
// when the character is the comment character, '#'.
type recognizeComment struct {
	l *lexer // The lexer
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
	// XXX
}
