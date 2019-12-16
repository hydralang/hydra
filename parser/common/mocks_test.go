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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/testutils"
)

func TestMockScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &MockScanner{})
}

func TestMockScannerNext(t *testing.T) {
	a := assert.New(t)
	s := &MockScanner{}
	s.On("Next").Return(AugChar{C: 'c'})

	result := s.Next()

	a.Equal(AugChar{C: 'c'}, result)
	s.AssertExpectations(t)
}

func TestMockScannerPush(t *testing.T) {
	s := &MockScanner{}
	s.On("Push", AugChar{C: 'c'})

	s.Push(AugChar{C: 'c'})

	s.AssertExpectations(t)
}

func TestMockLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*Lexer)(nil), &MockLexer{})
}

func TestMockLexerNextToken(t *testing.T) {
	a := assert.New(t)
	l := &MockLexer{}
	l.On("Next").Return(&Token{Sym: &Symbol{Name: "sym"}})

	result := l.Next()

	a.Equal(&Token{Sym: &Symbol{Name: "sym"}}, result)
	l.AssertExpectations(t)
}

func TestMockLexerNextNil(t *testing.T) {
	a := assert.New(t)
	l := &MockLexer{}
	l.On("Next").Return(nil)

	result := l.Next()

	a.Nil(result)
	l.AssertExpectations(t)
}

func TestMockLexerPush(t *testing.T) {
	l := &MockLexer{}
	l.On("Push", &Token{Sym: &Symbol{Name: "sym"}})

	l.Push(&Token{Sym: &Symbol{Name: "sym"}})

	l.AssertExpectations(t)
}

func TestMockParserImplementsParser(t *testing.T) {
	assert.Implements(t, (*Parser)(nil), &MockParser{})
}

func TestMockParserExpressionNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Expression", 5).Return(nil, nil)

	result, err := p.Expression(5)

	a.Nil(result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserExpressionNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	e := &MockExpression{}
	p.On("Expression", 5).Return(e, nil)

	result, err := p.Expression(5)

	testutils.AssertPtrEqual(a, e, result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserExpressionErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Expression", 5).Return(nil, errors.New("an error"))

	result, err := p.Expression(5)

	a.Nil(result)
	a.Error(err, "an error")
	p.AssertExpectations(t)
}

func TestMockParserStatementNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Statement").Return(nil, nil)

	result, err := p.Statement()

	a.Nil(result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserStatementNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	s := &MockStatement{}
	p.On("Statement").Return(s, nil)

	result, err := p.Statement()

	testutils.AssertPtrEqual(a, s, result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserStatementErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Statement").Return(nil, errors.New("an error"))

	result, err := p.Statement()

	a.Nil(result)
	a.Error(err, "an error")
	p.AssertExpectations(t)
}

func TestMockParserModuleNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Module").Return(nil, nil)

	result, err := p.Module()

	a.Nil(result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserModuleNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	s := &MockStatement{}
	p.On("Module").Return(s, nil)

	result, err := p.Module()

	testutils.AssertPtrEqual(a, s, result)
	a.NoError(err)
	p.AssertExpectations(t)
}

func TestMockParserModuleErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	p.On("Module").Return(nil, errors.New("an error"))

	result, err := p.Module()

	a.Nil(result)
	a.Error(err, "an error")
	p.AssertExpectations(t)
}

func TestMockParserTableImplementsParserTable(t *testing.T) {
	assert.Implements(t, (*ParserTable)(nil), &MockParserTable{})
}

func TestMockParserTableBindingPower(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("BindingPower", p, tok).Return(5)

	result := pt.BindingPower(p, tok)

	a.Equal(5, result)
	pt.AssertExpectations(t)
}

func TestMockParserTableExprFirstNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("ExprFirst", p, tok).Return(nil, nil)

	result, err := pt.ExprFirst(p, tok)

	a.Nil(result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableExprFirstNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	e := &MockExpression{}
	pt.On("ExprFirst", p, tok).Return(e, nil)

	result, err := pt.ExprFirst(p, tok)

	testutils.AssertPtrEqual(a, e, result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableExprFirstErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("ExprFirst", p, tok).Return(nil, errors.New("an error"))

	result, err := pt.ExprFirst(p, tok)

	a.Nil(result)
	a.Error(err, "an error")
	pt.AssertExpectations(t)
}

func TestMockParserTableExprNextNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	l := &MockExpression{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("ExprNext", p, l, tok).Return(nil, nil)

	result, err := pt.ExprNext(p, l, tok)

	a.Nil(result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableExprNextNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	l := &MockExpression{}
	tok := &Token{}
	pt := &MockParserTable{}
	e := &MockExpression{}
	pt.On("ExprNext", p, l, tok).Return(e, nil)

	result, err := pt.ExprNext(p, l, tok)

	testutils.AssertPtrEqual(a, e, result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableExprNextErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	l := &MockExpression{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("ExprNext", p, l, tok).Return(nil, errors.New("an error"))

	result, err := pt.ExprNext(p, l, tok)

	a.Nil(result)
	a.Error(err, "an error")
	pt.AssertExpectations(t)
}

func TestMockParserTableStatementNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("Statement", p, tok).Return(nil, nil)

	result, err := pt.Statement(p, tok)

	a.Nil(result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableStatementNonNil(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	s := &MockStatement{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("Statement", p, tok).Return(s, nil)

	result, err := pt.Statement(p, tok)

	testutils.AssertPtrEqual(a, s, result)
	a.NoError(err)
	pt.AssertExpectations(t)
}

func TestMockParserTableStatementErr(t *testing.T) {
	a := assert.New(t)
	p := &MockParser{}
	tok := &Token{}
	pt := &MockParserTable{}
	pt.On("Statement", p, tok).Return(nil, errors.New("an error"))

	result, err := pt.Statement(p, tok)

	a.Nil(result)
	a.Error(err, "an error")
	pt.AssertExpectations(t)
}

func TestMockExpressionImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &MockExpression{})
}

func TestMockExpressionGetLoc(t *testing.T) {
	a := assert.New(t)
	e := &MockExpression{}
	e.On("GetLoc").Return(Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	})

	result := e.GetLoc()

	a.Equal(Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}, result)
	e.AssertExpectations(t)
}

func TestMockStatementImplementsStatement(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &MockStatement{})
}

func TestMockStatementGetLoc(t *testing.T) {
	a := assert.New(t)
	s := &MockStatement{}
	s.On("GetLoc").Return(Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	})

	result := s.GetLoc()

	a.Equal(Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}, result)
	s.AssertExpectations(t)
}
