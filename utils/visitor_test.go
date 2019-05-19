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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type visitTest struct {
	children []Visitable
	ctxt     int
	last     bool
	visited  bool
	err      error
}

func (v *visitTest) Children() []Visitable {
	return v.children
}

func visitPredTest(ctxt interface{}, v Visitable, last bool) (interface{}, error) {
	obj := v.(*visitTest)

	obj.ctxt = ctxt.(int)
	obj.last = last
	obj.visited = true

	if obj.err != nil {
		return nil, obj.err
	}
	return obj.ctxt + 1, nil
}

func TestVisitUnvisitable(t *testing.T) {
	a := assert.New(t)
	obj := &mock.Mock{}

	err := Visit(obj, visitPredTest, 0)

	a.Equal(ErrNotVisitable, err)
}

func TestVisitBase(t *testing.T) {
	a := assert.New(t)
	obj := &visitTest{}

	err := Visit(obj, visitPredTest, 0)

	a.NoError(err)
	a.Equal(0, obj.ctxt)
	a.True(obj.last)
	a.True(obj.visited)
}

func TestVisitRootError(t *testing.T) {
	a := assert.New(t)
	obj := &visitTest{
		err: assert.AnError,
	}

	err := Visit(obj, visitPredTest, 0)

	a.Equal(assert.AnError, err)
	a.Equal(0, obj.ctxt)
	a.True(obj.last)
	a.True(obj.visited)
}

func TestVisitChildren(t *testing.T) {
	a := assert.New(t)
	l3c1 := &visitTest{}
	l3c2 := &visitTest{}
	l3c3 := &visitTest{}
	l3c4 := &visitTest{}
	l2c1 := &visitTest{
		children: []Visitable{l3c1, l3c2},
	}
	l2c2 := &visitTest{
		children: []Visitable{l3c3, l3c4},
	}
	l2c3 := &visitTest{}
	root := &visitTest{
		children: []Visitable{l2c1, l2c2, l2c3},
	}

	err := Visit(root, visitPredTest, 0)

	a.NoError(err)
	a.Equal(0, root.ctxt)
	a.True(root.last)
	a.True(root.visited)
	a.Equal(1, l2c1.ctxt)
	a.False(l2c1.last)
	a.True(l2c1.visited)
	a.Equal(1, l2c2.ctxt)
	a.False(l2c2.last)
	a.True(l2c2.visited)
	a.Equal(1, l2c3.ctxt)
	a.True(l2c3.last)
	a.True(l2c3.visited)
	a.Equal(2, l3c1.ctxt)
	a.False(l3c1.last)
	a.True(l3c1.visited)
	a.Equal(2, l3c2.ctxt)
	a.True(l3c2.last)
	a.True(l3c2.visited)
	a.Equal(2, l3c3.ctxt)
	a.False(l3c3.last)
	a.True(l3c3.visited)
	a.Equal(2, l3c4.ctxt)
	a.True(l3c4.last)
	a.True(l3c4.visited)
}

func TestVisitChildrenError(t *testing.T) {
	a := assert.New(t)
	l3c1 := &visitTest{}
	l3c2 := &visitTest{}
	l3c3 := &visitTest{
		err: assert.AnError,
	}
	l3c4 := &visitTest{}
	l2c1 := &visitTest{
		children: []Visitable{l3c1, l3c2},
	}
	l2c2 := &visitTest{
		children: []Visitable{l3c3, l3c4},
	}
	l2c3 := &visitTest{}
	root := &visitTest{
		children: []Visitable{l2c1, l2c2, l2c3},
	}

	err := Visit(root, visitPredTest, 0)

	a.Equal(assert.AnError, err)
	a.Equal(0, root.ctxt)
	a.True(root.last)
	a.True(root.visited)
	a.Equal(1, l2c1.ctxt)
	a.False(l2c1.last)
	a.True(l2c1.visited)
	a.Equal(1, l2c2.ctxt)
	a.False(l2c2.last)
	a.True(l2c2.visited)
	a.Equal(0, l2c3.ctxt)
	a.False(l2c3.last)
	a.False(l2c3.visited)
	a.Equal(2, l3c1.ctxt)
	a.False(l3c1.last)
	a.True(l3c1.visited)
	a.Equal(2, l3c2.ctxt)
	a.True(l3c2.last)
	a.True(l3c2.visited)
	a.Equal(2, l3c3.ctxt)
	a.False(l3c3.last)
	a.True(l3c3.visited)
	a.Equal(0, l3c4.ctxt)
	a.False(l3c4.last)
	a.False(l3c4.visited)
}
