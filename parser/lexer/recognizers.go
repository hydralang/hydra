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
	"github.com/hydralang/hydra/parser/common"
)

// Recognizer is a type describing recognizers.  A recognizer is
// initialized with the lexer object and implements the logic
// necessary to recognize a sequence of characters from the scanner.
//
// Note: some recognizers implement additional state; for instance,
// the string recognizer has state designed to interact with the
// recognizer for identifiers, to allow string flags to be recognized
// and processed.
type Recognizer interface {
	// Recognize is called to recognize a lexical construct.  Will
	// be called with the first character, and should push zero or
	// more tokens onto the lexer's tokens queue.
	Recognize(ch common.AugChar)
}

// RecogInit is a function that initializes a recognizer.  It will be
// passed the lexer object, and must return a Recognizer.
type RecogInit func(l *lexer) Recognizer
