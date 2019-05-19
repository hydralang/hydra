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
)

// FlagSet is an interface representing sets of named flags.
type FlagSet interface {
	// FlagName retrieves the name for a flag corresponding to the
	// specified bit.
	FlagName(b uint) string

	// Flags retrieves a list of flags corresponding to the
	// specified bit set.  The flags will be in order from lowest
	// set bit to highest.
	Flags(flags interface{}) []string
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

// Flags retrieves a list of flags corresponding to the specified bit
// set.  The flags will be in order from lowest set bit to highest.
func (fs FlagSet8) Flags(flags interface{}) []string {
	// Convert flags to the correct type
	mask := flags.(uint8)

	// Initialize the list of flag names
	flagNames := make([]string, bits.OnesCount8(mask))

	// Collect all the names
	idx := 0
	lo := uint(bits.TrailingZeros8(mask))
	hi := uint(bits.Len8(mask))
	for i := lo + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
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

// Flags retrieves a list of flags corresponding to the specified bit
// set.  The flags will be in order from lowest set bit to highest.
func (fs FlagSet16) Flags(flags interface{}) []string {
	// Convert flags to the correct type
	mask := flags.(uint16)

	// Initialize the list of flag names
	flagNames := make([]string, bits.OnesCount16(mask))

	// Collect all the names
	idx := 0
	lo := uint(bits.TrailingZeros16(mask))
	hi := uint(bits.Len16(mask))
	for i := lo + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
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

// Flags retrieves a list of flags corresponding to the specified bit
// set.  The flags will be in order from lowest set bit to highest.
func (fs FlagSet32) Flags(flags interface{}) []string {
	// Convert flags to the correct type
	mask := flags.(uint32)

	// Initialize the list of flag names
	flagNames := make([]string, bits.OnesCount32(mask))

	// Collect all the names
	idx := 0
	lo := uint(bits.TrailingZeros32(mask))
	hi := uint(bits.Len32(mask))
	for i := lo + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
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

// Flags retrieves a list of flags corresponding to the specified bit
// set.  The flags will be in order from lowest set bit to highest.
func (fs FlagSet64) Flags(flags interface{}) []string {
	// Convert flags to the correct type
	mask := flags.(uint64)

	// Initialize the list of flag names
	flagNames := make([]string, bits.OnesCount64(mask))

	// Collect all the names
	idx := 0
	lo := uint(bits.TrailingZeros64(mask))
	hi := uint(bits.Len64(mask))
	for i := lo + 1; i <= hi; i++ {
		flagNames[idx] = fs.FlagName(i)
		idx++
	}

	return flagNames
}
