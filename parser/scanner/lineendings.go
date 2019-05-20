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

package scanner

import "github.com/hydralang/hydra/parser/common"

// lineEnding is a function that processes line ending characters.
type lineEnding func(ch rune) rune

// leUnknown handles line endings when the style is not yet known.
func (s *scanner) leUnknown(ch rune) rune {
	// If it's a newline, that's easy
	if ch == '\n' {
		s.le = s.leNewline
		return '\n'
	}

	// Peek at the next character
	ch, err := s.nextChar()
	if ch == common.EOF || err != nil {
		// EOF; assume '\r' was end-of-line
		s.le = s.leCarriage

		// Push back the error
		s.err = err

		return '\n'
	} else if ch == '\n' {
		// We're in both style
		s.le = s.leBoth
		return '\n'
	}

	// We're in carriage return style
	s.le = s.leCarriage

	// Push back the character
	s.pushed = ch

	return '\n'
}

// leCarriage handles the carriage return line ending style.
func (s *scanner) leCarriage(ch rune) rune {
	if ch == '\r' {
		return '\n'
	}

	// Got '\n' in carriage return style?
	return ' '
}

// leNewline handles the newline line ending style.
func (s *scanner) leNewline(ch rune) rune {
	return ch
}

// leBoth handles the carriage return-newline line ending style.
func (s *scanner) leBoth(ch rune) rune {
	if ch == '\r' {
		// Get the next character, and ignore '\r' if next is
		// '\n'
		ch, err := s.nextChar()
		if ch == common.EOF || err != nil {
			// Put back the error
			s.err = err

			// EOF or error; we'll return the '\r'
			return '\r'
		}

		// If next is a newline, that's what we'll return
		if ch == '\n' {
			return '\n'
		}

		// Push the character back and return the '\r'
		s.pushed = ch
		return '\r'
	}

	return ch
}
