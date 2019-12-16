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

package common

import (
	"fmt"

	"github.com/hydralang/hydra/utils"
)

// Scanner is an interface describing a scanner.  A scanner reads a
// source character rune by character rune, returning augmented
// characters.
type Scanner interface {
	// Next retrieves the next rune from the file.  An EOF
	// augmented character is returned on end of file, and an Err
	// augmented character is returned in the event of an error.
	Next() AugChar

	// Push pushes back a single augmented character onto the
	// scanner.  Any number of characters may be pushed back.
	Push(ch AugChar)
}

// Lexer is an interface describing a lexer.  A lexer pulls characters
// from a scanner and converts them to tokens, which may then be used
// by the parser.
type Lexer interface {
	// Next retrieves the next token from the scanner.  If the end
	// of file is reached, an EOF token is returned; if an error
	// occurs while scanning or lexically analyzing the file, an
	// error token is returned with the error as the token's
	// semantic value.  After either an EOF token or an error
	// token, nil will be returned.
	Next() *Token

	// Push pushes a single token back onto the lexer.  Any number
	// of tokens may be pushed back.
	Push(tok *Token)
}

// Parser is an interface describing a parser.  A parser utilizes a
// lexer to tokenize the input stream and convert it into an
// appropriate abstract syntax tree.
type Parser interface {
	// Expression parses a single expression from the output of
	// the lexer.  Note that this does not necessarily consume the
	// entire input.  The method is called with a "right binding
	// power", which is used to determine operator precedence.  An
	// initial call should set this parameter to 0; calls by token
	// parsers (see ParserTable) may pass different values,
	// typically their left binding power.
	Expression(rbp int) (Expression, error)

	// Statement parses a single statement from the output of the
	// lexer.  Note that this does not necessarily consume the
	// entire input.
	Statement() (Statement, error)

	// Module parses a module, or collection of statements, from
	// the output of the lexer.  This is intended to consume the
	// entire input.
	Module() (Statement, error)
}

// ParserTable is an interface describing the table of symbols the
// parser uses during parsing.
type ParserTable interface {
	// BindingPower retrieves the binding power for a particular
	// expression token.  It is passed the token.  It returns the
	// binding power or an error.
	BindingPower(p Parser, t *Token) int

	// ExprFirst is called for the first expression token.  It is
	// passed the token.  It returns an expression or an error.
	ExprFirst(p Parser, t *Token) (Expression, error)

	// ExprNext is called for subsequent expression tokens.  It is
	// passed the left and right tokens.  It returns an expression
	// or an error.
	ExprNext(p Parser, l Expression, r *Token) (Expression, error)

	// Statement is called for statement tokens.  It returns a
	// statement or an error.
	Statement(p Parser, t *Token) (Statement, error)
}

// Expression is an interface describing an expression node in the
// abstract syntax tree.
type Expression interface {
	fmt.Stringer
	utils.Visitable

	// GetLoc retrieves the location of the expression.
	GetLoc() Location
}

// Statement is an interface describing a statement node in the
// abstract syntax tree.
type Statement interface {
	// GetLoc retrieves the location of the statement.
	GetLoc() Location
}
