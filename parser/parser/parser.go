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
	"errors"

	"github.com/hydralang/hydra/ast"
	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/lexer"
)

// parser is an implementation of Parser.
type parser struct {
	l    common.Lexer    // The lexer for the source
	opts *common.Options // The parser options
}

// Parse prepares a new parser from the parser options and the lexer.
// If the lexer is nil, one will be constructed from the options.
func Parse(opts *common.Options, l common.Lexer) (common.Parser, error) {
	// Construct the lexer
	if l == nil {
		var err error
		l, err = lexer.Lex(opts, nil)
		if err != nil {
			return nil, err
		}
	}

	// Construct the parser object
	p := &parser{
		l:    l,
		opts: opts,
	}

	return p, nil
}

// Expression parses a single expression from the output of the lexer.
// Note that this does not necessarily consume the entire input.  The
// method is called with a "right binding power", which is used to
// determine operator precedence.  An initial call should set this
// parameter to 0; calls by token parsers (see ParserTable) may pass
// different values, typically their left binding power.
func (p *parser) Expression(rbp int) (ast.Expression, error) {
	// Get the parse table
	pt := p.opts.Prof.ParseTab

	// First, grab a token off the lexer
	tok := p.l.Next()

	// Do the evaluation of the first expression
	expr, err := pt.ExprFirst(p, tok)
	if err != nil {
		p.l.Push(tok)
		return nil, err
	}

	// Now handle the rest of the input
	for tok = p.l.Next(); rbp < pt.BindingPower(p, tok); tok = p.l.Next() {
		expr, err = pt.ExprNext(p, expr, tok)
		if err != nil {
			p.l.Push(tok)
			return nil, err
		}
	}

	p.l.Push(tok)
	return expr, nil
}

// Statement parses a single statement from the output of the lexer.
// Note that this does not necessarily consume the entire input.
func (p *parser) Statement() (ast.Statement, error) {
	return nil, errors.New("not implemented")
}

// Module parses a module, or collection of statements, from the
// output of the lexer.  This is intended to consume the entire input.
func (p *parser) Module() (ast.Statement, error) {
	return nil, errors.New("not implemented")
}
