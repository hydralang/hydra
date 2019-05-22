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

func TestSimpleEscape(t *testing.T) {
	a := assert.New(t)

	esc := SimpleEscape('c')
	r, loc, err := esc(AugChar{
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
	}, nil, 0)

	a.NoError(err)
	a.Equal('c', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 2},
	}, loc)
}

func TestHexEscapeBase(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '6',
		Class: CharHexDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 6,
	}).Once()
	s.On("Next").Return(AugChar{
		C:     '1',
		Class: CharHexDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 3},
			E:    FilePos{L: 3, C: 4},
		},
		Val: 1,
	}).Once()

	esc := HexEscape(2)
	r, loc, err := esc(AugChar{
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
	}, s, 0)

	a.NoError(err)
	a.Equal('a', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 4},
	}, loc)
	s.AssertExpectations(t)
}

func TestHexEscapeErr(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '6',
		Class: CharHexDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 6,
	}).Once()
	s.On("Next").Return(AugChar{
		C:     Err,
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 3},
			E:    FilePos{L: 3, C: 4},
		},
		Val: assert.AnError,
	}).Once()

	esc := HexEscape(2)
	r, loc, err := esc(AugChar{
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
	}, s, 0)

	a.Equal(assert.AnError, err)
	a.Equal(rune(0), r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 4},
	}, loc)
	s.AssertExpectations(t)
}

func TestHexEscapeBadEscape(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '6',
		Class: CharHexDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 6,
	}).Once()
	s.On("Next").Return(AugChar{
		C:     'z',
		Class: CharIDStart,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 3},
			E:    FilePos{L: 3, C: 4},
		},
	}).Once()

	esc := HexEscape(2)
	r, loc, err := esc(AugChar{
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
	}, s, 0)

	a.Equal(ErrBadEscape, err)
	a.Equal(rune(0), r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 3},
		E:    FilePos{L: 3, C: 4},
	}, loc)
	s.AssertExpectations(t)
}

func TestOctEscape1Digit(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	ch := AugChar{
		C:     'c',
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
	}
	s.On("Next").Return(ch).Once()
	s.On("Push", ch)

	r, loc, err := OctEscape(AugChar{
		C:     '1',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
		Val: 1,
	}, s, 0)

	a.NoError(err)
	a.Equal('\x01', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 2},
	}, loc)
	s.AssertExpectations(t)
}

func TestOctEscape2Digit(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '7',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 7,
	}).Once()
	ch := AugChar{
		C:     'c',
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 3},
			E:    FilePos{L: 3, C: 4},
		},
	}
	s.On("Next").Return(ch).Once()
	s.On("Push", ch)

	r, loc, err := OctEscape(AugChar{
		C:     '3',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
		Val: 3,
	}, s, 0)

	a.NoError(err)
	a.Equal('\x1f', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 3},
	}, loc)
	s.AssertExpectations(t)
}

func TestOctEscape3Digit(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '7',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 7,
	}).Once()
	s.On("Next").Return(AugChar{
		C:     '7',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 3},
			E:    FilePos{L: 3, C: 4},
		},
		Val: 7,
	}).Once()

	r, loc, err := OctEscape(AugChar{
		C:     '3',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
		Val: 3,
	}, s, 0)

	a.NoError(err)
	a.Equal('\xff', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 4},
	}, loc)
	s.AssertExpectations(t)
}

func TestOctEscape2DigitMax(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     '0',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: 0,
	}).Once()

	r, loc, err := OctEscape(AugChar{
		C:     '4',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
		Val: 4,
	}, s, 0)

	a.NoError(err)
	a.Equal('\x20', r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 1},
		E:    FilePos{L: 3, C: 3},
	}, loc)
	s.AssertExpectations(t)
}

func TestOctEscapeErr(t *testing.T) {
	a := assert.New(t)
	s := &mockScanner{}
	s.On("Next").Return(AugChar{
		C:     Err,
		Class: 0,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 2},
			E:    FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	}).Once()

	r, loc, err := OctEscape(AugChar{
		C:     '3',
		Class: CharOctDigit,
		Loc: Location{
			File: "file",
			B:    FilePos{L: 3, C: 1},
			E:    FilePos{L: 3, C: 2},
		},
		Val: 3,
	}, s, 0)

	a.Equal(assert.AnError, err)
	a.Equal(rune(0), r)
	a.Equal(Location{
		File: "file",
		B:    FilePos{L: 3, C: 2},
		E:    FilePos{L: 3, C: 3},
	}, loc)
	s.AssertExpectations(t)
}
