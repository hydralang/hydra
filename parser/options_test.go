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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/parser/common"
)

func TestFilename(t *testing.T) {
	a := assert.New(t)
	opts := &common.Options{}

	opt := Filename("file")
	opt(opts)

	a.Equal("file", opts.Filename)
}

func TestEncoding(t *testing.T) {
	a := assert.New(t)
	opts := &common.Options{}

	opt := Encoding("enc")
	opt(opts)

	a.Equal("enc", opts.Encoding)
}

func TestTabStop(t *testing.T) {
	a := assert.New(t)
	opts := &common.Options{}

	opt := TabStop(4)
	opt(opts)

	a.Equal(4, opts.TabStop)
}
