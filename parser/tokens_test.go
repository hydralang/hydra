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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolString(t *testing.T) {
	a := assert.New(t)
	sym := Symbol{Name: "sym"}

	result := sym.String()

	a.Equal("sym", result)
}

func TestTokenStringBase(t *testing.T) {
	a := assert.New(t)
	sym := Symbol{Name: "sym"}
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	tok := Token{Sym: sym, Loc: loc}

	result := tok.String()

	a.Equal("file:3:2: <sym> token", result)
}

func TestTokenStringWithValue(t *testing.T) {
	a := assert.New(t)
	sym := Symbol{Name: "sym"}
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	tok := Token{Sym: sym, Loc: loc, Val: "value"}

	result := tok.String()

	a.Equal("file:3:2: <sym> token: value", result)
}
