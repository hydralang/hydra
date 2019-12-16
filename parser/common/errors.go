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

	"github.com/hydralang/hydra/utils"
)

// ErrDanglingOpen generates an error for a dangling open operator
// with no corresponding close operator.
func ErrDanglingOpen(tok *Token) error {
	return fmt.Errorf("%w; expected \"%s\"", utils.ErrDanglingOpen, tok.Sym.Close)
}

// ErrNoOpen generates an error for a close operator with no
// corresponding open operator.
func ErrNoOpen(sym *Symbol) error {
	return fmt.Errorf("%w \"%s\"", utils.ErrNoOpen, sym.Name)
}

// ErrOpMismatch generates an error for a close operator that doesn't
// match the open operator.
func ErrOpMismatch(openTok *Token, close *Symbol) error {
	return fmt.Errorf("%w: operator \"%s\" does not match operator \"%s\" at %s", utils.ErrOpMismatch, close.Name, openTok.Sym.Name, openTok.Loc)
}
