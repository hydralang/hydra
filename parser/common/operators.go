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
	"strings"
	"unicode/utf8"

	"github.com/hydralang/hydra/utils"
)

// Operators is a structure for describing an operator tree.  The
// lexer uses the operator tree to match operators, while allowing for
// backtracking; this enables selecting the longest match.
type Operators struct {
	prefix   string              // The operator prefix at this node
	Sym      *Symbol             // The operator at this node
	root     *Operators          // Root of the operator tree
	parent   *Operators          // Parent of this node
	children map[rune]*Operators // Tree node children
}

// NewOperators constructs an Operators tree with all the specified
// operators.
func NewOperators(ops ...*Symbol) *Operators {
	// Construct the root
	root := &Operators{children: map[rune]*Operators{}}

	// Add each of the symbols to it
	for _, op := range ops {
		root.Add(op)
	}

	return root
}

// doCopy constructs a new Operators tree that is a copy of this one.
// Each node is copied in turn.
func (o *Operators) doCopy(root, parent *Operators) *Operators {
	// Construct the new node
	new := &Operators{
		prefix:   o.prefix,
		Sym:      o.Sym,
		root:     root,
		parent:   parent,
		children: map[rune]*Operators{},
	}

	// If root is nil, that means we're to become the root
	if root == nil {
		root = new
	}

	// Construct each child
	for r, child := range o.children {
		new.children[r] = child.doCopy(root, new)
	}

	return new
}

// Copy constructs a copy of this Operators tree.  The copy will
// contain just the subtree rooted at this node, if this node is not
// the root.
func (o *Operators) Copy() *Operators {
	return o.doCopy(nil, nil)
}

// prune removes empty nodes of the operator tree.
func (o *Operators) prune() {
	// Step through the tree towards the root
	node := o
	for node != nil {
		// Only concerned about empty nodes
		if node.Sym != nil || (node.children != nil && len(node.children) > 0) {
			break
		}

		// Find this node in parent and delete it
		if node.parent != nil {
			r, _ := utf8.DecodeLastRuneInString(node.prefix)
			delete(node.parent.children, r)
		}

		// Step toward the root
		node = node.parent
	}
}

// Add adds an operator to the operator tree.
func (o *Operators) Add(op *Symbol) {
	// Delegate to the root
	if o.root != nil {
		o.root.Add(op)
		return
	}

	// Scan through symbol name rune by rune
	pos := 0
	node := o
	for pos < len(op.Name) {
		// Grab next rune
		r, w := utf8.DecodeRuneInString(op.Name[pos:])
		if r == utf8.RuneError && w == 1 {
			panic(utils.ErrBadRune)
		}

		// Advance the text position
		pos += w

		// Make sure the children map exists
		if node.children == nil {
			node.children = map[rune]*Operators{}
		}

		// See if the node has an entry for that rune
		if tmp, ok := node.children[r]; ok {
			node = tmp
		} else {
			// Construct a new one
			tmp = &Operators{
				prefix:   op.Name[:pos],
				root:     o,
				parent:   node,
				children: map[rune]*Operators{},
			}
			node.children[r] = tmp
			node = tmp
		}
	}

	// Is the operator already set?
	if node.Sym != nil {
		return
	}

	// Save the symbol
	node.Sym = op
}

// Remove removes an operator from the operator tree.
func (o *Operators) Remove(op *Symbol) {
	// Delegate to the root
	if o.root != nil {
		o.root.Remove(op)
		return
	}

	// Scan through symbol name rune by rune
	pos := 0
	node := o
	for pos < len(op.Name) {
		// Grab next rune
		r, w := utf8.DecodeRuneInString(op.Name[pos:])
		if r == utf8.RuneError && w == 1 {
			panic(utils.ErrBadRune)
		}

		// Advance the text position
		pos += w

		// Make sure the children map exists
		if node.children == nil {
			node.children = map[rune]*Operators{}
		}

		// Find the next node
		if tmp, ok := node.children[r]; ok {
			node = tmp
		} else {
			// Operator isn't in tree
			return
		}
	}

	// Blank the operator
	node.Sym = nil

	// Prune the node back
	node.prune()
}

// Next looks up the next node in the tree, given an operator rune.
// Returns nil if no corresponding node exists in the tree.
func (o *Operators) Next(r rune) *Operators {
	// Make sure the children map exists
	if o.children == nil {
		o.children = map[rune]*Operators{}
	}

	// See if the rune's in the tree
	child, ok := o.children[r]
	if !ok {
		return nil
	}

	return child
}

// String outputs the operator tree node as a string.
func (o *Operators) String() string {
	text := &strings.Builder{}

	// Is it a root node?
	if o.prefix != "" {
		r, _ := utf8.DecodeLastRuneInString(o.prefix)
		text.WriteString(fmt.Sprintf("'%c' (%d): ", r, r))
	}

	// Add the operator
	if o.Sym != nil {
		text.WriteString(fmt.Sprintf("Operator \"%s\"", o.Sym))
	} else {
		text.WriteString("(node)")
	}

	return text.String()
}

// Children implements the utils.Visitable interface, allowing an
// operator tree to be visualized using utils.Visualize().
func (o *Operators) Children() []utils.Visitable {
	// Make sure the children map exists
	if o.children == nil {
		o.children = map[rune]*Operators{}
	}

	// Construct the returned visitables
	result := make([]utils.Visitable, len(o.children))

	// Add the children to the list
	i := 0
	for _, child := range o.children {
		result[i] = child
		i++
	}

	return result
}
