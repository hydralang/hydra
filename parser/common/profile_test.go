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
	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/rangetable"

	"github.com/hydralang/hydra/testutils"
)

var (
	testIDStart = runes.In(rangetable.New(
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', '_', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
		'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
		'V', 'W', 'X', 'Y', 'Z',
	))
	testIDCont = runes.In(rangetable.New(
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', '_', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
		'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
		'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6',
		'7', '8', '9',
	))
	testStrFlags = map[rune]uint8{
		'r': StrRaw,
		'R': StrRaw,
		'b': StrBytes,
		'B': StrBytes,
	}
	testQuotes = map[rune]uint8{
		'"':  StrTriple,
		'\'': StrTriple,
	}
	testEscapes = map[rune]StrEscape{
		'\n': SimpleEscape(EOF),
		'0':  OctEscape,
		'1':  OctEscape,
		'2':  OctEscape,
		'3':  OctEscape,
		'4':  OctEscape,
		'5':  OctEscape,
		'6':  OctEscape,
		'7':  OctEscape,
		'\\': SimpleEscape('\\'),
		'a':  SimpleEscape('\a'),
		'b':  SimpleEscape('\b'),
		'e':  SimpleEscape('\x1b'),
		'f':  SimpleEscape('\f'),
		'n':  SimpleEscape('\n'),
		'r':  SimpleEscape('\r'),
		't':  SimpleEscape('\t'),
		'u':  HexEscape(4),
		'U':  HexEscape(8),
		'v':  SimpleEscape('\v'),
		'x':  HexEscape(2),
	}
	testKeywords = Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}
	testOperators = NewOperators(
		&Symbol{Name: "+"},
		&Symbol{Name: "-"},
		&Symbol{Name: "*"},
		&Symbol{Name: "/"},
		&Symbol{Name: "+="},
		&Symbol{Name: "-="},
		&Symbol{Name: "*="},
		&Symbol{Name: "/="},
		&Symbol{Name: "<"},
		&Symbol{Name: "<="},
		&Symbol{Name: ">"},
		&Symbol{Name: ">="},
		&Symbol{Name: "<<"},
		&Symbol{Name: ">>"},
		&Symbol{Name: "<<="},
		&Symbol{Name: ">>="},
		&Symbol{Name: "="},
		&Symbol{Name: "=="},
		&Symbol{Name: "!="},
		&Symbol{Name: "(", Close: ")"},
		&Symbol{Name: ")", Open: "("},
		&Symbol{Name: "[", Close: "]"},
		&Symbol{Name: "]", Open: "["},
		&Symbol{Name: "{", Close: "}"},
		&Symbol{Name: "}", Open: "{"},
	)
	testProfile = &Profile{
		IDStart:   testIDStart,
		IDCont:    testIDCont,
		StrFlags:  testStrFlags,
		Quotes:    testQuotes,
		Escapes:   testEscapes,
		Keywords:  testKeywords,
		Norm:      norm.NFKC,
		Operators: testOperators,
	}
)

func TestProfileCopy(t *testing.T) {
	a := assert.New(t)
	result := testProfile.Copy()

	testutils.AssertPtrNotEqual(a, testProfile, result)
	testutils.AssertPtrEqual(a, testProfile.IDStart, result.IDStart)
	testutils.AssertPtrEqual(a, testProfile.IDCont, result.IDCont)
	a.Equal(testProfile.StrFlags, result.StrFlags)
	a.Equal(testProfile.Quotes, result.Quotes)
	testutils.AssertPtrEqual(a, testProfile.Escapes, result.Escapes)
	a.Equal(testProfile.Keywords, result.Keywords)
	testutils.AssertPtrNotEqual(a, testProfile.Keywords, result.Keywords)
	a.Equal(testProfile.Norm, result.Norm)
	a.Equal(testOperators, result.Operators)
	testutils.AssertPtrNotEqual(a, testProfile.Operators, result.Operators)
}
