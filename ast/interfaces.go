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

// Expression is an interface describing an expression node in the
// abstract syntax tree.
type Expression interface {
	fmt.Stringer
	utils.Visitable

	// GetLoc retrieves the location of the expression.
	GetLoc() utils.Location
}

// Statement is an interface describing a statement node in the
// abstract syntax tree.
type Statement interface {
	// GetLoc retrieves the location of the statement.
	GetLoc() utils.Location
}
