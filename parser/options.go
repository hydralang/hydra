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

package parser

import (
	"io"
	"regexp"
	"unicode/utf8"
)

// guessBlock is a block size for guessing a file encoding based on
// comments.
const guessBlock = 1024

// defaultFilename is a default filename to use if it cannot be
// determined from the source.
const defaultFilename = "<input>"

// defaultEncoding is a default encoding to assume for the source
// file.
const defaultEncoding = "utf-8"

// encodingRE is a regular expression that matches a coding
// declaration comment.
var encodingRE = regexp.MustCompile(
	`^(?:\s*#[^\r\n]*(?:\r?\n|\r))?\s*#[^\r\n]*coding[=:]\s*([-\w.]+)`,
)

// bom is the Unicode byte order mark (BOM), used to identify if a
// file is in Unicode.
const bom = '\ufeff'

// guessEncoding takes a string and attempts to determine the encoding
// of that string.  It first checks to see if the string has a BOM,
// implying UTF-8.  If not present, it then uses encodingRE to look
// for a coding system declaration in the first two lines of the
// string.
func guessEncoding(source []byte) string {
	// First, if it starts with a BOM, it's the default encoding
	ch, _ := utf8.DecodeRune(source)
	if ch == bom {
		return defaultEncoding
	}

	// OK, check if it matches the encoding RE
	match := encodingRE.FindSubmatch(source)
	if match != nil {
		return string(match[1])
	}

	// OK, choose the default
	return defaultEncoding
}

// Options contains the options for the parser.
type Options struct {
	Source   io.Reader // The source from which to read
	Filename string    // The name of the file being parsed
	Encoding string    // The encoding of the source
}

// namer is an interface with a single Name() method.  This matches
// the signature of the Name() method for os.File, and allows us to
// query the source for its name.
type namer interface {
	// Name retrieves the name of the object.
	Name() string
}

// Parse parses a series of options into the Options structure.
func (o *Options) Parse(opts ...Option) {
	// Just apply each option in turn
	for _, opt := range opts {
		opt(o)
	}

	// Set up default file name
	if o.Filename == "" {
		switch obj := o.Source.(type) {
		case namer: // Get the filename from the source
			o.Filename = obj.Name()

		default: // Use a default name
			o.Filename = defaultFilename
		}
	}

	// Set up default encoding
	if o.Encoding == "" {
		switch obj := o.Source.(type) {
		case io.Seeker: // Try to guess the encoding
			curLoc, err := obj.Seek(0, io.SeekCurrent)
			if err != nil {
				// Not seekable, use default encoding
				o.Encoding = defaultEncoding
				break
			}

			// Read in a block to guess encoding
			buf := make([]byte, guessBlock)
			n, err := o.Source.Read(buf)

			// Reset to the "beginning" of the file;
			// ignoring error...
			obj.Seek(curLoc, io.SeekStart)

			// Was our read successful?
			if n == 0 || err != nil {
				// Use default encoding
				o.Encoding = defaultEncoding
				break
			}

			// Guess the encoding
			o.Encoding = guessEncoding(buf)

		default: // Use a default encoding
			o.Encoding = defaultEncoding
		}
	}
}

// Option type for option functions.  Each function mutates a
// passed-in Options structure to set the specific option.
type Option func(opts *Options)

// Filename sets the filename being scanned.  If not set, an attempt
// is made to guess it from the source (depends on source having a
// Name() method returning a string), and a default is used if that
// fails.
func Filename(file string) Option {
	return func(opts *Options) {
		opts.Filename = file
	}
}

// Encoding sets the encoding for the file being scanned.  If not set,
// an attempt is made to guess it from the source (depends on source
// implementing io.Seeker), and a default of "utf-8" is used if that
// fails.
func Encoding(encoding string) Option {
	return func(opts *Options) {
		opts.Encoding = encoding
	}
}
