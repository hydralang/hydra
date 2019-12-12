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
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func equalFunc(a *assert.Assertions, expected, actual interface{}) {
	exp := reflect.ValueOf(expected).Pointer()
	act := reflect.ValueOf(actual).Pointer()
	a.Equal(exp, act)
}

type visStringTest struct {
	value string
}

func (v *visStringTest) String() string {
	return v.value
}

func (v *visStringTest) Children() []Visitable {
	return []Visitable{}
}

type visTest struct{}

func (v *visTest) Children() []Visitable {
	return []Visitable{}
}

func visPredString(v Visitable) (string, error) {
	return "constant", nil
}

func visPredError(v Visitable) (string, error) {
	return "", assert.AnError
}

func TestVisPredRoot(t *testing.T) {
	a := assert.New(t)
	ctxt := &visCtxt{
		prof:   VisASCII,
		vis:    visStringPred,
		buf:    &strings.Builder{},
		prefix: "",
	}
	obj := &visStringTest{value: "node"}

	nextCtxt, err := visPred(ctxt, obj, false)

	a.NoError(err)
	a.NotNil(nextCtxt)
	nextVisCtxt, ok := nextCtxt.(*visCtxt)
	a.True(ok)
	a.Equal(ctxt.prof, nextVisCtxt.prof)
	equalFunc(a, ctxt.vis, nextVisCtxt.vis)
	a.Equal(ctxt.buf, nextVisCtxt.buf)
	a.Equal("-- node\n", nextVisCtxt.buf.String())
	a.Equal("   ", nextVisCtxt.prefix)
}

func TestVisPredLast(t *testing.T) {
	a := assert.New(t)
	ctxt := &visCtxt{
		prof:   VisASCII,
		vis:    visStringPred,
		buf:    &strings.Builder{},
		prefix: "  ",
	}
	obj := &visStringTest{value: "node"}

	nextCtxt, err := visPred(ctxt, obj, true)

	a.NoError(err)
	a.NotNil(nextCtxt)
	nextVisCtxt, ok := nextCtxt.(*visCtxt)
	a.True(ok)
	a.Equal(ctxt.prof, nextVisCtxt.prof)
	equalFunc(a, ctxt.vis, nextVisCtxt.vis)
	a.Equal(ctxt.buf, nextVisCtxt.buf)
	a.Equal("  `- node\n", nextVisCtxt.buf.String())
	a.Equal("     ", nextVisCtxt.prefix)
}

func TestVisPredMiddle(t *testing.T) {
	a := assert.New(t)
	ctxt := &visCtxt{
		prof:   VisASCII,
		vis:    visStringPred,
		buf:    &strings.Builder{},
		prefix: "  ",
	}
	obj := &visStringTest{value: "node"}

	nextCtxt, err := visPred(ctxt, obj, false)

	a.NoError(err)
	a.NotNil(nextCtxt)
	nextVisCtxt, ok := nextCtxt.(*visCtxt)
	a.True(ok)
	a.Equal(ctxt.prof, nextVisCtxt.prof)
	equalFunc(a, ctxt.vis, nextVisCtxt.vis)
	a.Equal(ctxt.buf, nextVisCtxt.buf)
	a.Equal("  +- node\n", nextVisCtxt.buf.String())
	a.Equal("  |  ", nextVisCtxt.prefix)
}

func TestVisPredError(t *testing.T) {
	a := assert.New(t)
	ctxt := &visCtxt{
		prof:   VisASCII,
		vis:    visPredError,
		buf:    &strings.Builder{},
		prefix: "",
	}
	obj := &visStringTest{value: "node"}

	nextCtxt, err := visPred(ctxt, obj, false)

	a.Equal(assert.AnError, err)
	a.Nil(nextCtxt)
}

func TestVisStringPredStringer(t *testing.T) {
	a := assert.New(t)
	obj := &visStringTest{value: "node"}

	str, err := visStringPred(obj)

	a.NoError(err)
	a.Equal("node", str)
}

func TestVisStringNonStringer(t *testing.T) {
	a := assert.New(t)
	obj := &visTest{}

	str, err := visStringPred(obj)

	a.Equal(ErrNoVis, err)
	a.Equal("", str)
}

func TestVisProfile(t *testing.T) {
	a := assert.New(t)
	opts := &visCtxt{}

	opt := VisProfile(VisASCII)
	opt(opts)

	a.Equal(VisASCII, opts.prof)
}

func TestVisPredicate(t *testing.T) {
	a := assert.New(t)
	opts := &visCtxt{}

	opt := VisPredicate(visPredError)
	opt(opts)

	equalFunc(a, visPredError, opts.vis)
}

func TestVisualizeBase(t *testing.T) {
	a := assert.New(t)
	obj := &visStringTest{value: "node"}

	result, err := Visualize(obj)

	a.NoError(err)
	a.Equal("\u2500\u2500 node\n", result)
}

func TestVisualizeAlternate(t *testing.T) {
	a := assert.New(t)
	obj := &visStringTest{value: "node"}

	result, err := Visualize(
		obj, VisProfile(VisASCII), VisPredicate(visPredString),
	)

	a.NoError(err)
	a.Equal("-- constant\n", result)
}

func TestVisualizeError(t *testing.T) {
	a := assert.New(t)
	obj := &visStringTest{value: "node"}

	result, err := Visualize(obj, VisPredicate(visPredError))

	a.Equal(assert.AnError, err)
	a.Equal("", result)
}

type mockAnnotatedStringer struct {
	mock.Mock
}

func (m *mockAnnotatedStringer) Children() []Visitable {
	args := m.MethodCalled("Children")

	return args.Get(0).([]Visitable)
}

func (m *mockAnnotatedStringer) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

func TestAnnotatedChildren(t *testing.T) {
	a := assert.New(t)
	children := []Visitable{
		&mockAnnotatedStringer{},
		&mockAnnotatedStringer{},
	}
	wrapped := &mockAnnotatedStringer{}
	wrapped.On("Children").Return(children)
	obj := &Annotated{
		Wrapped: wrapped,
	}

	result := obj.Children()

	a.Equal(children, result)
	wrapped.AssertExpectations(t)
}

func TestAnnotatedString(t *testing.T) {
	a := assert.New(t)
	wrapped := &mockAnnotatedStringer{}
	wrapped.On("String").Return("ren")
	obj := &Annotated{
		Wrapped:    wrapped,
		Annotation: "child",
	}

	result := obj.String()

	a.Equal("children", result)
	wrapped.AssertExpectations(t)
}
