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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
	"github.com/hydralang/hydra/utils"
)

func TestLexerLastTokNil(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	result := l.lastTok()

	a.Nil(result)
}

func TestLexerLastTokReturned(t *testing.T) {
	a := assert.New(t)
	tok := &common.Token{
		Sym: common.TokIdent,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Val: "tok",
	}
	l := &lexer{
		prevTok: tok,
	}

	result := l.lastTok()

	a.Equal(result, tok)
}

func TestLexerLastTokQueued(t *testing.T) {
	a := assert.New(t)
	tok1 := &common.Token{
		Sym: common.TokIdent,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Val: "tok1",
	}
	tok2 := &common.Token{
		Sym: common.TokIdent,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Val: "tok2",
	}
	tok3 := &common.Token{
		Sym: common.TokIdent,
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Val: "tok3",
	}
	l := &lexer{
		prevTok: tok3,
	}
	l.tokens.PushBack(tok1)
	l.tokens.PushBack(tok2)

	result := l.lastTok()

	a.Equal(result, tok2)
}

func TestLexerPushTokBase(t *testing.T) {
	a := assert.New(t)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{
		prevTok: &common.Token{Sym: common.TokEOF},
	}
	l.indent.PushBack(1)

	result := l.pushTok(common.TokIdent, loc, "val")

	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: loc,
		Val: "val",
	}, result)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: loc,
		Val: "val",
	}, l.tokens.Front().Value.(*common.Token))
	a.Equal(1, l.indent.Len())
}

func TestLexerPushTokDuplicateNewline(t *testing.T) {
	a := assert.New(t)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{
		prevTok: &common.Token{Sym: common.TokNewline},
	}
	l.indent.PushBack(1)

	result := l.pushTok(common.TokNewline, loc, nil)

	a.Nil(result)
	a.Equal(0, l.tokens.Len())
	a.Equal(1, l.indent.Len())
}

func TestLexerPushTokInitialNewline(t *testing.T) {
	a := assert.New(t)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)

	result := l.pushTok(common.TokNewline, loc, nil)

	a.Nil(result)
	a.Equal(0, l.tokens.Len())
	a.Equal(1, l.indent.Len())
}

func TestLexerPushTokDedentEOF(t *testing.T) {
	a := assert.New(t)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)
	l.indent.PushBack(5)

	result := l.pushTok(common.TokEOF, loc, nil)

	a.Equal(&common.Token{
		Sym: common.TokEOF,
		Loc: loc,
	}, result)
	a.Equal(2, l.tokens.Len())
	elem := l.tokens.Front()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
		Val: nil,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokEOF,
		Loc: loc,
		Val: nil,
	}, elem.Value.(*common.Token))
	a.Equal(1, l.indent.Len())
}

func TestLexerPushTokIndent(t *testing.T) {
	a := assert.New(t)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)

	result := l.pushTok(common.TokIdent, loc, "val")

	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: loc,
		Val: "val",
	}, result)
	a.Equal(2, l.tokens.Len())
	elem := l.tokens.Front()
	a.Equal(&common.Token{
		Sym: common.TokIndent,
		Loc: loc,
		Val: nil,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: loc,
		Val: "val",
	}, elem.Value.(*common.Token))
	a.Equal(2, l.indent.Len())
	elem = l.indent.Front()
	a.Equal(1, elem.Value.(int))
	elem = elem.Next()
	a.Equal(2, elem.Value.(int))
}

func TestLexerPushErr(t *testing.T) {
	a := assert.New(t)
	opts := &common.Options{
		Encoding: "utf-8",
	}
	s, _ := scanner.Scan(opts)
	loc := utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}
	l := &lexer{s: s}

	l.pushErr(loc, assert.AnError)

	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: loc,
		Val: assert.AnError,
	}, l.tokens.Front().Value.(*common.Token))
	a.Nil(l.s)
}
