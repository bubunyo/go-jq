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
