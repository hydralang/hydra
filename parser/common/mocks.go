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
