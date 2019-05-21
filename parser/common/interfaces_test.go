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

type mockScanner struct {
	mock.Mock
}

func (m *mockScanner) Next() AugChar {
	args := m.MethodCalled("Next")

	return args.Get(0).(AugChar)
}

func (m *mockScanner) Push(ch AugChar) {
	m.MethodCalled("Push", ch)
}

type mockLexer struct {
	mock.Mock
}

func (m *mockLexer) Next() *Token {
	args := m.MethodCalled("Next")

	tok := args.Get(0)
	if tok == nil {
		return nil
	}
	return tok.(*Token)
}

func (m *mockLexer) Push(tok *Token) {
	m.MethodCalled("Push", tok)
}
