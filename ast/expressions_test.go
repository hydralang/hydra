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

func TestConstantImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Constant{})
}

func TestConstantGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Constant{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestConstantChildren(t *testing.T) {
	a := assert.New(t)
	obj := &Constant{}

	result := obj.Children()

	a.Equal([]utils.Visitable{}, result)
}

func TestConstantString(t *testing.T) {
	a := assert.New(t)
	obj := &Constant{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Val: "value",
	}

	result := obj.String()

	a.Equal("file:3:2: value", result)
}

func TestVariableImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Variable{})
}

func TestVariableGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Variable{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestVariableChildren(t *testing.T) {
	a := assert.New(t)
	obj := &Variable{}

	result := obj.Children()

	a.Equal([]utils.Visitable{}, result)
}

func TestVariableString(t *testing.T) {
	a := assert.New(t)
	obj := &Variable{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Name: "variable",
	}

	result := obj.String()

	a.Equal("file:3:2: <variable>", result)
}

func TestUnaryImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Unary{})
}

func TestUnaryGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Unary{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestUnaryChildren(t *testing.T) {
	a := assert.New(t)
	node := &MockExpression{}
	obj := &Unary{
		Node: node,
	}

	result := obj.Children()

	a.Equal([]utils.Visitable{
		utils.Annotated{
			Wrapped:    node,
			Annotation: "Node: ",
		},
	}, result)
}

func TestUnaryString(t *testing.T) {
	a := assert.New(t)
	obj := &Unary{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Op: "+",
	}

	result := obj.String()

	a.Equal("file:3:2: +", result)
}

func TestBinaryImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Binary{})
}

func TestBinaryGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Binary{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestBinaryChildren(t *testing.T) {
	a := assert.New(t)
	left := &MockExpression{}
	right := &MockExpression{}
	obj := &Binary{
		Left:  left,
		Right: right,
	}

	result := obj.Children()

	a.Equal([]utils.Visitable{
		utils.Annotated{
			Wrapped:    left,
			Annotation: "Left : ",
		},
		utils.Annotated{
			Wrapped:    right,
			Annotation: "Right: ",
		},
	}, result)
}

func TestBinaryString(t *testing.T) {
	a := assert.New(t)
	obj := &Binary{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Op: "+",
	}

	result := obj.String()

	a.Equal("file:3:2: +", result)
}

func TestAttributeImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Attribute{})
}

func TestAttributeGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Attribute{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestAttributeChildren(t *testing.T) {
	a := assert.New(t)
	expr := &MockExpression{}
	obj := &Attribute{
		Expr: expr,
	}

	result := obj.Children()

	a.Equal([]utils.Visitable{
		utils.Annotated{
			Wrapped:    expr,
			Annotation: "Expr: ",
		},
	}, result)
}

func TestAttributeString(t *testing.T) {
	a := assert.New(t)
	obj := &Attribute{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
		Attr: "attr",
	}

	result := obj.String()

	a.Equal("file:3:2: .attr", result)
}

func TestCallImplementsExpression(t *testing.T) {
	assert.Implements(t, (*Expression)(nil), &Call{})
}

func TestCallGetLoc(t *testing.T) {
	a := assert.New(t)
	obj := &Call{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.GetLoc()

	a.Equal(utils.Location{
		File: "file",
		B:    utils.FilePos{L: 3, C: 2},
		E:    utils.FilePos{L: 3, C: 3},
	}, result)
}

func TestCallChildren(t *testing.T) {
	a := assert.New(t)
	fun := &MockExpression{}
	args := []Expression{
		&MockExpression{},
		&MockExpression{},
		&MockExpression{},
	}
	kwargs := map[string]Expression{
		"keyword": &MockExpression{},
	}
	obj := &Call{
		Func:   fun,
		Args:   args,
		KwArgs: kwargs,
	}

	result := obj.Children()

	a.Equal([]utils.Visitable{
		utils.Annotated{
			Wrapped:    fun,
			Annotation: "Func: ",
		},
		utils.Annotated{
			Wrapped:    args[0],
			Annotation: "[0]: ",
		},
		utils.Annotated{
			Wrapped:    args[1],
			Annotation: "[1]: ",
		},
		utils.Annotated{
			Wrapped:    args[2],
			Annotation: "[2]: ",
		},
		utils.Annotated{
			Wrapped:    kwargs["keyword"],
			Annotation: "'keyword': ",
		},
	}, result)
}

func TestCallString(t *testing.T) {
	a := assert.New(t)
	obj := &Call{
		Loc: utils.Location{
			File: "file",
			B:    utils.FilePos{L: 3, C: 2},
			E:    utils.FilePos{L: 3, C: 3},
		},
	}

	result := obj.String()

	a.Equal("file:3:2: Call", result)
}
