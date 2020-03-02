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

func TestRecognizeOperatorImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &recognizeOperator{})
}

func TestRecogOperator(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	l := &lexer{
		opts: opts,
	}

	result := recogOperator(l)

	r, ok := result.(*recognizeOperator)
	a.True(ok)
	a.Equal(l, r.l)
	a.Equal(opts.Prof.Operators, r.node)
}

func TestRecognizeOperatorPushFrameBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	r := &recognizeOperator{}
	ch := common.AugChar{
		C: '!',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}

	r.pushFrame(ch, opts.Prof.Operators)

	a.Equal(1, r.st.Len())
	frame := r.st.Front().Value.(*opFrame)
	a.Equal(&opFrame{
		ch:   ch,
		loc:  ch.Loc,
		node: opts.Prof.Operators,
	}, frame)
	a.Equal(opts.Prof.Operators, r.node)
}

func TestRecognizeOperatorPushFramePushes(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	r := &recognizeOperator{}
	ch1 := common.AugChar{
		C: '$',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch2 := common.AugChar{
		C: '$',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	node1 := opts.Prof.Operators.Next('$')
	r.st.PushBack(&opFrame{
		ch:   ch1,
		loc:  ch1.Loc,
		node: node1,
	})
	node2 := node1.Next('$')

	r.pushFrame(ch2, node2)

	a.Equal(2, r.st.Len())
	elem := r.st.Front()
	frame := elem.Value.(*opFrame)
	a.Equal(&opFrame{
		ch:   ch1,
		loc:  ch1.Loc,
		node: node1,
	}, frame)
	elem = elem.Next()
	frame = elem.Value.(*opFrame)
	a.Equal(&opFrame{
		ch:   ch2,
		loc:  ch2.Loc,
		node: node2,
	}, frame)
	a.Equal(node2, r.node)
}

func TestRecognizeOperatorPushFrameReplaces(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	r := &recognizeOperator{}
	ch1 := common.AugChar{
		C: '$',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 3},
		},
	}
	ch2 := common.AugChar{
		C: '$',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
	}
	ch3 := common.AugChar{
		C: '$',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}
	node1 := opts.Prof.Operators.Next('$')
	r.st.PushBack(&opFrame{
		ch:   ch1,
		loc:  ch1.Loc,
		node: node1,
	})
	node2 := node1.Next('$')
	r.st.PushBack(&opFrame{
		ch:   ch2,
		loc:  ch2.Loc,
		node: node2,
	})
	node3 := node2.Next('$')

	r.pushFrame(ch3, node3)

	a.Equal(1, r.st.Len())
	elem := r.st.Front()
	frame := elem.Value.(*opFrame)
	a.Equal(&opFrame{
		ch: ch3,
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 2},
			E:    common.FilePos{L: 3, C: 5},
		},
		node: node3,
	}, frame)
	a.Equal(node3, r.node)
}

func TestRecognizeOperatorEmitBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: '$',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 3},
				E:    common.FilePos{L: 3, C: 4},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 4},
		},
		node: opts.Prof.Operators.Next('$').Next('$').Next('$'),
	})

	r.emit()

	a.Equal(s, l.s)
	a.Equal(0, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "$$$"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 4},
		},
		Val: "$$$",
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal(common.EOF, ch.C)
}

func TestRecognizeOperatorEmitExtraChars(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: '$',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 3},
				E:    common.FilePos{L: 3, C: 4},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 4},
		},
		node: opts.Prof.Operators.Next('$').Next('$').Next('$'),
	})
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: '!',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 4},
				E:    common.FilePos{L: 3, C: 5},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
		node: opts.Prof.Operators,
	})

	r.emit()

	a.Equal(s, l.s)
	a.Equal(0, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "$$$"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 4},
		},
		Val: "$$$",
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal(common.AugChar{
		C: '!',
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 4},
			E:    common.FilePos{L: 3, C: 5},
		},
	}, ch)
}

func TestRecognizeOperatorEmitNoSym(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: '$',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 3},
				E:    common.FilePos{L: 3, C: 4},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 4},
		},
		node: opts.Prof.Operators.Next('$').Next('$'),
	})

	r.emit()

	a.Nil(l.s)
	a.Equal(0, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 3},
			E:    common.FilePos{L: 3, C: 4},
		},
		Val: common.ErrBadOp,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeOperatorEmitOpen(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: '(',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 1},
				E:    common.FilePos{L: 3, C: 2},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		node: opts.Prof.Operators.Next('('),
	})

	r.emit()

	a.Equal(s, l.s)
	a.Equal(1, l.pair.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "(", Close: ")"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "(",
	}, l.pair.Back().Value.(*common.Token))
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "(", Close: ")"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: "(",
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal(common.EOF, ch.C)
}

func TestRecognizeOperatorEmitClose(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	l.pair.PushBack(&common.Token{
		Sym: &common.Symbol{Name: "(", Close: ")"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: "(",
	})
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: ')',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 1},
				E:    common.FilePos{L: 3, C: 2},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		node: opts.Prof.Operators.Next(')'),
	})

	r.emit()

	a.Equal(s, l.s)
	a.Equal(0, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: ")", Open: "("},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		Val: ")",
	}, l.tokens.Front().Value.(*common.Token))
	ch := s.Next()
	a.Equal(common.EOF, ch.C)
}

func TestRecognizeOperatorEmitCloseMismatch(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	l.pair.PushBack(&common.Token{
		Sym: &common.Symbol{Name: "[", Close: "]"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: "(",
	})
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: ')',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 1},
				E:    common.FilePos{L: 3, C: 2},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		node: opts.Prof.Operators.Next(')'),
	})

	r.emit()

	a.Nil(l.s)
	a.Equal(1, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokError, tok.Sym)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 2},
	}, tok.Loc)
	a.EqualError(tok.Val.(error), "close operator \")\" does not match open operator \"[\" at file:1:1")
}

func TestRecognizeOperatorEmitCloseNoOpen(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{l: l}
	r.st.PushBack(&opFrame{
		ch: common.AugChar{
			C: ')',
			Loc: common.Location{
				File: "file",
				B:    common.FilePos{L: 3, C: 1},
				E:    common.FilePos{L: 3, C: 2},
			},
		},
		loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 3, C: 1},
			E:    common.FilePos{L: 3, C: 2},
		},
		node: opts.Prof.Operators.Next(')'),
	})

	r.emit()

	a.Nil(l.s)
	a.Equal(0, l.pair.Len())
	a.Equal(1, l.tokens.Len())
	tok := l.tokens.Front().Value.(*common.Token)
	a.Equal(common.TokError, tok.Sym)
	a.Equal(common.Location{
		File: "file",
		B:    common.FilePos{L: 3, C: 1},
		E:    common.FilePos{L: 3, C: 2},
	}, tok.Loc)
	a.EqualError(tok.Val.(error), "unexpected close operator \")\"")
}

func TestRecognizeOperatorRecognizeBase(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("$$$"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{
		l:    l,
		node: opts.Prof.Operators,
	}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "$$$"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 4},
		},
		Val: "$$$",
	}, l.tokens.Front().Value.(*common.Token))
	ch = s.Next()
	a.Equal(common.EOF, ch.C)
}

func TestRecognizeOperatorRecognizeNonop(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("$$$a"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{
		l:    l,
		node: opts.Prof.Operators,
	}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "$$$"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 4},
		},
		Val: "$$$",
	}, l.tokens.Front().Value.(*common.Token))
	ch = s.Next()
	a.Equal('a', ch.C)
}

func TestRecognizeOperatorRecognizeOpLookalike(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("$$$@"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{
		l:    l,
		node: opts.Prof.Operators,
	}
	ch := s.Next()

	r.Recognize(ch)

	a.Equal(s, l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: &common.Symbol{Name: "$$$"},
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 4},
		},
		Val: "$$$",
	}, l.tokens.Front().Value.(*common.Token))
	ch = s.Next()
	a.Equal('@', ch.C)
}

func TestRecognizeOperatorRecognizeNotOp(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader("@"))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{
		l:    l,
		node: opts.Prof.Operators,
	}
	ch := s.Next()

	r.Recognize(ch)

	a.Nil(l.s)
	a.Equal(1, l.tokens.Len())
	a.Equal(&common.Token{
		Sym: common.TokError,
		Loc: common.Location{
			File: "file",
			B:    common.FilePos{L: 1, C: 1},
			E:    common.FilePos{L: 1, C: 2},
		},
		Val: common.ErrBadOp,
	}, l.tokens.Front().Value.(*common.Token))
}

func TestRecognizeOperatorRecognizeErr(t *testing.T) {
	a := assert.New(t)
	opts := makeOptions(strings.NewReader(""))
	s, _ := scanner.Scan(opts)
	l := &lexer{
		s:    s,
		opts: opts,
	}
	l.indent.PushBack(1)
	r := &recognizeOperator{
		l:    l,
		node: opts.Prof.Operators,
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
		C: '$',
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
