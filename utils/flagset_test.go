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
)

const (
	FLAG8F1 uint8 = 1 << iota
	FLAG8F2
	FLAG8F3
)

var names8 = FlagSet8{
	FLAG8F1: "FLAG8F1",
	FLAG8F2: "FLAG8F2",
	FLAG8F3: "FLAG8F3",
}

func TestFlags8(t *testing.T) {
	a := assert.New(t)
	flags := FLAG8F2 | FLAG8F3 | (1 << 3)

	result := Flags(names8, flags, ", ")

	a.Equal("FLAG8F2, FLAG8F3, 4 (0x08)", result)
}

const (
	FLAG16F1 uint16 = 1 << iota
	FLAG16F2
	FLAG16F3
)

var names16 = FlagSet16{
	FLAG16F1: "FLAG16F1",
	FLAG16F2: "FLAG16F2",
	FLAG16F3: "FLAG16F3",
}

func TestFlags16(t *testing.T) {
	a := assert.New(t)
	flags := FLAG16F2 | FLAG16F3 | (1 << 3)

	result := Flags(names16, flags, ", ")

	a.Equal("FLAG16F2, FLAG16F3, 4 (0x0008)", result)
}

const (
	FLAG32F1 uint32 = 1 << iota
	FLAG32F2
	FLAG32F3
)

var names32 = FlagSet32{
	FLAG32F1: "FLAG32F1",
	FLAG32F2: "FLAG32F2",
	FLAG32F3: "FLAG32F3",
}

func TestFlags32(t *testing.T) {
	a := assert.New(t)
	flags := FLAG32F2 | FLAG32F3 | (1 << 3)

	result := Flags(names32, flags, ", ")

	a.Equal("FLAG32F2, FLAG32F3, 4 (0x00000008)", result)
}

const (
	FLAG64F1 uint64 = 1 << iota
	FLAG64F2
	FLAG64F3
)

var names64 = FlagSet64{
	FLAG64F1: "FLAG64F1",
	FLAG64F2: "FLAG64F2",
	FLAG64F3: "FLAG64F3",
}

func TestFlags64(t *testing.T) {
	a := assert.New(t)
	flags := FLAG64F2 | FLAG64F3 | (1 << 3)

	result := Flags(names64, flags, ", ")

	a.Equal("FLAG64F2, FLAG64F3, 4 (0x0000000000000008)", result)
}
