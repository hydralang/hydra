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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/hydra/utils"
)

// MockExpression is a mock object for expressions.
type MockExpression struct {
	mock.Mock
}

// GetLoc retrieves the location of the expression.
func (m *MockExpression) GetLoc() utils.Location {
	args := m.MethodCalled("GetLoc")

	return args.Get(0).(utils.Location)
}

// Children implements the utils.Visitable interface.
func (m *MockExpression) Children() []utils.Visitable {
	args := m.MethodCalled("Children")

	return args.Get(0).([]utils.Visitable)
}

// String implements the fmt.Stringer interface.
func (m *MockExpression) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

// MockStatement is a mock object for statements.
type MockStatement struct {
	mock.Mock
}

// GetLoc retrieves the location of the statement.
func (m *MockStatement) GetLoc() utils.Location {
	args := m.MethodCalled("GetLoc")

	return args.Get(0).(utils.Location)
}
