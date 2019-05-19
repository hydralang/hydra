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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessEncodingBOM(t *testing.T) {
	a := assert.New(t)

	result := guessEncoding([]byte("\ufeffthis is some text"))

	a.Equal(defaultEncoding, result)
}

func TestGuessEncodingEmacsLine1(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal("some-system", result)
}

func TestGuessEncodingEmacsLine2CRLFIndent(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# this is a test\r\n  # -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal("some-system", result)
}

func TestGuessEncodingEmacsLine2CRIndent(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# this is a test\r  # -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal("some-system", result)
}

func TestGuessEncodingEmacsLine2LFIndent(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# this is a test\n  # -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal("some-system", result)
}

func TestGuessEncodingVimLine2CRLFIndent(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# this is a test\r\n  # vim:fileencoding=some-system",
	)

	result := guessEncoding(src)

	a.Equal("some-system", result)
}

func TestGuessEncodingEmacsLine1NotComment(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"this is a test\r\n  # -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal(defaultEncoding, result)
}

func TestGuessEncodingEmacsLine2NotComment(t *testing.T) {
	a := assert.New(t)
	src := []byte(
		"# this is a test\r\n   -*- coding: some-system -*-",
	)

	result := guessEncoding(src)

	a.Equal(defaultEncoding, result)
}

func TestOptionsParseDefaults(t *testing.T) {
	a := assert.New(t)
	obj := &Options{}

	obj.Parse()

	a.Equal(defaultFilename, obj.Filename)
	a.Equal(defaultEncoding, obj.Encoding)
}

type tnamer struct{}

func (n tnamer) Name() string {
	return "name"
}

func (n tnamer) Read(buf []byte) (int, error) {
	return 0, nil
}

func TestOptionsParseNamer(t *testing.T) {
	a := assert.New(t)
	obj := &Options{Source: tnamer{}}

	obj.Parse()

	a.Equal("name", obj.Filename)
	a.Equal(defaultEncoding, obj.Encoding)
}

func TestOptionsParseSeeker(t *testing.T) {
	a := assert.New(t)
	obj := &Options{Source: strings.NewReader("# coding: other")}

	obj.Parse()

	a.Equal(defaultFilename, obj.Filename)
	a.Equal("other", obj.Encoding)
}

func TestOptionsParseSeekerNoContent(t *testing.T) {
	a := assert.New(t)
	obj := &Options{Source: strings.NewReader("")}

	obj.Parse()

	a.Equal(defaultFilename, obj.Filename)
	a.Equal(defaultEncoding, obj.Encoding)
}

type tseeker struct{}

func (n tseeker) Seek(offset int64, whence int) (int64, error) {
	return 0, assert.AnError
}

func (n tseeker) Read(buf []byte) (int, error) {
	return 0, nil
}

func TestOptionsParseSeekerUnseekable(t *testing.T) {
	a := assert.New(t)
	obj := &Options{Source: tseeker{}}

	obj.Parse()

	a.Equal(defaultFilename, obj.Filename)
	a.Equal(defaultEncoding, obj.Encoding)
}

func TestOptionsParseOptions(t *testing.T) {
	a := assert.New(t)
	obj := &Options{}

	obj.Parse(Filename("file"), Encoding("other"))

	a.Equal("file", obj.Filename)
	a.Equal("other", obj.Encoding)
}

func TestFilename(t *testing.T) {
	a := assert.New(t)
	opts := &Options{}

	opt := Filename("file")
	opt(opts)

	a.Equal("file", opts.Filename)
}

func TestEncoding(t *testing.T) {
	a := assert.New(t)
	opts := &Options{}

	opt := Encoding("enc")
	opt(opts)

	a.Equal("enc", opts.Encoding)
}
