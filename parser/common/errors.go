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
	"errors"
	"fmt"
)

// Various errors that may occur during parsing.
var (
	ErrSplitEntity       = errors.New("entity split across files")
	ErrBadRune           = errors.New("illegal UTF-8 encoding")
	ErrBadIndent         = errors.New("inconsistent indentation")
	ErrBadOp             = errors.New("bad operator character")
	ErrMixedIndent       = errors.New("mixed whitespace types in indent")
	ErrDanglingBackslash = errors.New("dangling backslash")
	ErrBadNumber         = errors.New("bad character for number literal")
	ErrBadEscape         = errors.New("bad escape sequence")
	ErrBadStrChar        = errors.New("invalid character for string")
	ErrUnclosedStr       = errors.New("unclosed string literal")
	ErrBadIdent          = errors.New("bad identifier character")
)

// ErrDanglingOpen generates an error for a dangling open operator
// with no corresponding close operator.
func ErrDanglingOpen(tok *Token) error {
	return fmt.Errorf("unexpected EOF; expected \"%s\"", tok.Sym.Close)
}
