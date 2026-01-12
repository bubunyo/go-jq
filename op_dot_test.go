package jq_test

import (
	"testing"

	"github.com/bubunyo/go-jq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkDot(t *testing.B) {
	op := jq.Dot("hello")
	data := []byte(`{"hello":"world"}`)

	for i := 0; i < t.N; i++ {
		_, err := op.Apply(data)
		require.NoError(t, err)
	}
}

func TestDot(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Key      string
		Expected string
		HasError bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Key:      "hello",
			Expected: `"world"`,
		},
		"key not found": {
			In:       `{"hello":"world"}`,
			Key:      "junk",
			HasError: true,
		},
		"unclosed value": {
			In:       `{"hello":"world`,
			Key:      "hello",
			HasError: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			op := jq.Dot(tc.Key)
			data, err := op.Apply([]byte(tc.In))
			if tc.HasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, string(data))
			}
		})
	}
}
