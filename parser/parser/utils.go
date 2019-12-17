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
)

// LiteralExprFirst is an ExprFirst function that returns a literal
// value.
func LiteralExprFirst(p common.Parser, t *common.Token) (ast.Expression, error) {
	return &ast.Constant{
		BaseExpression: ast.BaseExpression{
			Loc: t.Loc,
		},
		Val: t.Val,
	}, nil
}

// InfixExprNext is a factory function for ExprNext functions for
// basic infix operators--e.g., "+", "-", "*", "/", etc.
func InfixExprNext(op string, lbp int) ExprNext {
	return func(p common.Parser, l ast.Expression, r *common.Token) (ast.Expression, error) {
		right, err := p.Expression(lbp)
		if err != nil {
			return nil, err
		}

		return &ast.Binary{
			BaseExpression: ast.BaseExpression{
				Loc: l.GetLoc().ThruEnd(r.Loc),
			},
			Op:    op,
			Left:  l,
			Right: right,
		}, nil
	}
}

// PrefixExprFirst is a factory function for ExprFirst functions for
// prefixed unary operators--e.g., "~", "-", etc.
func PrefixExprFirst(op string, lbp int) ExprFirst {
	return func(p common.Parser, t *common.Token) (ast.Expression, error) {
		expr, err := p.Expression(lbp)
		if err != nil {
			return nil, err
		}

		return &ast.Unary{
			BaseExpression: ast.BaseExpression{
				Loc: expr.GetLoc(),
			},
			Op:   op,
			Node: expr,
		}, nil
	}
}

// InfixRightExprNext is a factory function for ExprNext functions for
// right-associative infix operators--e.g., "**".
func InfixRightExprNext(op string, lbp int) ExprNext {
	return func(p common.Parser, l ast.Expression, r *common.Token) (ast.Expression, error) {
		right, err := p.Expression(lbp - 1)
		if err != nil {
			return nil, err
		}

		return &ast.Binary{
			BaseExpression: ast.BaseExpression{
				Loc: l.GetLoc().ThruEnd(r.Loc),
			},
			Op:    op,
			Left:  l,
			Right: right,
		}, nil
	}
}
