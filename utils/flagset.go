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
	"fmt"
	"math/bits"
	"strings"
)

// flagSet is an interface representing sets of named flags.
type flagSet interface {
	// Retrieves the name for a flag corresponding to the
	// specified bit.
	FlagName(b uint) string
}

// FlagSet8 is a mapping from flags to flag names.
type FlagSet8 map[uint8]string

// FlagName retrieves the name for a flag corresponding to the
// specified bit.
func (fs FlagSet8) FlagName(b uint) string {
	var flag uint8 = 1
	flag <<= b - 1

	name, ok := fs[flag]
	if !ok {
		return fmt.Sprintf("%d (0x%02x)", b, flag)
	}
	return name
}

// FlagSet16 is a mapping from flags to flag names.
type FlagSet16 map[uint16]string

// FlagName retrieves the name for a flag corresponding to the
// specified bit.
func (fs FlagSet16) FlagName(b uint) string {
	var flag uint16 = 1
	flag <<= b - 1

	name, ok := fs[flag]
	if !ok {
		return fmt.Sprintf("%d (0x%04x)", b, flag)
	}
	return name
}

// FlagSet32 is a mapping from flags to flag names.
type FlagSet32 map[uint32]string

// FlagName retrieves the name for a flag corresponding to the
// specified bit.
func (fs FlagSet32) FlagName(b uint) string {
	var flag uint32 = 1
	flag <<= b - 1

	name, ok := fs[flag]
	if !ok {
		return fmt.Sprintf("%d (0x%08x)", b, flag)
	}
	return name
}

// FlagSet64 is a mapping from flags to flag names.
type FlagSet64 map[uint64]string

// FlagName retrieves the name for a flag corresponding to the
// specified bit.
func (fs FlagSet64) FlagName(b uint) string {
	var flag uint64 = 1
	flag <<= b - 1

	name, ok := fs[flag]
	if !ok {
		return fmt.Sprintf("%d (0x%016x)", b, flag)
	}
	return name
}

// flags8 assembles the flag names for 8-bit flag sets.
func flags8(fs FlagSet8, flags uint8) []string {
	flagNames := make([]string, bits.OnesCount8(flags))

	// Collect all the names
	idx := 0
	var low uint = uint(bits.TrailingZeros8(flags))
	var hi uint = uint(bits.Len8(flags))
	for i := low + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
}

// flags16 assembles the flag names for 16-bit flag sets.
func flags16(fs FlagSet16, flags uint16) []string {
	flagNames := make([]string, bits.OnesCount16(flags))

	// Collect all the names
	idx := 0
	var low uint = uint(bits.TrailingZeros16(flags))
	var hi uint = uint(bits.Len16(flags))
	for i := low + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
}

// flags32 assembles the flag names for 32-bit flag sets.
func flags32(fs FlagSet32, flags uint32) []string {
	flagNames := make([]string, bits.OnesCount32(flags))

	// Collect all the names
	idx := 0
	var low uint = uint(bits.TrailingZeros32(flags))
	var hi uint = uint(bits.Len32(flags))
	for i := low + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
}

// flags64 assembles the flag names for 64-bit flag sets.
func flags64(fs FlagSet64, flags uint64) []string {
	flagNames := make([]string, bits.OnesCount64(flags))

	// Collect all the names
	idx := 0
	var low uint = uint(bits.TrailingZeros64(flags))
	var hi uint = uint(bits.Len64(flags))
	for i := low + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
}

// Flags constructs a string representing the flags set, joined by the
// specified joiner.
func Flags(fs flagSet, flags interface{}, joiner string) string {
	var flagNames []string

	// Delegate to the appropriate flag compiler
	switch set := fs.(type) {
	case FlagSet8:
		flagNames = flags8(set, flags.(uint8))
	case FlagSet16:
		flagNames = flags16(set, flags.(uint16))
	case FlagSet32:
		flagNames = flags32(set, flags.(uint32))
	case FlagSet64:
		flagNames = flags64(set, flags.(uint64))
	}

	return strings.Join(flagNames, joiner)
}
