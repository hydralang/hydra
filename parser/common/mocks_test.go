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
	"testing"

	"github.com/stretchr/testify/assert"
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
