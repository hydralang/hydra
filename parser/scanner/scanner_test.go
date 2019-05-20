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

package scanner

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func makeOptions(src io.Reader) *common.Options {
	return &common.Options{
		Source:   src,
		Filename: "file",
		Encoding: "utf-8",
		IDStart:  testIDStart,
		IDCont:   testIDCont,
		StrFlags: testStrFlags,
		Quotes:   testQuotes,
		TabStop:  8,
	}
}

func TestScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &scanner{})
}

func TestScanDefaultEncoding(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(bytes.NewReader([]byte{69, 108, 78, 105, 110, 204, 131, 111}))

	result, err := Scan(opts)

	a.NoError(err)
	a.NotNil(result)
	s, ok := result.(*scanner)
	a.True(ok)
	a.Equal(rune(utf8.RuneSelf), rune(s.buf[0]))
	a.Equal(0, s.pos)
	a.Equal(0, s.end)
	testutils.AssertPtrEqual(a, s.leUnknown, s.le)
	a.Equal(common.Err, s.pushed)
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 1},
	}, s.loc)
	buf := [20]byte{}
	n, err := s.source.Read(buf[:])
	a.NoError(err)
	a.Equal(8, n)
	a.Equal([]byte{69, 108, 78, 105, 110, 204, 131, 111}, buf[:n])
}

func TestScanISO8859_1(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(bytes.NewReader([]byte{69, 108, 78, 105, 241, 111}))
	opts.Encoding = "iso-8859-1"

	result, err := Scan(opts)

	a.NoError(err)
	a.NotNil(result)
	s, ok := result.(*scanner)
	a.True(ok)
	a.Equal(rune(utf8.RuneSelf), rune(s.buf[0]))
	a.Equal(0, s.pos)
	a.Equal(0, s.end)
	testutils.AssertPtrEqual(a, s.leUnknown, s.le)
	a.Equal(common.Err, s.pushed)
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 1},
	}, s.loc)
	buf := [20]byte{}
	n, err := s.source.Read(buf[:])
	a.NoError(err)
	a.Equal(7, n)
	a.Equal([]byte{69, 108, 78, 105, 195, 177, 111}, buf[:n])
}

func TestScanNoSuchEncoding(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(bytes.NewReader([]byte{69, 108, 78, 105, 110, 204, 131, 111}))
	opts.Encoding = "no-such-encoding"

	result, err := Scan(opts)

	a.NotNil(err)
	a.Nil(result)
}

func TestScannerNextCharBufferedASCII(t *testing.T) {
	a := assert.New(t)
	src := strings.NewReader("test")
	s := &scanner{
		source: src,
		end:    8,
	}
	copy(s.buf[0:], []byte{'b', 'u', 'f', 'f', 'e', 'r', 'e', 'd', utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('b', r)
	a.Equal(src, s.source)
	a.Equal([]byte("uffered"), s.buf[s.pos:s.end])
	a.Equal(1, s.pos)
	a.Equal(8, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharBufferedMultiByte(t *testing.T) {
	a := assert.New(t)
	src := strings.NewReader("test")
	s := &scanner{
		source: src,
		end:    5,
	}
	copy(s.buf[0:], []byte{195, 177, 'i', 'n', 'o', utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('\xf1', r)
	a.Equal(src, s.source)
	a.Equal([]byte("ino"), s.buf[s.pos:s.end])
	a.Equal(2, s.pos)
	a.Equal(5, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharBufferedBadChar(t *testing.T) {
	a := assert.New(t)
	src := strings.NewReader("test")
	s := &scanner{
		source: src,
		end:    5,
	}
	copy(s.buf[0:], []byte{128, 'n', 'i', 'n', 'o', utf8.RuneSelf})

	r, err := s.nextChar()

	a.Equal(common.ErrBadRune, err)
	a.Equal(common.Err, r)
	a.Nil(s.source)
	a.Equal([]byte("nino"), s.buf[s.pos:s.end])
	a.Equal(1, s.pos)
	a.Equal(5, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharEOF(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		source: nil,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal(common.EOF, r)
	a.Nil(s.source)
	a.Equal([]byte{}, s.buf[s.pos:s.end])
	a.Equal(0, s.pos)
	a.Equal(0, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharError(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		source: nil,
		err:    assert.AnError,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	r, err := s.nextChar()

	a.Equal(assert.AnError, err)
	a.Equal(common.Err, r)
	a.Nil(s.source)
	a.Equal([]byte{}, s.buf[s.pos:s.end])
	a.Equal(0, s.pos)
	a.Equal(0, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharEmpty(t *testing.T) {
	a := assert.New(t)
	src := strings.NewReader("test")
	s := &scanner{
		source: src,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('t', r)
	a.Equal(src, s.source)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Equal(1, s.pos)
	a.Equal(4, s.end)
	a.Nil(s.err)
}

func TestScannerNextCharSplitMulti(t *testing.T) {
	a := assert.New(t)
	src := strings.NewReader("\xb1ino")
	s := &scanner{
		source: src,
		end:    1,
	}
	copy(s.buf[0:], []byte{195, utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('\xf1', r)
	a.Equal(src, s.source)
	a.Equal([]byte("ino"), s.buf[s.pos:s.end])
	a.Equal(2, s.pos)
	a.Equal(5, s.end)
	a.Nil(s.err)
}

type mockReader struct {
	mock.Mock
}

func (r *mockReader) Read(b []byte) (int, error) {
	args := r.MethodCalled("Read")

	data := args.Get(0).([]byte)
	copy(b, data)

	return args.Int(1), args.Error(2)
}

func TestScannerNextCharEndOfReader(t *testing.T) {
	a := assert.New(t)
	src := &mockReader{}
	src.On("Read").Return([]byte("test"), 4, io.EOF)
	s := &scanner{
		source: src,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('t', r)
	a.Nil(s.source)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Equal(1, s.pos)
	a.Equal(4, s.end)
	a.Nil(s.err)
	src.AssertExpectations(t)
}

func TestScannerNextCharReadError(t *testing.T) {
	a := assert.New(t)
	src := &mockReader{}
	src.On("Read").Return([]byte("test"), 4, assert.AnError)
	s := &scanner{
		source: src,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	r, err := s.nextChar()

	a.NoError(err)
	a.Equal('t', r)
	a.Nil(s.source)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Equal(1, s.pos)
	a.Equal(4, s.end)
	a.Equal(assert.AnError, s.err)
	src.AssertExpectations(t)
}

func TestScannerPush(t *testing.T) {
	a := assert.New(t)
	s := &scanner{}
	ch := common.AugChar{
		C:     'c',
		Class: 5,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 1},
		},
		Val: nil,
	}

	s.Push(ch)

	a.Equal(1, s.queue.Len())
	a.Equal(ch, s.queue.Front().Value.(common.AugChar))
}

func TestScannerNextPushed(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		end:    4,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 1},
		},
	}
	copy(s.buf[0:], []byte{'t', 'e', 's', 't', utf8.RuneSelf})
	s.le = s.leNewline
	s.queue.PushFront(common.AugChar{
		C:     'p',
		Class: common.CharIDStart | common.CharIDCont,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 2, C: 1},
			E:    common.FilePos{L: 2, C: 2},
		},
		Val: nil,
	})
	ch := s.Next()

	a.Equal(common.AugChar{
		C:     'p',
		Class: common.CharIDStart | common.CharIDCont,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 2, C: 1},
			E:    common.FilePos{L: 2, C: 2},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte("test"), s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 1},
	}, s.loc)
}

func TestScannerNextLEPushed(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: 'p',
		end:    4,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 1},
		},
	}
	copy(s.buf[0:], []byte{'t', 'e', 's', 't', utf8.RuneSelf})
	s.le = s.leNewline

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     'p',
		Class: common.CharIDStart | common.CharIDCont,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte("test"), s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 2},
	}, s.loc)
}

func TestScannerNextEOF(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})
	s.le = s.leNewline

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     common.EOF,
		Class: 0,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte{}, s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 2},
		E:    common.FilePos{L: 1, C: 2},
	}, s.loc)
}

func TestScannerNextError(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		err:    assert.AnError,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})
	s.le = s.leNewline

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     common.Err,
		Class: 0,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: assert.AnError,
	}, ch)
	a.Equal([]byte{}, s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 2},
		E:    common.FilePos{L: 1, C: 2},
	}, s.loc)
}

func TestScannerNextCharacter(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		end:    4,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 1},
		},
	}
	copy(s.buf[0:], []byte{'t', 'e', 's', 't', utf8.RuneSelf})
	s.le = s.leNewline

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     't',
		Class: common.CharIDStart | common.CharIDCont,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 2},
	}, s.loc)
}

func TestScannerNextNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		end:    4,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	}
	copy(s.buf[0:], []byte{'\n', 'e', 's', 't', utf8.RuneSelf})
	s.le = s.leNewline

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     '\n',
		Class: common.CharWS | common.CharNL,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 2, C: 1},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 4},
		E:    common.FilePos{L: 2, C: 1},
	}, s.loc)
}

func leSwap(ch rune) rune {
	if ch == '\n' {
		return '\r'
	}
	return '\n'
}

func TestScannerNextCarriageSwapped(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		opts:   makeOptions(bytes.NewReader([]byte{})),
		pushed: common.Err,
		end:    4,
		loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	}
	copy(s.buf[0:], []byte{'\r', 'e', 's', 't', utf8.RuneSelf})
	s.le = leSwap

	ch := s.Next()

	a.Equal(common.AugChar{
		C:     '\n',
		Class: common.CharWS | common.CharNL,
		Loc: common.Location{
			File: "filename",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 2, C: 1},
		},
		Val: nil,
	}, ch)
	a.Equal([]byte("est"), s.buf[s.pos:s.end])
	a.Nil(s.err)
	a.Equal(common.Location{
		File: "filename",
		B:    common.FilePos{L: 1, C: 4},
		E:    common.FilePos{L: 2, C: 1},
	}, s.loc)
}
