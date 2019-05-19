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
	"errors"
	"fmt"
	"strings"
)

// ErrNoVis indicates that a node cannot be visualized.
var ErrNoVis = errors.New("node has no String() method")

// VisProf represents a visual profile for visualizing a tree.
type VisProf struct {
	start  rune // Joiner character to start a tree at the root
	last   rune // Joiner character for last child
	branch rune // Joiner character for branch
	skip   rune // Character indicating node not a branch
	into   rune // Character leading into a node
}

// Pre-configured visual profiles for trees.
var (
	VisASCII = VisProf{
		start:  '-',
		last:   '`',
		branch: '+',
		skip:   '|',
		into:   '-',
	}
	VisRounded = VisProf{
		start:  '\u2500',
		last:   '\u2570',
		branch: '\u251c',
		skip:   '\u2502',
		into:   '\u2500',
	}
	VisSquare = VisProf{
		start:  '\u2500',
		last:   '\u2514',
		branch: '\u251c',
		skip:   '\u2502',
		into:   '\u2500',
	}
)

// VisPred is a predicate function called on a given Visitable to
// produce a string representation of a single node in a tree.
type VisPred func(v Visitable) (string, error)

// visCtxt is a context to pass to the visPred function.
type visCtxt struct {
	prof   VisProf          // The visual profile to utilize
	vis    VisPred          // The predicate function to visualize the node
	buf    *strings.Builder // Buffer for the output
	prefix string           // The prefix for this level of the tree
}

// visPred is the visualizer predicate function.
func visPred(ctxt interface{}, v Visitable, last bool) (interface{}, error) {
	// Get the actual visualization context
	vis := ctxt.(*visCtxt)

	// Select the joiner and construct the next context
	var joiner rune
	var nextVis *visCtxt
	if vis.prefix == "" {
		joiner = vis.prof.start
		nextVis = &visCtxt{
			prof:   vis.prof,
			vis:    vis.vis,
			buf:    vis.buf,
			prefix: "   ",
		}
	} else if last {
		joiner = vis.prof.last
		nextVis = &visCtxt{
			prof:   vis.prof,
			vis:    vis.vis,
			buf:    vis.buf,
			prefix: fmt.Sprintf("%s   ", vis.prefix),
		}
	} else {
		joiner = vis.prof.branch
		nextVis = &visCtxt{
			prof:   vis.prof,
			vis:    vis.vis,
			buf:    vis.buf,
			prefix: fmt.Sprintf("%s%c  ", vis.prefix, vis.prof.skip),
		}
	}

	// Add the visitable to the buffer
	nodeVis, err := vis.vis(v)
	if err != nil {
		return nil, err
	}
	vis.buf.WriteString(fmt.Sprintf(
		"%s%c%c %s\n",
		vis.prefix, joiner, vis.prof.into, nodeVis,
	))

	return nextVis, nil
}

// visStringPred is a default predicate for visualizing a node in a
// tree.  It simply calls the node's String() method.
func visStringPred(v Visitable) (string, error) {
	// Convert visitable to a stringer
	obj, ok := v.(fmt.Stringer)
	if !ok {
		return "", ErrNoVis
	}

	return obj.String(), nil
}

// VisOption is an option for the tree visualizer.
type VisOption func(opts *visCtxt)

// VisProfile sets the visual profile to use for the visualization.
func VisProfile(prof VisProf) VisOption {
	return func(opts *visCtxt) {
		opts.prof = prof
	}
}

// VisPredicate sets the predicate for visualizing a tree node.
func VisPredicate(vis VisPred) VisOption {
	return func(opts *visCtxt) {
		opts.vis = vis
	}
}

// Visualize constructs a visualization of a tree implementing the
// Visitable interface.
func Visualize(v interface{}, opts ...VisOption) (string, error) {
	// Initialize the visualization context
	ctxt := &visCtxt{
		prof:   VisRounded,
		vis:    visStringPred,
		buf:    &strings.Builder{},
		prefix: "",
	}

	// Apply options
	for _, opt := range opts {
		opt(ctxt)
	}

	// Visualize the tree
	err := Visit(v, visPred, ctxt)
	if err != nil {
		return "", err
	}

	// Return the visualization
	return ctxt.buf.String(), nil
}
