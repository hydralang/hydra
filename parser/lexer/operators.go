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
)

// opFrame describes a single element of the operator stack.  The
// operator stack allows for backtracking, by keeping track of
// characters that may need to be pushed back to the scanner.
type opFrame struct {
	ch   common.AugChar    // The character
	loc  common.Location   // The location of the whole token so far
	node *common.Operators // The operator tree node
}

// recognizeOperator is a recognizer for operators.  It should be
// called when the character is not of any other recognized class.
type recognizeOperator struct {
	l    *lexer            // The lexer
	st   list.List         // The stack of opFrame items
	node *common.Operators // Next operator tree node
}

// recogOperator constructs a recognizer for operators.
func recogOperator(l *lexer) Recognizer {
	return &recognizeOperator{
		l:    l,
		node: l.opts.Prof.Operators,
	}
}

// pushFrame pushes an operator stack frame onto the operator stack.
// If the new stack frame includes a symbol, the stack will be
// cleared.
func (r *recognizeOperator) pushFrame(ch common.AugChar, node *common.Operators) {
	// Construct the location
	loc := ch.Loc

	// If there's a symbol in the node, prepare the stack
	if node.Sym != nil && r.st.Len() > 0 {
		// Reconstruct the location
		loc = r.st.Front().Value.(*opFrame).loc.ThruEnd(ch.Loc)

		// Clear the stack
		r.st.Init()
	}

	// Push an entry onto the stack
	r.st.PushBack(&opFrame{
		ch:   ch,
		loc:  loc,
		node: node,
	})

	// Keep track of last node
	r.node = node
}

// emit pushes a token onto the lexer token stack and pushes any
// unprocessed operator characters back onto the scanner.
func (r *recognizeOperator) emit() {
	// Get the first element from the operator stack
	elem := r.st.Front()
	frame := elem.Value.(*opFrame)

	// Check if there's a token
	if frame.node.Sym == nil {
		r.l.pushErr(frame.ch.Loc, common.ErrBadOp)
		return
	}

	// Check for pairing violations
	if frame.node.Sym.Open != "" {
		if r.l.pair.Len() > 0 {
			// See if it's a match
			openElem := r.l.pair.Back()
			open := openElem.Value.(*common.Token)
			if open.Sym.Close == frame.node.Sym.Name {
				r.l.pair.Remove(openElem)
			} else {
				r.l.pushErr(frame.loc, common.ErrOpMismatch(open, frame.node.Sym))
				return
			}
		} else {
			// No opener
			r.l.pushErr(frame.loc, common.ErrNoOpen(frame.node.Sym))
			return
		}
	}

	// Emit the operator
	tok := r.l.pushTok(frame.node.Sym, frame.loc, frame.node.Sym.Name)

	// Push a pairing if necessary
	if frame.node.Sym.Close != "" {
		r.l.pair.PushBack(tok)
	}

	// Pop off the front element
	r.st.Remove(elem)

	// Push back additional characters
	for elem = r.st.Back(); elem != nil; elem = elem.Prev() {
		r.l.s.Push(elem.Value.(*opFrame).ch)
	}
}

// Recognize is called to recognize a operator.  Will be called with
// the first character, and should push zero or more tokens onto the
// lexer's tokens queue.
func (r *recognizeOperator) Recognize(ch common.AugChar) {
	// Loop over operator characters
	for ; ; ch = r.l.s.Next() {
		// Handle the error case
		if ch.C == common.Err {
			r.l.pushErr(ch.Loc, ch.Val.(error))
			return
		} else if ch.C == common.EOF || ch.Class != 0 {
			// Done processing the operator
			break
		}

		// Look up the next node
		nextNode := r.node.Next(ch.C)
		if nextNode == nil {
			// Done processing the operator
			break
		}

		// Add to the operator stack
		r.pushFrame(ch, nextNode)
	}

	// If there is no stack, we couldn't match the op
	if r.st.Len() == 0 {
		r.l.pushErr(ch.Loc, common.ErrBadOp)
		return
	}

	// Push back the character we stopped on
	r.l.s.Push(ch)

	// Emit the token
	r.emit()
}
