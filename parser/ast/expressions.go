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

import "github.com/hydralang/hydra/parser/common"

// Constant describes a constant expression node.
type Constant struct {
	Loc common.Location // Location of the constant
	Val interface{}     // Value of the constant
}

// Unary describes the action of a unary operator on another
// expression node.
type Unary struct {
	Loc  common.Location   // Location of the unary operator
	Op   string            // The operation to perform
	Node common.Expression // The expression node acted upon
}

// Binary describes the action of a binary operator on two expression
// nodes.
type Binary struct {
	Loc   common.Location   // Location of the binary operator
	Op    string            // The operation to perform
	Left  common.Expression // The left-hand expression
	Right common.Expression // The right-hand expression
}

// Call describes a call to a function.
type Call struct {
	Loc    common.Location              // Location of the function call
	Func   common.Expression            // The function to be called
	Args   []common.Expression          // The function arguments
	KwArgs map[string]common.Expression // The function keyword args
}
