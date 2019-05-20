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
	"fmt"
	"strings"
)

// FilePos specifies a position within a given file.
type FilePos struct {
	L int // The line number of the position
	C int // The column number of the position
}

// Location specifies the exact range of locations of some entity.
type Location struct {
	File string  // The name of the file
	B    FilePos // The beginning of the range
	E    FilePos // The end of the range
}

// Advance advances a location in place.  The current range end
// becomes the range beginning, and the range end is the sum of the
// new range beginning and the provided offset.
func (l *Location) Advance(offset FilePos) {
	// Begin by advancing the beginning
	l.B = l.E

	// Now advance the ending; if L is incremented, C will be
	// reset to 1
	if offset.L > 0 {
		l.E.L += offset.L
		l.E.C = 1
	}
	l.E.C += offset.C
}

// AdvanceTab advances a location in place, as if by a tab character.
// The argument indicates the size of a tab stop.
func (l *Location) AdvanceTab(tabstop int) {
	l.Advance(FilePos{C: 1 + tabstop - l.E.C%tabstop})
}

// Thru creates a new Location that ranges from the beginning of this
// location to the beginning of another Location.
func (l Location) Thru(other Location) (Location, error) {
	// Location can't range across files
	if l.File != other.File {
		return Location{}, ErrSplitEntity
	}

	// Create and return the new location
	return Location{
		File: l.File,
		B:    l.B,
		E:    other.B,
	}, nil
}

// ThruEnd is similar to Thru, except that it creates a new Location
// that ranges from the beginning of this location to the ending of
// another Location.
func (l Location) ThruEnd(other Location) (Location, error) {
	// Location can't range across files
	if l.File != other.File {
		return Location{}, ErrSplitEntity
	}

	// Create and return the new location
	return Location{
		File: l.File,
		B:    l.B,
		E:    other.E,
	}, nil
}

// String constructs a string representation of the location.
func (l Location) String() string {
	text := strings.Builder{}

	// Add the beginning to the location
	text.WriteString(fmt.Sprintf("%s:%d:%d", l.File, l.B.L, l.B.C))

	// Is it split across lines or wider than one column?
	if l.B.L != l.E.L {
		text.WriteString(fmt.Sprintf("-%d:%d", l.E.L, l.E.C))
	} else if l.E.C-l.B.C > 1 {
		text.WriteString(fmt.Sprintf("-%d", l.E.C))
	}

	return text.String()
}
