package jq_test

import (
	"testing"

	"github.com/bubunyo/go-jq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkChain(t *testing.B) {
	op := jq.Chain(jq.Dot("a"), jq.Dot("b"))
	data := []byte(`{"a":{"b":"value"}}`)

	for i := 0; i < t.N; i++ {
		_, err := op.Apply(data)
		require.NoError(t, err)
	}
}

func TestChain(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Op       jq.Op
		Expected string
		HasError bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Op:       jq.Chain(jq.Dot("hello")),
			Expected: `"world"`,
		},
		"nested": {
			In:       `{"a":{"b":"world"}}`,
			Op:       jq.Chain(jq.Dot("a"), jq.Dot("b")),
			Expected: `"world"`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := tc.Op.Apply([]byte(tc.In))
			if tc.HasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, string(data))
			}
		})
	}
}
