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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/hydra/parser/common"
)

func TestMockRecognizerImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &mockRecognizer{})
}

type mockRecognizer struct {
	mock.Mock
	l *lexer
}

func (m *mockRecognizer) RecogMock(l *lexer) Recognizer {
	m.l = l

	return m
}

func (m *mockRecognizer) Recognize(ch common.AugChar) {
	args := m.MethodCalled("Recognize", ch)

	tmpTok := args.Get(0)
	if tmpTok != nil {
		m.l.pushTok(
			tmpTok.(*common.Symbol),
			args.Get(1).(common.Location),
			args.Get(2),
		)
	}
}

type mockRecs struct {
	rComment *mockRecognizer
	rNumber  *mockRecognizer
	rIdent   *mockRecognizer
	rString  *mockRecognizer
	rOp      *mockRecognizer
}

type saveRecs struct {
	rComment RecogInit
	rNumber  RecogInit
	rIdent   RecogInit
	rString  RecogInit
	rOp      RecogInit
}

func newMockRecs() *mockRecs {
	return &mockRecs{
		rComment: &mockRecognizer{},
		rNumber:  &mockRecognizer{},
		rIdent:   &mockRecognizer{},
		rString:  &mockRecognizer{},
		rOp:      &mockRecognizer{},
	}
}

func (mr *mockRecs) Install() *saveRecs {
	prev := &saveRecs{
		rComment: rComment,
		rNumber:  rNumber,
		rIdent:   rIdent,
		rString:  rString,
		rOp:      rOp,
	}

	rComment = mr.rComment.RecogMock
	rNumber = mr.rNumber.RecogMock
	rIdent = mr.rIdent.RecogMock
	rString = mr.rString.RecogMock
	rOp = mr.rOp.RecogMock

	return prev
}

func (sr *saveRecs) Install() {
	rComment = sr.rComment
	rNumber = sr.rNumber
	rIdent = sr.rIdent
	rString = sr.rString
	rOp = sr.rOp
}

func (mr *mockRecs) AssertExpectations(t *testing.T) {
	mr.rComment.AssertExpectations(t)
	mr.rNumber.AssertExpectations(t)
	mr.rIdent.AssertExpectations(t)
	mr.rString.AssertExpectations(t)
	mr.rOp.AssertExpectations(t)
}
