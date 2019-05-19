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

package utils

import (
	"container/list"
	"errors"
)

// ErrNotVisitable indicates that the object cannot be visited.
var ErrNotVisitable = errors.New("object has no Children() method")

// Visitable is an interface for describing visitable nodes.
type Visitable interface {
	// Children returns the list of child nodes.
	Children() []Visitable
}

// VisitPred is a predicate function called on each of the child nodes
// of a Visitable.  It is called with a context, the Visitable being
// visited, and a boolean indicating if the Visitable is the last
// child of its parent.  It should return a context to be passed to
// the VisitPred for its children and an error.
type VisitPred func(ctxt interface{}, v Visitable, last bool) (interface{}, error)

// visitState describes the current state of the visit.
type visitState struct {
	ctxt     interface{} // The context to use for VisitPred
	children []Visitable // A list of the Visitable's children
	child    int         // The next child to be visited
}

// Visit visits every node in a Visitable in depth-first fashion,
// calling the specified VisitPred with the specified context.  If a
// VisitPred call returns an error, the error is returned.
func Visit(visitable interface{}, visit VisitPred, ctxt interface{}) error {
	// Convert to a visitable
	v, ok := visitable.(Visitable)
	if !ok {
		return ErrNotVisitable
	}

	// Visit the root Visitable first
	ctxt, err := visit(ctxt, v, true)
	if err != nil {
		return err
	}

	// Get the children
	children := v.Children()
	if len(children) == 0 {
		// No children to visit
		return nil
	}

	// Now initialize the stack
	stack := list.List{}
	stack.PushBack(&visitState{
		ctxt:     ctxt,
		children: children,
	})

	// Start working the stack
	for stack.Len() > 0 {
		// Collect the item we're working and extract the data
		elem := stack.Back()
		item := elem.Value.(*visitState)
		ctxt := item.ctxt
		v := item.children[item.child]
		last := item.child >= len(item.children)-1

		// Call the predicate on the next child
		ctxt, err := visit(ctxt, v, last)
		if err != nil {
			return err
		}

		// If it was the last child, remove the stack element
		if last {
			stack.Remove(elem)
		} else {
			// Increment the child value
			item.child++
		}

		// Get its children
		children := v.Children()
		if len(children) > 0 {
			// Adding the child to the stack to visit next
			stack.PushBack(&visitState{
				ctxt:     ctxt,
				children: children,
			})
		}
	}

	return nil
}
