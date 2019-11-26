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

import "github.com/stretchr/testify/mock"

// MockScanner is a mock object for scanners.
type MockScanner struct {
	mock.Mock
}

// Next retrieves the next rune from the file.  An EOF augmented
// character is returned on end of file, and an Err augmented
// character is returned in the event of an error.
func (m *MockScanner) Next() AugChar {
	args := m.MethodCalled("Next")

	return args.Get(0).(AugChar)
}

// Push pushes back a single augmented character onto the scanner.
// Any number of characters may be pushed back.
func (m *MockScanner) Push(ch AugChar) {
	m.MethodCalled("Push", ch)
}

// MockLexer is a mock object for lexers.
type MockLexer struct {
	mock.Mock
}

// Next retrieves the next token from the scanner.  If the end of file
// is reached, an EOF token is returned; if an error occurs while
// scanning or lexically analyzing the file, an error token is
// returned with the error as the token's semantic value.  After
// either an EOF token or an error token, nil will be returned.
func (m *MockLexer) Next() *Token {
	args := m.MethodCalled("Next")

	tok := args.Get(0)
	if tok == nil {
		return nil
	}

	return tok.(*Token)
}

// Push pushes a single token back onto the lexer.  Any number of
// tokens may be pushed back.
func (m *MockLexer) Push(tok *Token) {
	m.MethodCalled("Push", tok)
}

// MockParser is a mock object for parsers.
type MockParser struct {
	mock.Mock
}

// Expression parses a single expression from the output of the lexer.
// Note that this does not necessarily consume the entire input.  The
// method is called with a "right binding power", which is used to
// determine operator precedence.  An initial call should set this
// parameter to 0; calls by token parsers (see ParserTable) may pass
// different values, typically their left binding power.
func (m *MockParser) Expression(rbp int) (Expression, error) {
	args := m.MethodCalled("Expression", rbp)

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Expression), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statement parses a single statement from the output of the lexer.
// Note that this does not necessarily consume the entire input.
func (m *MockParser) Statement() (Statement, error) {
	args := m.MethodCalled("Statement")

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Statement), args.Error(1)
	}

	return nil, args.Error(1)
}

// Module parses a module, or collection of statements, from the
// output of the lexer.  This is intended to consume the entire input.
func (m *MockParser) Module() (Statement, error) {
	args := m.MethodCalled("Module")

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Statement), args.Error(1)
	}

	return nil, args.Error(1)
}

// MockParserTable is a mock object for parser tables.
type MockParserTable struct {
	mock.Mock
}

// ExprFirst is called for the first expression token.  It is passed
// the token.  It returns an expression or an error.
func (m *MockParserTable) ExprFirst(p Parser, t *Token) (Expression, error) {
	args := m.MethodCalled("ExprFirst", p, t)

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Expression), args.Error(1)
	}

	return nil, args.Error(1)
}

// ExprNext is called for subsequent expression tokens.  It is passed
// the left and right tokens.  It returns an expression or an error.
func (m *MockParserTable) ExprNext(p Parser, l Expression, r *Token) (Expression, error) {
	args := m.MethodCalled("ExprNext", p, l, r)

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Expression), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statement is called for statement tokens.  It returns a statement
// or an error.
func (m *MockParserTable) Statement(p Parser, t *Token) (Statement, error) {
	args := m.MethodCalled("Statement", p, t)

	tmp := args.Get(0)
	if tmp != nil {
		return tmp.(Statement), args.Error(1)
	}

	return nil, args.Error(1)
}

// MockExpression is a mock object for expressions.
type MockExpression struct {
	mock.Mock
}

// MockStatement is a mock object for statements.
type MockStatement struct {
	mock.Mock
}
