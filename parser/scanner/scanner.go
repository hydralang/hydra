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

import (
	"container/list"
	"io"
	"unicode/utf8"

	"github.com/hydralang/hydra/parser"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

// scanBuf is the size of the read buffer to utilize.
const scanBuf = 4096

// Scanner is an interface describing a scanner.  A scanner reads a
// source character rune by character rune, returning augmented
// characters.
type Scanner interface {
	// Next retrieves the next rune from the file.  An EOF
	// augmented character is returned on end of file, and an Err
	// augmented character is returned in the event of an error.
	Next() parser.AugChar

	// Push pushes back a single augmented character onto the
	// scanner.  Any number of characters may be pushed back.
	Push(ch parser.AugChar)
}

// scanner is an implementation of Scanner.
type scanner struct {
	source io.Reader         // The reader, including encoding
	opts   *parser.Options   // The parser options
	buf    [scanBuf + 1]byte // The read buffer
	pos    int               // The current index into the read buffer
	end    int               // The end of the buffer
	le     lineEnding        // The processor for line ending style
	pushed rune              // One char pushback for line endings
	err    error             // Deferred error
	loc    parser.Location   // Location of head of read buffer
	queue  list.List         // List of pushed-back chars
}

// Scan prepares a new scanner from the parser options.
func Scan(opts *parser.Options) (Scanner, error) {
	// Set up the encoding transform to apply to the input
	enc, err := ianaindex.IANA.Encoding(opts.Encoding)
	if err != nil {
		return nil, err
	}

	// Construct our scanner object
	s := &scanner{
		source: transform.NewReader(opts.Source, enc.NewDecoder()),
		opts:   opts,
		pushed: parser.Err, // sentinel for nothing there
		loc: parser.Location{
			File: opts.Filename,
			B:    parser.FilePos{L: 1, C: 1},
			E:    parser.FilePos{L: 1, C: 1},
		},
	}

	// Set up the buffer and line ending processor
	s.buf[0] = utf8.RuneSelf
	s.le = s.leUnknown

	return s, nil
}

// nextChar retrieves the next rune from the file.  Returns EOF at end
// of file, and Err (and a non-nil error) if an error occurred.  This
// is the inner portion of Next and does not handle pushed-back
// characters.
func (s *scanner) nextChar() (rune, error) {
	// Convert the next byte of the buffer into a rune; optimized
	// for the common case of bytes < 0x80.  (Note that much of
	// the following algorithm is adapted from
	// text/scanner/scanner.go.)
	ch, width := rune(s.buf[s.pos]), 1

	// Character is either part of a multi-byte character or the
	// sentinel for end of buffer
	if ch >= utf8.RuneSelf {
		// Do we have enough in the buffer?
		for s.pos+utf8.UTFMax > s.end && !utf8.FullRune(s.buf[s.pos:s.end]) {
			// If we reached EOF, return it or deferred
			// read error
			if s.source == nil {
				if s.err == nil {
					return parser.EOF, nil
				}

				// Save the error for return and clear
				// it
				err := s.err
				s.err = nil

				return parser.Err, err
			}

			// Don't have enough, start by shifting the
			// unread portion of the buffer to the
			// beginning
			copy(s.buf[0:], s.buf[s.pos:s.end])

			// Now read more bytes into the buffer
			bufLen := s.end - s.pos
			readLen, err := s.source.Read(s.buf[bufLen:scanBuf])
			s.pos = 0
			s.end = bufLen + readLen
			s.buf[s.end] = utf8.RuneSelf // sentinel

			// Handle any errors returned by Read()
			if err != nil {
				// Done with the source; signals end
				// of input
				s.source = nil

				// Save the error
				if err != io.EOF {
					s.err = err
				}
			}

		}

		// We know buffer has at least one byte; try again
		ch = rune(s.buf[s.pos])
		if ch >= utf8.RuneSelf {
			// Not ASCII subset of UTF-8
			ch, width = utf8.DecodeRune(s.buf[s.pos:s.end])

			// Handle erroroneous encodings
			if ch == utf8.RuneError && width == 1 {
				// Advance the location
				s.pos += width

				// Done with the source; signals end
				// of input
				s.source = nil

				return parser.Err, parser.ErrBadRune
			}
		}
	}

	// Advance the buffer position
	s.pos += width

	return ch, nil
}

// Push pushes back a single augmented character onto the scanner.
// Any number of characters may be pushed back.
func (s *scanner) Push(ch parser.AugChar) {
	// Push the character onto the queue
	s.queue.PushFront(ch)
}

// Next retrieves the next rune from the file.  An EOF augmented
// character is returned on end of file, and an Err augmented
// character is returned in the event of an error.
func (s *scanner) Next() parser.AugChar {
	// Handle characters pushed back by Push
	if s.queue.Len() > 0 {
		// Pop the first element off
		elem := s.queue.Front()
		s.queue.Remove(elem)

		// Return the character
		return elem.Value.(parser.AugChar)
	}

	// OK, get the next character to process
	var ch rune
	var err error
	if s.pushed != parser.Err {
		ch = s.pushed
		s.pushed = parser.Err

		// No need to handle line endings, because that's the
		// only thing that can push a character back
	} else {
		ch, err = s.nextChar()

		// Handle line endings
		if ch == '\r' || ch == '\n' {
			ch = s.le(ch)
		}
	}

	// Advance the location as needed
	s.opts.Advance(ch, &s.loc)

	// Classify the character and return it
	return s.opts.Classify(ch, s.loc, err)
}
