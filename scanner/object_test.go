package scanner_test

import (
	"testing"

	"github.com/bubunyo/jq/scanner"
)

func BenchmarkObject(t *testing.B) {
	data := []byte(`{"hello":"world"}`)

	for i := 0; i < t.N; i++ {
		end, err := scanner.Object(data, 0)
		if err != nil {
			t.FailNow()
			return
		}

		if end == 0 {
			t.FailNow()
			return
		}
	}
}

func TestObject(t *testing.T) {
	testCases := map[string]struct {
		In     string
		Out    string
		HasErr bool
	}{
		"simple": {
			In:  `{"hello":"world"}`,
			Out: `{"hello":"world"}`,
		},
		"empty": {
			In:  `{}`,
			Out: `{}`,
		},
		"spaced": {
			In:  ` { "hello" : "world" } `,
			Out: ` { "hello" : "world" }`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			end, err := scanner.Object([]byte(tc.In), 0)
			if tc.HasErr {
				if err == nil {
					t.FailNow()
				}
			} else {
				data := tc.In[0:end]
				if string(data) != tc.Out {
					t.FailNow()
				}
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}
