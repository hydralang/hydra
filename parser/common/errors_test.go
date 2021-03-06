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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrDanglingOpen(t *testing.T) {
	a := assert.New(t)
	tok := &Token{Sym: &Symbol{Name: "(", Close: ")"}}

	result := ErrDanglingOpen(tok)

	a.EqualError(result, "unexpected EOF; expected \")\"")
}

func TestErrNoOpen(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: ")"}

	result := ErrNoOpen(sym)

	a.EqualError(result, "unexpected close operator \")\"")
}

func TestErrOpMismatch(t *testing.T) {
	a := assert.New(t)
	tok := &Token{
		Sym: &Symbol{Name: "["},
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
	}
	sym := &Symbol{Name: ")"}

	result := ErrOpMismatch(tok, sym)

	a.EqualError(result, "close operator \")\" does not match open operator \"[\" at file:3:2")
}
