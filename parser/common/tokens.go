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
	"fmt"
	"strings"
)

// Symbol represents a defined symbol, or token type.  This could
// indicate something with a fixed value, like an operator, or
// something that has semantic value, such as a number literal.
type Symbol struct {
	Name string // The name of the symbol, for display purposes
}

// String constructs a string representation of a symbol--e.g., the
// symbol name.
func (s Symbol) String() string {
	return s.Name
}

// Token represents a single token emitted by the lexer.
type Token struct {
	Sym Symbol      // The token type
	Loc Location    // The location range of the token
	Val interface{} // The semantic value of the token
}

// String constructs a string representation of a token.
func (t Token) String() string {
	text := strings.Builder{}

	// Add the prefix
	text.WriteString(fmt.Sprintf("%s: <%s> token", t.Loc, t.Sym))

	// Add the semantic value, if present
	if t.Val != nil {
		text.WriteString(fmt.Sprintf(": %v", t.Val))
	}

	return text.String()
}

// Standard token symbols
var (
	TokError   = Symbol{"Error"}
	TokEOF     = Symbol{"EOF"}
	TokNewline = Symbol{"Newline"}
	TokIndent  = Symbol{"Indent"}
	TokDedent  = Symbol{"Dedent"}
	TokIdent   = Symbol{"Ident"}
	TokInt     = Symbol{"Int"}
	TokFloat   = Symbol{"Float"}
	TokString  = Symbol{"String"}
	TokBytes   = Symbol{"Bytes"}
)
