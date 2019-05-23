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
	"testing"

	"github.com/hydralang/hydra/testutils"
	"github.com/hydralang/hydra/utils"
	"github.com/stretchr/testify/assert"
)

func TestOperatorsImplementsVisitable(t *testing.T) {
	assert.Implements(t, (*utils.Visitable)(nil), &Operators{})
}

func TestNewOperators(t *testing.T) {
	a := assert.New(t)
	sym1 := &Symbol{Name: "=="}
	sym2 := &Symbol{Name: "="}

	obj := NewOperators(sym1, sym2)

	a.Equal("", obj.prefix)
	a.Nil(obj.Sym)
	a.Nil(obj.root)
	a.Nil(obj.parent)
	a.Equal(1, len(obj.children))
	a.Contains(obj.children, '=')
	op2 := obj.children['=']
	a.Equal("=", op2.prefix)
	a.Equal(sym2, op2.Sym)
	a.Equal(obj, op2.root)
	a.Equal(obj, op2.parent)
	a.Equal(1, len(op2.children))
	a.Contains(op2.children, '=')
	op1 := op2.children['=']
	a.Equal("==", op1.prefix)
	a.Equal(sym1, op1.Sym)
	a.Equal(obj, op1.root)
	a.Equal(op2, op1.parent)
	a.Equal(0, len(op1.children))
}

func TestOperatorsCopy(t *testing.T) {
	a := assert.New(t)
	sym1 := &Symbol{Name: "<"}
	sym2 := &Symbol{Name: ">"}
	sym3 := &Symbol{Name: "<<"}
	sym4 := &Symbol{Name: ">>"}
	tree := NewOperators(sym1, sym2, sym3, sym4)

	result := tree.Copy()

	a.Equal(tree, result)
	testutils.AssertPtrNotEqual(a, tree, result)
	for r, child := range result.children {
		a.Equal(tree.children[r], child)
		testutils.AssertPtrNotEqual(a, tree.children[r], child)
		testutils.AssertPtrEqual(a, result, child.root)
		testutils.AssertPtrEqual(a, result, child.parent)
		for rr, grandchild := range child.children {
			a.Equal(tree.children[r].children[rr], grandchild)
			testutils.AssertPtrNotEqual(a, tree.children[r].children[rr], grandchild)
			testutils.AssertPtrEqual(a, result, grandchild.root)
			testutils.AssertPtrEqual(a, result.children[r], grandchild.parent)
		}
	}
}

func TestOperatorsPruneEmptyTree(t *testing.T) {
	a := assert.New(t)
	tree := &Operators{}

	tree.prune()

	a.Nil(tree.Sym)
	a.Nil(tree.parent)
	a.Nil(tree.children)
}

func TestOperatorsPruneOpNode(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "="}
	node := &Operators{
		prefix: "=",
		Sym:    sym,
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.parent = tree

	node.prune()

	a.Equal(sym, node.Sym)
	a.Equal(tree, node.parent)
	a.Nil(node.children)
	a.Nil(tree.Sym)
	a.Nil(tree.parent)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
}

func TestOperatorsPruneChildrenNode(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	child := &Operators{
		prefix: "==",
		Sym:    sym,
	}
	node := &Operators{
		prefix: "=",
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.parent = tree

	node.prune()

	a.Equal(sym, child.Sym)
	a.Equal(node, child.parent)
	a.Nil(child.children)
	a.Nil(node.Sym)
	a.Equal(tree, node.parent)
	a.Equal(map[rune]*Operators{
		'=': child,
	}, node.children)
	a.Nil(tree.Sym)
	a.Nil(tree.parent)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
}

func TestOperatorsPruneChildNode(t *testing.T) {
	a := assert.New(t)
	node := &Operators{
		prefix: "=",
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.parent = tree

	node.prune()

	a.Nil(node.Sym)
	a.Equal(tree, node.parent)
	a.Nil(node.children)
	a.Nil(tree.Sym)
	a.Nil(tree.parent)
	a.Equal(map[rune]*Operators{}, tree.children)
}

func TestOperatorsPruneRecurse(t *testing.T) {
	a := assert.New(t)
	child := &Operators{
		prefix: "==",
	}
	node := &Operators{
		prefix: "=",
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.parent = tree

	child.prune()

	a.Nil(child.Sym)
	a.Equal(node, child.parent)
	a.Nil(child.children)
	a.Nil(node.Sym)
	a.Equal(tree, node.parent)
	a.Equal(map[rune]*Operators{}, node.children)
	a.Nil(tree.Sym)
	a.Nil(tree.parent)
	a.Equal(map[rune]*Operators{}, tree.children)
}

func TestOperatorsAddAbsent(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	node := &Operators{
		prefix: "=",
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.root = tree
	node.parent = tree

	tree.Add(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Nil(node.Sym)
	a.NotNil(node.children)
	a.Contains(node.children, '=')
	child := node.children['=']
	a.Equal("==", child.prefix)
	a.Equal(sym, child.Sym)
	a.Equal(map[rune]*Operators{}, child.children)
}

func TestOperatorsAddPresent(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	child := &Operators{
		prefix: "==",
		Sym:    &Symbol{Name: "!="},
	}
	node := &Operators{
		prefix: "=",
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	child.root = tree
	node.root = tree
	node.parent = tree

	tree.Add(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Nil(node.Sym)
	a.Equal(map[rune]*Operators{
		'=': child,
	}, node.children)
	a.Equal("==", child.prefix)
	a.NotEqual(sym, child.Sym)
}

func TestOperatorsAddBadRune(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: string([]byte{'=', '\xff'})}
	tree := &Operators{}

	a.PanicsWithValue(ErrBadRune, func() { tree.Add(sym) })

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
}

func TestOperatorsAddDelegate(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	node := &Operators{
		prefix: "=",
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.root = tree
	node.parent = tree

	node.Add(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Nil(node.Sym)
	a.NotNil(node.children)
	a.Contains(node.children, '=')
	child := node.children['=']
	a.Equal("==", child.prefix)
	a.Equal(sym, child.Sym)
	a.Equal(map[rune]*Operators{}, child.children)
}

func TestOperatorsRemoveAbsent(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	node := &Operators{
		prefix: "=",
		Sym:    &Symbol{Name: "="},
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.root = tree
	node.parent = tree

	tree.Remove(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Equal(&Symbol{Name: "="}, node.Sym)
	a.NotNil(node.children)
	a.Equal(map[rune]*Operators{}, node.children)
}

func TestOperatorsRemovePresent(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	child := &Operators{
		prefix: "==",
		Sym:    sym,
	}
	node := &Operators{
		prefix: "=",
		Sym:    &Symbol{Name: "="},
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	child.root = tree
	node.root = tree
	node.parent = tree

	tree.Remove(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Equal(&Symbol{Name: "="}, node.Sym)
	a.NotNil(node.children)
	a.Equal(map[rune]*Operators{}, node.children)
}

func TestOperatorsRemovePrune(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	child := &Operators{
		prefix: "==",
		Sym:    sym,
	}
	node := &Operators{
		prefix: "=",
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	child.root = tree
	node.root = tree
	node.parent = tree

	tree.Remove(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{}, tree.children)
}

func TestRemoveBadRune(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: string([]byte{'=', '\xff'})}
	node := &Operators{
		prefix: "=",
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	node.root = tree
	node.parent = tree

	a.PanicsWithValue(ErrBadRune, func() { tree.Remove(sym) })

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
}

func TestOperatorsRemoveDelegate(t *testing.T) {
	a := assert.New(t)
	sym := &Symbol{Name: "=="}
	child := &Operators{
		prefix: "==",
		Sym:    sym,
	}
	node := &Operators{
		prefix: "=",
		Sym:    &Symbol{Name: "="},
		children: map[rune]*Operators{
			'=': child,
		},
	}
	child.parent = node
	tree := &Operators{
		children: map[rune]*Operators{
			'=': node,
		},
	}
	child.root = tree
	node.root = tree
	node.parent = tree

	node.Remove(sym)

	a.Equal("", tree.prefix)
	a.Nil(tree.Sym)
	a.Equal(map[rune]*Operators{
		'=': node,
	}, tree.children)
	a.Equal("=", node.prefix)
	a.Equal(&Symbol{Name: "="}, node.Sym)
	a.NotNil(node.children)
	a.Equal(map[rune]*Operators{}, node.children)
}

func TestOperatorsNextPresent(t *testing.T) {
	a := assert.New(t)
	child := &Operators{}
	node := &Operators{
		children: map[rune]*Operators{
			'=': child,
		},
	}

	next := node.Next('=')

	a.NotNil(next)
	a.Equal(child, next)
	a.Equal(map[rune]*Operators{
		'=': child,
	}, node.children)
}

func TestOperatorsNextAbsent(t *testing.T) {
	a := assert.New(t)
	node := &Operators{
		children: map[rune]*Operators{},
	}

	next := node.Next('=')

	a.Nil(next)
	a.Equal(map[rune]*Operators{}, node.children)
}

func TestOperatorsNextNoChildren(t *testing.T) {
	a := assert.New(t)
	node := &Operators{}

	next := node.Next('=')

	a.Nil(next)
	a.NotNil(node.children)
	a.Equal(map[rune]*Operators{}, node.children)
}

func TestOperatorsStringRoot(t *testing.T) {
	a := assert.New(t)
	node := &Operators{}

	result := node.String()

	a.Equal("(node)", result)
}

func TestOperatorsStringNode(t *testing.T) {
	a := assert.New(t)
	node := &Operators{
		prefix: "-=",
	}

	result := node.String()

	a.Equal("'=' (61): (node)", result)
}

func TestOperatorsStringSymbol(t *testing.T) {
	a := assert.New(t)
	node := &Operators{
		Sym: &Symbol{Name: "=>"},
	}

	result := node.String()

	a.Equal("Operator \"=>\"", result)
}

func TestOperatorsChildrenBase(t *testing.T) {
	a := assert.New(t)
	node1 := &Operators{
		Sym: &Symbol{Name: "=>"},
	}
	node2 := &Operators{
		Sym: &Symbol{Name: "=<"},
	}
	tree := &Operators{
		children: map[rune]*Operators{
			'>': node1,
			'<': node2,
		},
	}

	result := tree.Children()

	a.Equal(2, len(result))
	a.Contains(result, node1)
	a.Contains(result, node2)
}

func TestOperatorsChildrenCreateChildren(t *testing.T) {
	a := assert.New(t)
	tree := &Operators{}

	result := tree.Children()

	a.Equal(0, len(result))
	a.NotNil(tree.children)
	a.Equal(map[rune]*Operators{}, tree.children)
}
