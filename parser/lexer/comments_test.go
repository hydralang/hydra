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

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
	"github.com/stretchr/testify/assert"
)

func TestRecognizeCommentImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &recognizeComment{})
}

func TestRecogComment(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	result := recogComment(l)

	r, ok := result.(*recognizeComment)
	a.True(ok)
	a.Equal(l, r.l)
	a.Equal(common.Location{}, r.loc)
	a.Nil(r.buf)
}

func TestRecognizeCommentRecognizeBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("# this is a test\n"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeComment{l: l}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(0, l.tokens.Len())
	a.Equal(common.AugChar{
		C:     '\n',
		Class: common.CharWS | common.CharNL,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 17},
			E:    common.FilePos{L: 2, C: 1},
		},
	}, s.Next())
}

func TestRecognizeCommentRecognizeDocComment(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("##  this is a test\n"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeComment{l: l}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokDocComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 19},
		},
		Val: "  this is a test",
	}, l.tokens.Front().Value.(*common.Token))
	a.Equal(common.AugChar{
		C:     '\n',
		Class: common.CharWS | common.CharNL,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 19},
			E:    common.FilePos{L: 2, C: 1},
		},
	}, s.Next())
}

func TestRecognizeCommentRecognizeDocCommentEOF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("##  this is a test"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeComment{l: l}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokDocComment,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 19},
		},
		Val: "  this is a test",
	}, l.tokens.Front().Value.(*common.Token))
	a.Equal(common.AugChar{
		C:     common.EOF,
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 19},
			E:    common.FilePos{L: 1, C: 19},
		},
	}, s.Next())
}

func TestRecognizeCommentRecognizeErr(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeComment{l: l}
	s.Push(common.AugChar{
		C: common.Err,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	})
	ch := common.AugChar{
		C: '#',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
		Val: assert.AnError,
	}, l.tokens.Front().Value.(*common.Token))
}
