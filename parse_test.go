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

package jq_test

import (
	"testing"

	"github.com/bubunyo/jq"
)

func TestParse(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Op       string
		Expected string
		HasError bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Op:       ".hello",
			Expected: `"world"`,
		},
		"nested": {
			In:       `{"a":{"b":"world"}}`,
			Op:       ".a.b",
			Expected: `"world"`,
		},
		"index": {
			In:       `["a","b","c"]`,
			Op:       ".[1]",
			Expected: `"b"`,
		},
		"range": {
			In:       `["a","b","c"]`,
			Op:       ".[1:2]",
			Expected: `["b","c"]`,
		},
		"from": {
			In:       `["a","b","c","d"]`,
			Op:       ".[1:]",
			Expected: `["b","c","d"]`,
		},
		"to": {
			In:       `["a","b","c","d"]`,
			Op:       ".[:2]",
			Expected: `["a","b","c"]`,
		},
		"all": {
			In:       `["a","b","c","d"]`,
			Op:       ".[]",
			Expected: `["a","b","c","d"]`,
		},
		"nested index": {
			In:       `{"abc":"-","def":["a","b","c"]}`,
			Op:       ".def.[1]",
			Expected: `"b"`,
		},
		"nested range": {
			In:       `{"abc":"-","def":["a","b","c"]}`,
			Op:       ".def.[1:2]",
			Expected: `["b","c"]`,
		},

		// Pipe operator tests
		"pipe with dots": {
			In:       `{"a":{"b":"value"}}`,
			Op:       ".a|.b",
			Expected: `"value"`,
		},
		"multiple pipes": {
			In:       `{"a":{"b":{"c":"nested"}}}`,
			Op:       ".a|.b|.c",
			Expected: `"nested"`,
		},
		"pipe with array index": {
			In:       `[{"name":"alice"},{"name":"bob"}]`,
			Op:       ".[1]|.name",
			Expected: `"bob"`,
		},
		"pipe array operations": {
			In:       `[["a","b","c"],["d","e","f"]]`,
			Op:       ".[0]|.[1]",
			Expected: `"b"`,
		},
		"pipe with range": {
			In:       `{"items":["a","b","c","d"]}`,
			Op:       ".items|.[1:2]",
			Expected: `["b","c"]`,
		},
		"pipe mixed operations": {
			In:       `{"users":[{"id":1},{"id":2}]}`,
			Op:       ".users|.[0]|.id",
			Expected: `1`,
		},
		"pipe with from": {
			In:       `[{"a":1},{"a":2},{"a":3}]`,
			Op:       ".[1:]|.[0]|.a",
			Expected: `2`,
		},
		"pipe with b64_decode": {
			In:       `{"hello":{"world":{"object":"ewogICAgICAgICJmaXJzdCI6ICJqb2UiCiAgICAgIH0="}}}`,
			Op:       ".hello.world.object|b64_decode|.first",
			Expected: `"joe"`,
		},

		// Error cases
		"empty pipe segment": {
			In:       `{"a":"value"}`,
			Op:       ".a||.b",
			HasError: true,
		},
		"pipe at start": {
			In:       `{"a":"value"}`,
			Op:       "|.a",
			HasError: true,
		},
		"pipe at end": {
			In:       `{"a":"value"}`,
			Op:       ".a|",
			HasError: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			op, err := jq.Parse(tc.Op)
			if tc.HasError {
				// For error cases, we expect either Parse to fail or Apply to fail
				if err != nil {
					// Parse failed as expected
					return
				}
				// Parse succeeded, check if Apply fails
				_, err = op.Apply([]byte(tc.In))
				if err == nil {
					t.FailNow()
				}
				return
			}

			// For non-error cases, Parse should succeed
			if err != nil {
				t.FailNow()
			}

			data, err := op.Apply([]byte(tc.In))
			if string(data) != tc.Expected {
				t.FailNow()
			}
			if err != nil {
				t.FailNow()
			}
		})
	}
}

func TestFindIndices(t *testing.T) {
	testCases := map[string]struct {
		In     string
		Expect []string
	}{
		"simple": {
			In:     "[0]",
			Expect: []string{"[0]", "0", "", ""},
		},
		"range": {
			In:     "[0:1]",
			Expect: []string{"[0:1]", "0", ":", "1"},
		},
		"from": {
			In:     "[1:]",
			Expect: []string{"[1:]", "1", ":", ""},
		},
		"to": {
			In:     "[:1]",
			Expect: []string{"[:1]", "", ":", "1"},
		},
	}
	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			matches := jq.FindIndices(tc.In)
			t.Logf("%#v", matches[0])
			if len(matches) == 0 {
				t.Log("no matches")
				t.FailNow()
			}
			if len(matches[0]) != len(tc.Expect) {
				t.Log("count mismatch")
				t.FailNow()
			}
			for k, v := range tc.Expect {
				if v != matches[0][k] {
					t.Log("expected mismatch")
					t.FailNow()
				}
			}
		})
	}
}
