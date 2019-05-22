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
	"unicode"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
	"github.com/stretchr/testify/assert"
)

func TestBufStringImplementsBuffer(t *testing.T) {
	assert.Implements(t, (*buffer)(nil), &bufString{})
}

func TestBufStringPutC(t *testing.T) {
	a := assert.New(t)
	obj := &bufString{}

	err := obj.putC('c')

	a.NoError(err)
	a.Equal("c", obj.String())
}

func TestBufStringPutCFails(t *testing.T) {
	a := assert.New(t)
	obj := &bufString{}

	err := obj.putC(unicode.MaxRune + 1)

	a.Equal(common.ErrBadStrChar, err)
	a.Equal("", obj.String())
}

func TestBufStringGet(t *testing.T) {
	a := assert.New(t)
	obj := &bufString{}
	obj.WriteString("test")

	result := obj.get()

	a.Equal("test", result)
}

func TestBufStringSym(t *testing.T) {
	a := assert.New(t)
	obj := &bufString{}

	result := obj.sym()

	a.Equal(common.TokString, result)
}

func TestBufBytesImplementsBuffer(t *testing.T) {
	assert.Implements(t, (*buffer)(nil), &bufBytes{})
}

func TestBufBytesPutC(t *testing.T) {
	a := assert.New(t)
	obj := &bufBytes{}

	err := obj.putC('c')

	a.NoError(err)
	a.Equal("c", obj.String())
}

func TestBufBytesPutCFails(t *testing.T) {
	a := assert.New(t)
	obj := &bufBytes{}

	err := obj.putC(0x100)

	a.Equal(common.ErrBadStrChar, err)
	a.Equal("", obj.String())
}

func TestBufBytesGet(t *testing.T) {
	a := assert.New(t)
	obj := &bufBytes{}
	obj.WriteString("test")

	result := obj.get()

	a.Equal([]byte("test"), result)
}

func TestBufBytesSym(t *testing.T) {
	a := assert.New(t)
	obj := &bufBytes{}

	result := obj.sym()

	a.Equal(common.TokBytes, result)
}

func TestRecognizeStringImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &recognizeString{})
}

func TestRecogString(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	result := recogString(l)

	r, ok := result.(*recognizeString)
	a.True(ok)
	a.Equal(l, r.l)
	a.Equal(uint8(0), r.flags)
	a.Equal(common.Location{}, r.loc)
	a.Equal(rune(0), r.q)
	a.Equal(0, r.qcnt)
	a.Equal(common.Location{}, r.runLoc)
}

func TestRecognizeStringSetFlagStrFlag(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	l := &lexer{
		opts: opts,
	}
	r := &recognizeString{l: l}
	ch := common.AugChar{
		C:     'b',
		Class: common.CharStrFlag,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}

	result := r.setFlag(ch)

	a.Equal(r, result)
	a.Equal(common.StrBytes, r.flags)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}, r.loc)
}

func TestRecognizeStringSetFlagQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	l := &lexer{
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrBytes,
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	ch := common.AugChar{
		C:     '"',
		Class: common.CharQuote,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}

	result := r.setFlag(ch)

	a.Equal(r, result)
	a.Equal(common.StrBytes|common.StrTriple, r.flags)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 2},
	}, r.loc)
}

func TestRecognizeStringSetFlagOther(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	l := &lexer{
		opts: opts,
	}
	r := &recognizeString{l: l}
	ch := common.AugChar{
		C:     '!',
		Class: 0,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}

	result := r.setFlag(ch)

	a.Nil(result)
}

func TestRecognizeStringEscapeRaw(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\a"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrRaw,
		buf:   &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.NoError(err)
	a.Equal("\\a", r.buf.get())
	a.Equal(common.Location{}, loc)
}

func TestRecognizeStringEscapeRawBadEscape(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrRaw,
		buf:   &bufString{},
	}
	ch := common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: 'a',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadStrChar, err)
	a.Equal("", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 2},
	}, loc)
}

func TestRecognizeStringEscapeRawErr(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrRaw,
		buf:   &bufString{},
	}
	ch := common.AugChar{
		C: '\\',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
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

	loc, err := r.escape(ch)

	a.Equal(assert.AnError, err)
	a.Equal("\\", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 2},
		E:    common.FilePos{L: 3, C: 3},
	}, loc)
}

func TestRecognizeStringEscapeRawEOF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrRaw,
		buf:   &bufString{},
	}
	ch := common.AugChar{
		C: '\\',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: common.EOF,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})

	loc, err := r.escape(ch)

	a.Equal(common.ErrUnclosedStr, err)
	a.Equal("\\", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 3},
	}, loc)
}

func TestRecognizeStringEscapeRawBadRune(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:     l,
		flags: common.StrRaw,
		buf:   &bufString{},
	}
	ch := common.AugChar{
		C: '\\',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadStrChar, err)
	a.Equal("\\", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 3},
	}, loc)
}

func TestRecognizeStringEscape(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\a"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.NoError(err)
	a.Equal("\a", r.buf.get())
	a.Equal(common.Location{}, loc)
}

func TestRecognizeStringEscapeQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.NoError(err)
	a.Equal("\"", r.buf.get())
	a.Equal(common.Location{}, loc)
}

func TestRecognizeStringEscapeNewline(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\\n"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.NoError(err)
	a.Equal("", r.buf.get())
	a.Equal(common.Location{}, loc)
}

func TestRecognizeStringEscapeQuoteBadQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	opts.Prof = opts.Prof.Copy()
	opts.Prof.Quotes = map[rune]uint8{
		unicode.MaxRune + 1: 0,
	}
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := common.AugChar{
		C: '\\',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
	}
	s.Push(common.AugChar{
		C:     unicode.MaxRune + 1,
		Class: common.CharQuote,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	})

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadStrChar, err)
	a.Equal("", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 3},
	}, loc)
}

func TestRecognizeStringEscapeEarlyEOF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\x"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadEscape, err)
	a.Equal("", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 3},
	}, loc)
}

func TestRecognizeStringEscapeBadChar(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\u0100"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufBytes{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadStrChar, err)
	a.Nil(r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 7},
	}, loc)
}

func TestRecognizeStringEscapeUnknown(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\\!"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	r := &recognizeString{
		l:   l,
		buf: &bufString{},
	}
	ch := l.s.Next()

	loc, err := r.escape(ch)

	a.Equal(common.ErrBadEscape, err)
	a.Equal("", r.buf.get())
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 1, C: 1},
		E:    common.FilePos{L: 1, C: 3},
	}, loc)
}

func TestRecognizeStringRecognizeEmpty(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 3},
		},
		Val: "",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeEmptyTriple(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"\"\"\"\"\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 7},
		},
		Val: "",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeBasic(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"spam\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 7},
		},
		Val: "spam",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeBasicTriple(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"\"\"spam\"\"\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 11},
		},
		Val: "spam",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeBytes(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("b\"spam\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	r.setFlag(l.s.Next())
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokBytes,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 8},
		},
		Val: []byte("spam"),
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeTripleInclusions(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"\"\"s\"p\"\"am\"\"\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 14},
		},
		Val: "s\"p\"\"am",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeTripleInclusionsBadQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l:     l,
		flags: common.StrTriple,
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: 'p',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 6},
			E:    common.FilePos{L: 1, C: 7},
		},
	})
	s.Push(common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 5},
			E:    common.FilePos{L: 1, C: 6},
		},
	})
	s.Push(common.AugChar{
		C: 's',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 5},
		},
	})
	s.Push(common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	})
	s.Push(common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 3},
		},
	})
	ch := common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}

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
		Val: common.ErrBadStrChar,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeReadError(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l:     l,
		flags: common.StrTriple,
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: common.Err,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 5},
		},
		Val: assert.AnError,
	})
	s.Push(common.AugChar{
		C: 'p',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	})
	s.Push(common.AugChar{
		C: 's',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 3},
		},
	})
	ch := common.AugChar{
		C: '"',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 5},
		},
		Val: assert.AnError,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeEOF(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l:     l,
		flags: common.StrTriple,
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}
	s.Push(common.AugChar{
		C: common.EOF,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 5},
		},
	})
	s.Push(common.AugChar{
		C: 'p',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	})
	s.Push(common.AugChar{
		C: 's',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 3},
		},
	})
	ch := common.AugChar{
		C: '"',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 5},
		},
		Val: common.ErrUnclosedStr,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeEscape(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"sp\\am\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 8},
		},
		Val: "sp\am",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeEscapeBadEscape(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"sp\\!am\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 1, C: 6},
		},
		Val: common.ErrBadEscape,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeNewline(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"sp\nam\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 4},
			E:    common.FilePos{L: 2, C: 1},
		},
		Val: common.ErrUnclosedStr,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeNewlineTripleQuote(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("\"\"\"sp\nam\"\"\""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	ch := l.s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokString,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 2, C: 6},
		},
		Val: "sp\nam",
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeStringRecognizeBadRune(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeString{
		l: l,
	}
	s.Push(common.AugChar{
		C: unicode.MaxRune + 1,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
	})
	s.Push(common.AugChar{
		C: 's',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 2},
			E:    common.FilePos{L: 1, C: 3},
		},
	})
	ch := common.AugChar{
		C: '"',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
	}

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 3},
			E:    common.FilePos{L: 1, C: 4},
		},
		Val: common.ErrBadStrChar,
	}, l.tokens.Front().Value.(*common.Token))
}
