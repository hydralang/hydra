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

package testutils

import (
	"reflect"

	"github.com/stretchr/testify/assert"
)

// AssertPtrEqual is a helper that asserts that two pointers are
// equal, without comparing the contents.
func AssertPtrEqual(a *assert.Assertions, expected, actual interface{}) {
	exp := reflect.ValueOf(expected)
	act := reflect.ValueOf(actual)
	a.Equal(exp.Pointer(), act.Pointer())
}

// AssertPtrNotEqual is a helper that asserts that two pointers are
// not equal, without comparing the contents.
func AssertPtrNotEqual(a *assert.Assertions, expected, actual interface{}) {
	exp := reflect.ValueOf(expected)
	act := reflect.ValueOf(actual)
	a.NotEqual(exp.Pointer(), act.Pointer())
}
