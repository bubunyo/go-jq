// Copyright (c) 2016 Matt Ho <matt.ho@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	testCases := map[string]struct {
		In   string
		From string
		To   string
	}{
		"simple": {
			In:   `[0]`,
			From: "0",
		},
		"range": {
			In:   `[0:1]`,
			From: "0",
			To:   "1",
		},
		"space before": {
			In:   ` [0:1]`,
			From: "0",
			To:   "1",
		},
		"space after": {
			In:   `[0:1] `,
			From: "0",
			To:   "1",
		},
		"space from": {
			In:   `[ 0 :1] `,
			From: "0",
			To:   "1",
		},
		"space to": {
			In:   `[0: 1 ] `,
			From: "0",
			To:   "1",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			matches := reArray.FindAllStringSubmatch(tc.In, -1)
			require.Len(t, matches, 1, "Should have exactly one match")
			require.Len(t, matches[0], 4, "Match should have 4 capture groups")
			assert.Equal(t, tc.From, matches[0][1], "From value should match")
			assert.Equal(t, tc.To, matches[0][3], "To value should match")
		})
	}
}
