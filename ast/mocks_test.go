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

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/utils"
)

func TestMockExpressionImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &MockExpression{})
}

func TestMockExpressionGetLoc(t *testing.T) {
	a := assert.New(t)
	e := &MockExpression{}
	e.On("GetLoc").Return(utils.Location{
		File: "file",
		B: utils.FilePos{
			L: 3,
			C: 2,
		},
		E: utils.FilePos{
			L: 3,
			C: 3,
		},
	})

	result := e.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B: utils.FilePos{
			L: 3,
			C: 2,
		},
		E: utils.FilePos{
			L: 3,
			C: 3,
		},
	}, result)
	e.AssertExpectations(t)
}

func TestMockExpressionChildren(t *testing.T) {
	a := assert.New(t)
	children := []utils.Visitable{
		&MockExpression{},
		&MockExpression{},
		&MockExpression{},
	}
	e := &MockExpression{}
	e.On("Children").Return(children)

	result := e.Children()

	a.Equal(children, result)
	e.AssertExpectations(t)
}

func TestMockExpressionString(t *testing.T) {
	a := assert.New(t)
	e := &MockExpression{}
	e.On("String").Return("string!")

	result := e.String()

	a.Equal("string!", result)
	e.AssertExpectations(t)
}

func TestMockStatementImplementsStatement(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &MockStatement{})
}

func TestMockStatementGetLoc(t *testing.T) {
	a := assert.New(t)
	s := &MockStatement{}
	s.On("GetLoc").Return(utils.Location{
		File: "file",
		B: utils.FilePos{
			L: 3,
			C: 2,
		},
		E: utils.FilePos{
			L: 3,
			C: 3,
		},
	})

	result := s.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B: utils.FilePos{
			L: 3,
			C: 2,
		},
		E: utils.FilePos{
			L: 3,
			C: 3,
		},
	}, result)
	s.AssertExpectations(t)
}
