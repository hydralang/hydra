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
	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/rangetable"
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
	testStrFlags = map[rune]interface{}{
		'r': nil,
		'R': nil,
		'b': nil,
		'B': nil,
	}
	testQuotes = map[rune]interface{}{
		'"':  nil,
		'\'': nil,
	}
)

var expected = map[rune]AugChar{
	' ': {
		C:     ' ',
		Class: CharWS,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'\n': {
		C:     '\n',
		Class: CharWS | CharNL,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'0': {
		C:     '0',
		Class: CharBinDigit | CharOctDigit | CharDecDigit | CharHexDigit | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 0,
	},
	'2': {
		C:     '2',
		Class: CharOctDigit | CharDecDigit | CharHexDigit | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 2,
	},
	'8': {
		C:     '8',
		Class: CharDecDigit | CharHexDigit | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 8,
	},
	'a': {
		C:     'a',
		Class: CharHexDigit | CharIDStart | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 10,
	},
	'A': {
		C:     'A',
		Class: CharHexDigit | CharIDStart | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 10,
	},
	'b': {
		C:     'b',
		Class: CharHexDigit | CharIDStart | CharIDCont | CharStrFlag,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 11,
	},
	'g': {
		C:     'g',
		Class: CharIDStart | CharIDCont,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'r': {
		C:     'r',
		Class: CharIDStart | CharIDCont | CharStrFlag,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'R': {
		C:     'R',
		Class: CharIDStart | CharIDCont | CharStrFlag,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'"': {
		C:     '"',
		Class: CharQuote,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
	'#': {
		C:     '#',
		Class: CharComment,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	},
}

func TestOptionsClassify(t *testing.T) {
	a := assert.New(t)
	opts := &Options{
		IDStart:  testIDStart,
		IDCont:   testIDCont,
		StrFlags: testStrFlags,
		Quotes:   testQuotes,
	}

	for r, exp := range expected {
		result := opts.Classify(r, Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		}, nil)

		a.Equal(exp, result)
	}
}

func TestOptionsClassifyEOF(t *testing.T) {
	a := assert.New(t)
	opts := &Options{
		IDStart:  testIDStart,
		IDCont:   testIDCont,
		StrFlags: testStrFlags,
		Quotes:   testQuotes,
	}

	result := opts.Classify(EOF, Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}, nil)

	a.Equal(AugChar{
		C:     EOF,
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: nil,
	}, result)
}

func TestOptionsClassifyErr(t *testing.T) {
	a := assert.New(t)
	opts := &Options{
		IDStart:  testIDStart,
		IDCont:   testIDCont,
		StrFlags: testStrFlags,
		Quotes:   testQuotes,
	}

	result := opts.Classify(Err, Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}, assert.AnError)

	a.Equal(AugChar{
		C:     Err,
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	}, result)
}

func TestOptionsAdvanceEOF(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance(EOF, &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 3},
	}, loc)
}

func TestOptionsAdvanceErr(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance(Err, &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 3},
	}, loc)
}

func TestOptionsAdvanceNewline(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance('\n', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 4, C: 1},
	}, loc)
}

func TestOptionsAdvanceTab8(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance('\t', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 9},
	}, loc)
}

func TestOptionsAdvanceTab4(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 4}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance('\t', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 5},
	}, loc)
}

func TestOptionsAdvanceFFBegin(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 2},
	}

	opts.Advance('\f', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 2},
	}, loc)
}

func TestOptionsAdvanceFFMiddle(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance('\f', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 4},
	}, loc)
}

func TestOptionsAdvanceOther(t *testing.T) {
	a := assert.New(t)
	opts := &Options{TabStop: 8}
	loc := Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}

	opts.Advance('o', &loc)

	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 4},
	}, loc)
}
