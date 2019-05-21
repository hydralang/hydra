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
	"io"
	"strings"
	"testing"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
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
	testStrFlags = map[rune]uint8{
		'r': common.StrRaw,
		'R': common.StrRaw,
		'b': common.StrBytes,
		'B': common.StrBytes,
	}
	testQuotes = map[rune]uint8{
		'"':  common.StrTriple,
		'\'': common.StrTriple,
	}
	testEscapes = map[rune]common.StrEscape{
		'\n': common.SimpleEscape(common.EOF),
		'0':  common.OctEscape,
		'1':  common.OctEscape,
		'2':  common.OctEscape,
		'3':  common.OctEscape,
		'4':  common.OctEscape,
		'5':  common.OctEscape,
		'6':  common.OctEscape,
		'7':  common.OctEscape,
		'\\': common.SimpleEscape('\\'),
		'a':  common.SimpleEscape('\a'),
		'b':  common.SimpleEscape('\b'),
		'e':  common.SimpleEscape('\x1b'),
		'f':  common.SimpleEscape('\f'),
		'n':  common.SimpleEscape('\n'),
		'r':  common.SimpleEscape('\r'),
		't':  common.SimpleEscape('\t'),
		'u':  common.HexEscape(4),
		'U':  common.HexEscape(8),
		'v':  common.SimpleEscape('\v'),
		'x':  common.HexEscape(2),
	}
	testProfile = &common.Profile{
		IDStart:  testIDStart,
		IDCont:   testIDCont,
		StrFlags: testStrFlags,
		Quotes:   testQuotes,
		Escapes:  testEscapes,
	}
)

func makeOptions(src io.Reader) *common.Options {
	return &common.Options{
		Source:   src,
		Filename: "file",
		Encoding: "utf-8",
		Prof:     testProfile,
		TabStop:  8,
	}
}

func TestLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*common.Lexer)(nil), &lexer{})
}

func TestLexWithScanner(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("test"))
	s, _ := scanner.Scan(opts)

	result, err := Lex(opts, s)

	a.NoError(err)
	a.NotNil(result)
	l, ok := result.(*lexer)
	a.True(ok)
	a.Equal(s, l.s)
	a.Equal(opts, l.opts)
	a.Equal(1, l.indent.Len())
	a.Equal(1, l.indent.Front().Value.(int))
	a.Equal(0, l.pair.Len())
	a.Equal(0, l.tokens.Len())
	a.Nil(l.prevTok)
}

func TestLexWithoutScanner(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("test"))

	result, err := Lex(opts, nil)

	a.NoError(err)
	a.NotNil(result)
	l, ok := result.(*lexer)
	a.True(ok)
	a.NotNil(l.s)
	a.Equal(opts, l.opts)
	a.Equal(1, l.indent.Len())
	a.Equal(1, l.indent.Front().Value.(int))
	a.Equal(0, l.pair.Len())
	a.Equal(0, l.tokens.Len())
	a.Nil(l.prevTok)
}

func TestLexScannerError(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("test"))
	opts.Encoding = ""

	result, err := Lex(opts, nil)

	a.Error(err)
	a.Nil(result)
}

func TestLexerNextEmpty(t *testing.T) {
	a := assert.New(t)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	l := &lexer{}
	l.indent.PushBack(1)

	result := l.Next()

	a.Nil(result)
	a.Equal(0, l.tokens.Len())
	a.Nil(l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextEnqueued(t *testing.T) {
	a := assert.New(t)
	tok1 := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: "tok1",
	}
	tok2 := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: "tok2",
	}
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	l := &lexer{}
	l.indent.PushBack(1)
	l.tokens.PushBack(tok1)
	l.tokens.PushBack(tok2)

	result := l.Next()

	a.Equal(tok1, result)
	a.Equal(1, l.tokens.Len())
	a.Equal(tok2, l.tokens.Front().Value.(*common.Token))
	a.Equal(tok1, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextError(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(common.AugChar{
		C:     common.Err,
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	})
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextEOF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(common.AugChar{
		C:     common.EOF,
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})
	expTok := &common.Token{
		Sym: common.TokEOF,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextEOFDangle(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	pairTok := &common.Token{
		Sym: &common.Symbol{
			Name:  "(",
			Close: ")",
		},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.pair.PushBack(pairTok)
	s.Push(common.AugChar{
		C:     common.EOF,
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: common.ErrDanglingOpen(pairTok),
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextComment(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rComment.On("Recognize", ch).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextDigit(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '5',
		Class: common.CharOctDigit | common.CharDecDigit | common.CharHexDigit | common.CharIDCont,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rNumber.On("Recognize", ch).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"number",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "number",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextPeriodDigit(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '.',
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch2 := common.AugChar{
		C:     '5',
		Class: common.CharOctDigit | common.CharDecDigit | common.CharHexDigit | common.CharIDCont,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	recs.rNumber.On("Recognize", ch1).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"number",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "number",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextIdent(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     'a',
		Class: common.CharHexDigit | common.CharIDStart | common.CharIDCont,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rIdent.On("Recognize", ch).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"ident",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "ident",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '"',
		Class: common.CharQuote,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rString.On("Recognize", ch).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"quote",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "quote",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextOp(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '!',
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rOp.On("Recognize", ch).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"op",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "op",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextPeriodOp(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '.',
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	recs.rOp.On("Recognize", ch1).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"op",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "op",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextOther(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '$',
		Class: common.CharIDCont,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: common.ErrBadOp,
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextContinuation(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '\\',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 3},
			E:    common.FilePos{L: 2, C: 4},
		},
	}
	ch2 := common.AugChar{
		C:     '\n',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 4},
			E:    common.FilePos{L: 3, C: 1},
		},
	}
	ch3 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	recs.rComment.On("Recognize", ch3).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextContinuationErr(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '\\',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 3},
			E:    common.FilePos{L: 2, C: 4},
		},
	}
	ch2 := common.AugChar{
		C:     common.Err,
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 4},
			E:    common.FilePos{L: 2, C: 4},
		},
		Val: assert.AnError,
	}
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 4},
			E:    common.FilePos{L: 2, C: 4},
		},
		Val: assert.AnError,
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextContinuationDangling(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '\\',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 3},
			E:    common.FilePos{L: 2, C: 4},
		},
	}
	ch2 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 4},
			E:    common.FilePos{L: 2, C: 4},
		},
	}
	l := &lexer{s: s}
	l.indent.PushBack(1)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 2, C: 4},
			E:    common.FilePos{L: 2, C: 4},
		},
		Val: common.ErrDanglingBackslash,
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextNewline(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch := common.AugChar{
		C:     '\n',
		Class: common.CharWS | common.CharNL,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 4, C: 1},
		},
	}
	l := &lexer{
		s: s,
		prevTok: &common.Token{
			Sym: common.TokEOF,
		},
	}
	l.indent.PushBack(1)
	s.Push(ch)
	expTok := &common.Token{
		Sym: common.TokNewline,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 4, C: 1},
		},
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespaceBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}
	ch5 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
	}
	recs.rComment.On("Recognize", ch5).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespaceLeadFFSkipped(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '\f',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 1},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	ch5 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}
	ch6 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
	}
	recs.rComment.On("Recognize", ch6).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	s.Push(ch6)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespaceMixed(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	ch1 := common.AugChar{
		C:     '\t',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 9},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 9},
			E:    common.FilePos{L: 3, C: 10},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 10},
			E:    common.FilePos{L: 3, C: 11},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 11},
			E:    common.FilePos{L: 3, C: 12},
		},
	}
	ch5 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 12},
			E:    common.FilePos{L: 3, C: 13},
		},
	}
	ch6 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 13},
			E:    common.FilePos{L: 3, C: 14},
		},
	}
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(13)
	s.Push(ch6)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 9},
		},
		Val: common.ErrMixedIndent,
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.Nil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespacePaired(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	pairTok := &common.Token{
		Sym: &common.Symbol{
			Name:  "(",
			Close: ")",
		},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	ch1 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}
	ch5 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
	}
	recs.rComment.On("Recognize", ch5).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	l.pair.PushBack(pairTok)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespacePairedLeadFFSkipped(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	pairTok := &common.Token{
		Sym: &common.Symbol{
			Name:  "(",
			Close: ")",
		},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	ch1 := common.AugChar{
		C:     '\f',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 1},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	ch5 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}
	ch6 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
	}
	recs.rComment.On("Recognize", ch6).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	l.pair.PushBack(pairTok)
	s.Push(ch6)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 5},
			E:    common.FilePos{L: 3, C: 6},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespacePairedMixed(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	pairTok := &common.Token{
		Sym: &common.Symbol{
			Name:  "(",
			Close: ")",
		},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	ch1 := common.AugChar{
		C:     '\t',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 9},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 9},
			E:    common.FilePos{L: 3, C: 10},
		},
	}
	ch3 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 10},
			E:    common.FilePos{L: 3, C: 11},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 11},
			E:    common.FilePos{L: 3, C: 12},
		},
	}
	ch5 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 12},
			E:    common.FilePos{L: 3, C: 13},
		},
	}
	ch6 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 13},
			E:    common.FilePos{L: 3, C: 14},
		},
	}
	recs.rComment.On("Recognize", ch6).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 13},
			E:    common.FilePos{L: 3, C: 14},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(13)
	l.pair.PushBack(pairTok)
	s.Push(ch6)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 13},
			E:    common.FilePos{L: 3, C: 14},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerNextWhitespacePairedNewline(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	recs := newMockRecs()
	oldRecs := recs.Install()
	defer oldRecs.Install()
	pairTok := &common.Token{
		Sym: &common.Symbol{
			Name:  "(",
			Close: ")",
		},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	ch1 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch2 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch3 := common.AugChar{
		C:     '\n',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 4, C: 1},
		},
	}
	ch4 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 4, C: 1},
			E:    common.FilePos{L: 4, C: 2},
		},
	}
	ch5 := common.AugChar{
		C:     ' ',
		Class: common.CharWS,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 4, C: 2},
			E:    common.FilePos{L: 4, C: 3},
		},
	}
	ch6 := common.AugChar{
		C:     '#',
		Class: common.CharComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 4, C: 3},
			E:    common.FilePos{L: 4, C: 4},
		},
	}
	recs.rComment.On("Recognize", ch6).Return(
		common.TokIdent,
		common.Location{
			File: "file",
			B:    common.FilePos{L: 4, C: 3},
			E:    common.FilePos{L: 4, C: 4},
		},
		"comment",
	)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	l.indent.PushBack(3)
	l.pair.PushBack(pairTok)
	s.Push(ch6)
	s.Push(ch5)
	s.Push(ch4)
	s.Push(ch3)
	s.Push(ch2)
	s.Push(ch1)
	expTok := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 4, C: 3},
			E:    common.FilePos{L: 4, C: 4},
		},
		Val: "comment",
	}

	result := l.Next()

	a.Equal(expTok, result)
	a.NotNil(l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(expTok, l.prevTok)
	recs.AssertExpectations(t)
}

func TestLexerPush(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}
	tok1 := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: "tok1",
	}
	tok2 := &common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: "tok2",
	}

	l.Push(tok1)
	l.Push(tok2)

	a.Equal(2, l.tokens.Len())
	elem := l.tokens.Front()
	a.Equal(tok2, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(tok1, elem.Value.(*common.Token))
}
