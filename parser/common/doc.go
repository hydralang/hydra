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

// Package common contains common definitions and routines for the
// Hydra parser.  This includes definitions of character classes,
// errors, locations, standard tokens, and the Profile, which enables
// dynamic changes to the way the parser functions.  The Profile, in
// particular, allows for relatively easy versioning of the Hydra
// language.
//
// Character classes are in classes.go; common errors, in errors.go.
// The Location class exists in locations.go; and options, which
// houses the Profile, is in options.go.  The Profile itself is
// defined in profile.go, and basic interfaces, such as the one
// defining a scanner, are in interfaces.go.
//
// The basic tokens are defined in tokens.go, with identifiers.go,
// operators.go, and strings.go containing the code for describing
// those token types.  (The identifiers.go file contains code
// associated with keywords, which are recognized by the identifiers
// recognizer in the lexer.)
package common
