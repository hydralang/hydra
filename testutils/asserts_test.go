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
	"testing"

	"github.com/stretchr/testify/assert"
)

type tStruct struct {
	Text string
}

func TestAssertPtrEqual(t *testing.T) {
	a := assert.New(t)
	obj1 := &tStruct{"obj1"}
	obj2 := obj1

	AssertPtrEqual(a, obj1, obj2)
}

func TestAssertPtrNotEqual(t *testing.T) {
	a := assert.New(t)
	obj1 := &tStruct{"obj1"}
	obj2 := &tStruct{"obj2"}

	AssertPtrNotEqual(a, obj1, obj2)
}
