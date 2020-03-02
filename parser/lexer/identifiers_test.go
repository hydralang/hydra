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

func TestRecognizeIdentifierImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &recognizeIdentifier{})
}

func TestRecogIdentifier(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	result := recogIdentifier(l)

	r, ok := result.(*recognizeIdentifier)
	a.True(ok)
	a.Equal(l, r.l)
	a.NotNil(r.s)
}

func TestRecognizeIdentifierRecognizeBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("Nino"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 5},
		},
		Val: "Nino",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeIdentifierRecognizeNormalizeUnneeded(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("Ni\u00f1o"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 5},
		},
		Val: "Ni\u00f1o",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeIdentifierRecognizeNormalizeNeeded(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("Nin\u0303o"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokIdent,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 6},
		},
		Val: "Ni\u00f1o",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeIdentifierRecognizeString(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("rb\"spam\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokBytes,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 9},
		},
		Val: []byte("spam"),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeIdentifierRecognizeKeyword(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("kw1"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "kw1"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 4},
		},
		Val: "kw1",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeIdentifierRecognizeErr(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
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
		C: 's',
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

func TestRecognizeIdentifierRecognizeBadIdent(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("Nino\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeIdentifier{
		l: l,
		s: recogString(l).(*recognizeString),
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 5},
			E:    common.FilePos{L: 1, C: 6},
		},
		Val: common.ErrBadIdent,
	}, l.tokens.Front().Value.(*common.Token))
}
