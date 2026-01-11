package jq_test

import (
	"testing"

	"github.com/bubunyo/jq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				assert.Error(t, err, "Expected an error but got none")
				return
			}

			// For non-error cases, Parse should succeed
			require.NoError(t, err, "Parse should not return an error")

			data, err := op.Apply([]byte(tc.In))
			require.NoError(t, err, "Apply should not return an error")
			assert.Equal(t, tc.Expected, string(data))
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
			require.NotEmpty(t, matches, "Expected at least one match")
			assert.Equal(t, tc.Expect, matches[0])
		})
	}
}
