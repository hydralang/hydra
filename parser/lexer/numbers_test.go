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
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
	"github.com/hydralang/hydra/utils"
)

func TestRecognizeNumberImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &recognizeNumber{})
}

func TestRecogNumber(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	result := recogNumber(l)

	r, ok := result.(*recognizeNumber)
	a.True(ok)
	a.Equal(l, r.l)
	a.NotNil(r.buf)
	a.Equal(NumInt|NumFloat|NumWhole, r.flags)
}

func TestRecogNumber0(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 2},
		},
		Val: &big.Int{},
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumber15(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("15"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 3},
		},
		Val: big.NewInt(15),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumberB10(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0b10"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 5},
		},
		Val: big.NewInt(2),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumberO15(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0o15"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 5},
		},
		Val: big.NewInt(13),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumberX15(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0x15"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 5},
		},
		Val: big.NewInt(21),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumber15underscore00(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("15_00"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 6},
		},
		Val: big.NewInt(1500),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecogNumber0point5(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0.5"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 4},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(0.5, flVal, 0.0001)
}

func TestRecogNumberPoint5(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(".5"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 3},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(0.5, flVal, 0.0001)
}

func TestRecogNumber1e2(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("1e2"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 4},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(1e2, flVal, 0.0001)
}

func TestRecogNumber1eMinus2(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("1e-2"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 5},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(1e-2, flVal, 0.0001)
}

func TestRecogNumber1ePlus2(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("1e+2"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 5},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(1e2, flVal, 0.0001)
}

func TestRecogNumberAllParts(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("1.51E+2"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokFloat, tok.Sym)
	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 1, C: 1},
		E:    utils.FilePos{L: 1, C: 8},
	}, tok.Loc)
	flVal, _ := tok.Val.(*big.Float).Float32()
	a.InEpsilon(1.51e2, flVal, 0.0001)
}

func TestRecogNumberWhitespace(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0 "))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 2},
		},
		Val: &big.Int{},
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal(' ', ch.C)
}

func TestRecogNumberOp(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0!"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokInt,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 1},
			E:    utils.FilePos{L: 1, C: 2},
		},
		Val: &big.Int{},
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal('!', ch.C)
}

func TestRecogNumberIdent(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("0a"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}
	l.indent.PushBack(1)
	r := recogNumber(l)

	r.Recognize(s.Next())

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 1, C: 2},
			E:    utils.FilePos{L: 1, C: 3},
		},
		Val: utils.ErrBadNumber,
	}, l.tokens.Front().Value.(*common.Token))
}
