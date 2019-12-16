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

import "github.com/hydralang/hydra/parser/common"

// Filename sets the filename being scanned.  If not set, an attempt
// is made to guess it from the source (depends on source having a
// Name() method returning a string), and a default is used if that
// fails.
func Filename(file string) common.Option {
	return func(opts *common.Options) {
		opts.Filename = file
	}
}

// Encoding sets the encoding for the file being scanned.  If not set,
// an attempt is made to guess it from the source (depends on source
// implementing io.Seeker), and a default of "utf-8" is used if that
// fails.
func Encoding(encoding string) common.Option {
	return func(opts *common.Options) {
		opts.Encoding = encoding
	}
}

// TabStop sets the size of a tab stop.  If not set, it defaults to 8.
func TabStop(tabstop int) common.Option {
	return func(opts *common.Options) {
		opts.TabStop = tabstop
	}
}
