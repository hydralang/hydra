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
)

func TestLocationAdvanceColumn(t *testing.T) {
	a := assert.New(t)
	loc := Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.Advance(FilePos{C: 2})

	a.Equal(FilePos{L: 3, C: 3}, loc.B)
	a.Equal(FilePos{L: 3, C: 5}, loc.E)
}

func TestLocationAdvanceLine(t *testing.T) {
	a := assert.New(t)
	loc := Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.Advance(FilePos{L: 1, C: 2})

	a.Equal(FilePos{L: 3, C: 3}, loc.B)
	a.Equal(FilePos{L: 4, C: 3}, loc.E)
}

func TestLocationAdvanceTab8(t *testing.T) {
	a := assert.New(t)
	loc := Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.AdvanceTab(8)

	a.Equal(FilePos{L: 3, C: 3}, loc.B)
	a.Equal(FilePos{L: 3, C: 9}, loc.E)
}

func TestLocationAdvanceTab4(t *testing.T) {
	a := assert.New(t)
	loc := Location{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.AdvanceTab(4)

	a.Equal(FilePos{L: 3, C: 3}, loc.B)
	a.Equal(FilePos{L: 3, C: 5}, loc.E)
}

func TestLocationThruBase(t *testing.T) {
	a := assert.New(t)
	loc1 := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := Location{File: "file", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.Thru(loc2)

	a.NoError(err)
	a.Equal("file", result.File)
	a.Equal(FilePos{L: 3, C: 2}, result.B)
	a.Equal(FilePos{L: 3, C: 5}, result.E)
}

func TestLocationThruSplit(t *testing.T) {
	a := assert.New(t)
	loc1 := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := Location{File: "other", B: FilePos{3, 5}, E: FilePos{3, 6}}

	_, err := loc1.Thru(loc2)

	a.Equal(ErrSplitEntity, err)
}

func TestLocationThruEndBase(t *testing.T) {
	a := assert.New(t)
	loc1 := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := Location{File: "file", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.ThruEnd(loc2)

	a.NoError(err)
	a.Equal("file", result.File)
	a.Equal(FilePos{L: 3, C: 2}, result.B)
	a.Equal(FilePos{L: 3, C: 6}, result.E)
}

func TestLocationThruEndSplit(t *testing.T) {
	a := assert.New(t)
	loc1 := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := Location{File: "other", B: FilePos{3, 5}, E: FilePos{3, 6}}

	_, err := loc1.ThruEnd(loc2)

	a.Equal(ErrSplitEntity, err)
}

func TestLocationString0Columns(t *testing.T) {
	a := assert.New(t)
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 2}}

	result := loc.String()

	a.Equal("file:3:2", result)
}

func TestLocationString1Column(t *testing.T) {
	a := assert.New(t)
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}

	result := loc.String()

	a.Equal("file:3:2", result)
}

func TestLocationString2Columns(t *testing.T) {
	a := assert.New(t)
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{3, 4}}

	result := loc.String()

	a.Equal("file:3:2-4", result)
}

func TestLocationString2Lines(t *testing.T) {
	a := assert.New(t)
	loc := Location{File: "file", B: FilePos{3, 2}, E: FilePos{4, 2}}

	result := loc.String()

	a.Equal("file:3:2-4:2", result)
}
