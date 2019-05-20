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

package lexer

import (
	"container/list"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/parser/scanner"
)

// Default recognizers.  Defining these as variables enables the
// recognizers to be mocked out for testing purposes.
var (
	rComment RecogInit = recogComment
	rNumber  RecogInit = recogNumber
	rIdent   RecogInit = recogIdentifier
	rString  RecogInit = recogString
	rOp      RecogInit = recogOperator
)

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
	Next() *common.Token

	// Push pushes a single token back onto the lexer.  Any number
	// of tokens may be pushed back.
	Push(tok *common.Token)
}

// lexer is an implementation of Lexer.
type lexer struct {
	s       scanner.Scanner // The scanner for the source
	opts    *common.Options // The parser options
	indent  list.List       // The indent stack
	pair    list.List       // The pairing stack
	tokens  list.List       // The token stack
	prevTok *common.Token   // Last token returned by lexer
}

// Lex prepares a new lexer from the parser options and the scanner.
// If the scanner is nil, one will be constructed from the options.
func Lex(opts *common.Options, s scanner.Scanner) (Lexer, error) {
	// Construct the scanner
	if s == nil {
		var err error
		s, err = scanner.Scan(opts)
		if err != nil {
			return nil, err
		}
	}

	// Construct the lexer object
	l := &lexer{
		s:    s,
		opts: opts,
	}

	// Push the starting column onto the indent stack
	l.indent.PushBack(1)

	return l, nil
}

// Next retrieves the next token from the scanner.  If the end of file
// is reached, an EOF token is returned; if an error occurs while
// scanning or lexically analyzing the file, an error token is
// returned with the error as the token's semantic value.  After
// either an EOF token or an error token, nil will be returned.
func (l *lexer) Next() *common.Token {
	// Pump some tokens onto the token stack
	for l.s != nil && l.tokens.Len() == 0 {
		// Get a character from the scanner
		ch := l.s.Next()

		// Handle EOF and error
		if ch.C == common.Err {
			l.pushErr(ch.Loc, ch.Val.(error))
			break
		} else if ch.C == common.EOF {
			// Warn about dangling pairs
			if l.pair.Len() > 0 {
				dangle := l.pair.Back().Value.(*common.Token)
				l.pushErr(dangle.Loc, common.ErrDanglingOpen(dangle))
				break
			}

			l.pushTok(common.TokEOF, ch.Loc, nil)
			l.s = nil
			break
		}

		// Handle newlines and whitespace
		if ch.Class&common.CharNL != 0 && l.pair.Len() == 0 {
			// Generate a newline token
			l.pushTok(common.TokNewline, ch.Loc, nil)
			break
		} else if ch.Class&common.CharWS != 0 {
			// Are we concerned about mixed spaces?
			errMixed := false

			// Set up the skipSpaces flags
			var skip uint8
			if l.pair.Len() > 0 {
				skip = SkipNL
			} else {
				prevTok := l.lastTok()
				if prevTok == nil || prevTok.Sym == common.TokNewline {
					skip = SkipLeadFF
					errMixed = true
				}
			}

			// Skip the whitespace
			mixed := l.skipSpaces(ch, skip)

			// Error out if it's mixed
			if errMixed && mixed {
				l.pushErr(ch.Loc, common.ErrMixedIndent)
				break
			}

			// Loop back around
			continue
		}

		// Handle backslash continuation
		if ch.C == '\\' {
			// Get next character and make sure it's
			// newline
			ch = l.s.Next()
			if ch.C == common.Err {
				// Hmm, got an error
				l.pushErr(ch.Loc, ch.Val.(error))
				break
			} else if ch.C != '\n' {
				l.pushErr(ch.Loc, common.ErrDanglingBackslash)
				break
			}

			// OK, continue to the next character
			continue
		}

		// Handle the case of ".n", where n is a decimal digit
		if ch.C == '.' {
			next := l.s.Next()
			l.s.Push(next)
			if next.Class&common.CharDecDigit != 0 {
				// Suck in a number
				rNumber(l).Recognize(ch)
				continue
			}
		}

		// Apply the correct recognizer
		if ch.Class&common.CharComment != 0 {
			rComment(l).Recognize(ch)
		} else if ch.Class&common.CharDecDigit != 0 {
			rNumber(l).Recognize(ch)
		} else if ch.Class&common.CharIDStart != 0 {
			rIdent(l).Recognize(ch)
		} else if ch.Class&common.CharQuote != 0 {
			rString(l).Recognize(ch)
		} else if ch.Class == 0 {
			rOp(l).Recognize(ch)
		} else {
			l.pushErr(ch.Loc, common.ErrBadOp)
			break
		}
	}

	// If there are no tokens, return nil
	if l.tokens.Len() == 0 {
		return nil
	}

	// Pop the first element off
	elem := l.tokens.Front()
	l.tokens.Remove(elem)

	// Return the token, and save it so we know what we returned
	// last
	l.prevTok = elem.Value.(*common.Token)
	return l.prevTok
}

// Push pushes a single token back onto the lexer.  Any number of
// tokens may be pushed back.
func (l *lexer) Push(tok *common.Token) {
	// Push the token onto the queue
	l.tokens.PushFront(tok)
}
