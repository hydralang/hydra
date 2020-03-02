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

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/testutils"
)

func TestKeywordsCopy(t *testing.T) {
	a := assert.New(t)

	result := testKeywords.Copy()

	a.Equal(testKeywords, result)
	testutils.AssertPtrNotEqual(a, testKeywords, result)
}

func TestKeywordsAddPresent(t *testing.T) {
	a := assert.New(t)
	obj := Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}

	obj.Add(&Symbol{Name: "kw1"})

	a.Equal(Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}, obj)
}

func TestKeywordsAddAbsent(t *testing.T) {
	a := assert.New(t)
	obj := Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}

	obj.Add(&Symbol{Name: "kw3"})

	a.Equal(Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
		"kw3": &Symbol{Name: "kw3"},
	}, obj)
}

func TestKeywordsRemovePresent(t *testing.T) {
	a := assert.New(t)
	obj := Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}

	obj.Remove(&Symbol{Name: "kw1"})

	a.Equal(Keywords{
		"kw2": &Symbol{Name: "kw2"},
	}, obj)
}

func TestKeywordsRemoveAbsent(t *testing.T) {
	a := assert.New(t)
	obj := Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}

	obj.Remove(&Symbol{Name: "kw3"})

	a.Equal(Keywords{
		"kw1": &Symbol{Name: "kw1"},
		"kw2": &Symbol{Name: "kw2"},
	}, obj)
}
