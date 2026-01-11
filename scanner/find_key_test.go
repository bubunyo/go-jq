package scanner_test

import (
	"testing"

	"github.com/bubunyo/jq/scanner"
)

func BenchmarkFindKey(t *testing.B) {
	data := []byte(`{"hello":"world"}`)

	for i := 0; i < t.N; i++ {
		out, err := scanner.FindKey(data, 0, []byte("hello"))
		if err != nil {
			t.FailNow()
			return
		}

		if string(out) != `"world"` {
			t.FailNow()
			return
		}
	}
}

func TestFindKey(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Key      string
		Expected string
		HasErr   bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Key:      "hello",
			Expected: `"world"`,
		},
		"spaced": {
			In:       ` { "hello" : "world" } `,
			Key:      "hello",
			Expected: `"world"`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := scanner.FindKey([]byte(tc.In), 0, []byte(tc.Key))
			if tc.HasErr {
				if err == nil {
					t.FailNow()
				}
			} else {
				if string(data) != tc.Expected {
					t.FailNow()
				}
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}
