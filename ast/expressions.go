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

// Constant describes a constant expression node.
type Constant struct {
	Loc utils.Location // Location of the constant
	Val interface{}    // Value of the constant
}

// GetLoc retrieves the location of the expression.
func (c *Constant) GetLoc() utils.Location {
	return c.Loc
}

// Children implements the utils.Visitable interface.
func (c *Constant) Children() []utils.Visitable {
	return []utils.Visitable{}
}

// String implements the fmt.Stringer interface.
func (c *Constant) String() string {
	return fmt.Sprintf("%s: %v", c.Loc, c.Val)
}

// Unary describes the action of a unary operator on another
// expression node.
type Unary struct {
	Loc  utils.Location // Location of the unary operator
	Op   string         // The operation to perform
	Node Expression     // The expression node acted upon
}

// GetLoc retrieves the location of the expression.
func (u *Unary) GetLoc() utils.Location {
	return u.Loc
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
	Loc   utils.Location // Location of the binary operator
	Op    string         // The operation to perform
	Left  Expression     // The left-hand expression
	Right Expression     // The right-hand expression
}

// GetLoc retrieves the location of the expression.
func (b *Binary) GetLoc() utils.Location {
	return b.Loc
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

// Call describes a call to a function.
type Call struct {
	Loc    utils.Location        // Location of the function call
	Func   Expression            // The function to be called
	Args   []Expression          // The function arguments
	KwArgs map[string]Expression // The function keyword args
}

// GetLoc retrieves the location of the expression.
func (c *Call) GetLoc() utils.Location {
	return c.Loc
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
