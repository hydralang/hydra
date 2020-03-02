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
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/hydra/parser/common"
	"github.com/hydralang/hydra/testutils"
)

func TestScannerLeUnknownNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    1,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{'o', utf8.RuneSelf})
	s.le = s.leUnknown

	result := s.leUnknown('\n')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
	testutils.AssertPtrEqual(a, s.leNewline, s.le)
}

func TestScannerLeUnknownCarriageEOF(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    0,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})
	s.le = s.leUnknown

	result := s.leUnknown('\r')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
	testutils.AssertPtrEqual(a, s.leCarriage, s.le)
}

func TestScannerLeUnknownCarriageError(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    0,
		pushed: common.Err,
		err:    assert.AnError,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})
	s.le = s.leUnknown

	result := s.leUnknown('\r')

	a.Equal('\n', result)
	a.Equal(assert.AnError, s.err)
	a.Equal(common.Err, s.pushed)
	testutils.AssertPtrEqual(a, s.leCarriage, s.le)
}

func TestScannerLeUnknownBoth(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    1,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{'\n', utf8.RuneSelf})
	s.le = s.leUnknown

	result := s.leUnknown('\r')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
	testutils.AssertPtrEqual(a, s.leBoth, s.le)
}

func TestScannerLeUnknownCarriage(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    1,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{'o', utf8.RuneSelf})
	s.le = s.leUnknown

	result := s.leUnknown('\r')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal('o', s.pushed)
	testutils.AssertPtrEqual(a, s.leCarriage, s.le)
}

func TestScannerLeCarriageCarriage(t *testing.T) {
	a := assert.New(t)
	s := &scanner{}

	result := s.leCarriage('\r')

	a.Equal('\n', result)
}

func TestScannerLeCarriageNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{}

	result := s.leCarriage('\n')

	a.Equal(' ', result)
}

func TestScannerLeNewlineCarriage(t *testing.T) {
	a := assert.New(t)
	s := &scanner{}

	result := s.leNewline('\r')

	a.Equal('\r', result)
}

func TestScannerLeNewlineNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{}

	result := s.leNewline('\n')

	a.Equal('\n', result)
}

func TestScannerLeBothCarriageEOF(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    0,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	result := s.leBoth('\r')

	a.Equal('\r', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
}

func TestScannerLeBothCarriageError(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    0,
		pushed: common.Err,
		err:    assert.AnError,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	result := s.leBoth('\r')

	a.Equal('\r', result)
	a.Equal(assert.AnError, s.err)
	a.Equal(common.Err, s.pushed)
}

func TestScannerLeBothCarriageNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    1,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{'\n', utf8.RuneSelf})

	result := s.leBoth('\r')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
}

func TestScannerLeBothCarriageOther(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    1,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{'o', utf8.RuneSelf})

	result := s.leBoth('\r')

	a.Equal('\r', result)
	a.Nil(s.err)
	a.Equal('o', s.pushed)
}

func TestScannerLeBothNewline(t *testing.T) {
	a := assert.New(t)
	s := &scanner{
		end:    0,
		pushed: common.Err,
	}
	copy(s.buf[0:], []byte{utf8.RuneSelf})

	result := s.leBoth('\n')

	a.Equal('\n', result)
	a.Nil(s.err)
	a.Equal(common.Err, s.pushed)
}
