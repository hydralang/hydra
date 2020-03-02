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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
)

func TestLexerSkipSpacesSpaces(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("     c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), 0)

	a.False(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestLexerSkipSpacesMixed(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("   \t c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), 0)

	a.True(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestLexerSkipSpacesLeadingFF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\f\f   c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), 0)

	a.True(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestLexerSkipSpacesSkipLeadingFF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\f\f   c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), SkipLeadFF)

	a.False(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestLexerSkipSpacesSkipLeadingMiddleFF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(" \f\f  c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), SkipLeadFF)

	a.True(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestLexerSkipSpacesNewline(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("  \n  c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), 0)

	a.False(mixed)
	next := s.Next()
	a.Equal('\n', next.C)
}

func TestLexerSkipSpacesNewlineSkipped(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("  \n  c"))
	s, _ := scanner.Scan(opts)
	l := &lexer{s: s}

	mixed := l.skipSpaces(s.Next(), SkipNL)

	a.True(mixed)
	next := s.Next()
	a.Equal('c', next.C)
}

func TestDoIndentSameColumn(t *testing.T) {
	a := assert.New(t)
	loc := common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)

	l.doIndent(1, loc)

	a.Equal(1, l.indent.Len())
	a.Equal(0, l.tokens.Len())
}

func TestDoIndentDeeperColumn(t *testing.T) {
	a := assert.New(t)
	loc := common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)

	l.doIndent(5, loc)

	a.Equal(2, l.indent.Len())
	elem := l.indent.Front()
	a.Equal(1, elem.Value.(int))
	elem = elem.Next()
	a.Equal(5, elem.Value.(int))
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokIndent,
		Loc: loc,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestDoIndentShallowerColumn(t *testing.T) {
	a := assert.New(t)
	loc := common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	l.indent.PushBack(9)
	l.indent.PushBack(13)

	l.doIndent(5, loc)

	a.Equal(2, l.indent.Len())
	elem := l.indent.Front()
	a.Equal(1, elem.Value.(int))
	elem = elem.Next()
	a.Equal(5, elem.Value.(int))
	a.Equal(2, l.tokens.Len())
	elem = l.tokens.Front()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
	}, elem.Value.(*common.Token))
}

func TestDoIndentShallowerColumnBadIndent(t *testing.T) {
	a := assert.New(t)
	loc := common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}
	l := &lexer{}
	l.indent.PushBack(1)
	l.indent.PushBack(5)
	l.indent.PushBack(9)
	l.indent.PushBack(13)

	l.doIndent(4, loc)

	a.Equal(1, l.indent.Len())
	a.Equal(1, l.indent.Front().Value.(int))
	a.Equal(4, l.tokens.Len())
	elem := l.tokens.Front()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokDedent,
		Loc: loc,
	}, elem.Value.(*common.Token))
	elem = elem.Next()
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: loc,
		Val: common.ErrBadIndent,
	}, elem.Value.(*common.Token))
}
