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

package parser

import (
	"github.com/hydralang/hydra/ast"
	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/utils"
)

// ExprFirst is called to process the first token in a sub-expression.
type ExprFirst func(p common.Parser, t *common.Token) (ast.Expression, error)

// ExprNext is called to process subsequent tokens in a
// sub-expression.
type ExprNext func(p common.Parser, l ast.Expression, r *common.Token) (ast.Expression, error)

// Statement is called to process statement tokens.
type Statement func(p common.Parser, t *common.Token) (ast.Statement, error)

// parserEntry is an entry in the parser table.
type parserEntry struct {
	Lbp       int       // Left binding power
	ExprFirst ExprFirst // Called to process first expr token
	ExprNext  ExprNext  // Called to process next expr token
	Statement Statement // Called to process statement tokens
}

// parserTable is an implementation of the ParserTable interface.
type parserTable map[string]parserEntry

// BindingPower retrieves the binding power for a particular
// expression token.  It is passed the token.  It returns the binding
// power or an error.
func (pt parserTable) BindingPower(p common.Parser, t *common.Token) int {
	ent, ok := pt[t.Sym.Name]
	if !ok || ent.Lbp <= 0 {
		return 0
	}

	return ent.Lbp
}

// ExprFirst is called for the first expression token.  It is passed
// the token and the associated left binding power of the token's
// symbol.  It returns an expression or an error.
func (pt parserTable) ExprFirst(p common.Parser, t *common.Token) (ast.Expression, error) {
	ent, ok := pt[t.Sym.Name]
	if !ok || ent.ExprFirst == nil {
		return nil, utils.ErrUnexpected
	}

	return ent.ExprFirst(p, t)
}

// ExprNext is called for subsequent expression tokens.  It is passed
// the left and right tokens and the associated left binding power of
// the left token, which is used to determine how to recurse.  It
// returns an expression or an error.
func (pt parserTable) ExprNext(p common.Parser, l ast.Expression, r *common.Token) (ast.Expression, error) {
	ent, ok := pt[r.Sym.Name]
	if !ok || ent.ExprNext == nil {
		return nil, utils.ErrUnexpected
	}

	return ent.ExprNext(p, l, r)
}

// Statement is called for statement tokens.  It returns a statement
// or an error.
func (pt parserTable) Statement(p common.Parser, t *common.Token) (ast.Statement, error) {
	ent, ok := pt[t.Sym.Name]
	if !ok || ent.Statement == nil {
		return nil, utils.ErrUnexpected
	}

	return ent.Statement(p, t)
}
