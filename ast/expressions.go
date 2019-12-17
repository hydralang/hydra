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
	"fmt"

	"github.com/hydralang/hydra/utils"
)

// BaseExpression describes basics for an expression.  Note that
// BaseExpression does not, nor is it intended to, implement
// Expression.
type BaseExpression struct {
	Loc utils.Location // Location of the expression
}

// GetLoc retrieves the location of the expression.
func (b *BaseExpression) GetLoc() utils.Location {
	return b.Loc
}

// Constant describes a constant expression node.
type Constant struct {
	BaseExpression
	Val interface{} // Value of the constant
}

// Children implements the utils.Visitable interface.
func (c *Constant) Children() []utils.Visitable {
	return []utils.Visitable{}
}

// String implements the fmt.Stringer interface.
func (c *Constant) String() string {
	return fmt.Sprintf("%s: %v", c.Loc, c.Val)
}

// Variable describes a reference to a variable.
type Variable struct {
	BaseExpression
	Name string // Name of the variable
}

// Children implements the utils.Visitable interface.
func (v *Variable) Children() []utils.Visitable {
	return []utils.Visitable{}
}

// String implements the fmt.Stringer interface.
func (v *Variable) String() string {
	return fmt.Sprintf("%s: <%s>", v.Loc, v.Name)
}

// Unary describes the action of a unary operator on another
// expression node.
type Unary struct {
	BaseExpression
	Op   string     // The operation to perform
	Node Expression // The expression node acted upon
}

// Children implements the utils.Visitable interface.
func (u *Unary) Children() []utils.Visitable {
	return []utils.Visitable{
		utils.Annotated{
			Wrapped:    u.Node,
			Annotation: "Node: ",
		},
	}
}

// String implements the fmt.Stringer interface.
func (u *Unary) String() string {
	return fmt.Sprintf("%s: %s", u.Loc, u.Op)
}

// Binary describes the action of a binary operator on two expression
// nodes.
type Binary struct {
	BaseExpression
	Op    string     // The operation to perform
	Left  Expression // The left-hand expression
	Right Expression // The right-hand expression
}

// Children implements the utils.Visitable interface.
func (b *Binary) Children() []utils.Visitable {
	return []utils.Visitable{
		utils.Annotated{
			Wrapped:    b.Left,
			Annotation: "Left : ",
		},
		utils.Annotated{
			Wrapped:    b.Right,
			Annotation: "Right: ",
		},
	}
}

// String implements the fmt.Stringer interface.
func (b *Binary) String() string {
	return fmt.Sprintf("%s: %s", b.Loc, b.Op)
}

// Attribute describes the action of the "." operator on an expression
// node.
type Attribute struct {
	BaseExpression
	Expr Expression // The expression to seek the attribute of
	Attr string     // The name of the attribute to seek
}

// Children implements the utils.Visitable interface.
func (a *Attribute) Children() []utils.Visitable {
	return []utils.Visitable{
		utils.Annotated{
			Wrapped:    a.Expr,
			Annotation: "Expr: ",
		},
	}
}

// String implements the fmt.Stringer interface.
func (a *Attribute) String() string {
	return fmt.Sprintf("%s: .%s", a.Loc, a.Attr)
}

// Call describes a call to a function.
type Call struct {
	BaseExpression
	Func   Expression            // The function to be called
	Args   []Expression          // The function arguments
	KwArgs map[string]Expression // The function keyword args
}

// Children implements the utils.Visitable interface.
func (c *Call) Children() []utils.Visitable {
	// Pre-allocate the children list
	children := make([]utils.Visitable, len(c.Args)+len(c.KwArgs)+1)
	idx := 0

	// Add the function
	children[idx] = utils.Annotated{
		Wrapped:    c.Func,
		Annotation: "Func: ",
	}
	idx++

	// Add the positional arguments
	for i, arg := range c.Args {
		children[idx] = utils.Annotated{
			Wrapped:    arg,
			Annotation: fmt.Sprintf("[%d]: ", i),
		}
		idx++
	}

	// Now the keyword arguments; note the ordering is not
	// guaranteed
	for key, arg := range c.KwArgs {
		children[idx] = utils.Annotated{
			Wrapped:    arg,
			Annotation: fmt.Sprintf("'%s': ", key),
		}
		idx++
	}

	return children
}

// String implements the fmt.Stringer interface.
func (c *Call) String() string {
	return fmt.Sprintf("%s: Call", c.Loc)
}
