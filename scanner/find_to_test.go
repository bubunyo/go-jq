package scanner_test

import (
	"testing"

	"github.com/bubunyo/jq/scanner"
)

func TestFindTo(t *testing.T) {
	testCases := map[string]struct {
		In       string
		To       int
		Expected string
		HasErr   bool
	}{
		"first": {
			In:       `["a","b","c","d","e"]`,
			To:       0,
			Expected: `["a"]`,
		},
		"second": {
			In:       `["a","b","c","d","e"]`,
			To:       1,
			Expected: `["a","b"]`,
		},
		"mixed": {
			In:       `["a",{"hello":"world"},"c","d","e"]`,
			To:       1,
			Expected: `["a",{"hello":"world"}]`,
		},
		"negative": {
			In:     `["a",{"hello":"world"},"c","d","e"]`,
			To:     -1,
			HasErr: true,
		},
		"out of bounds": {
			In:     `["a",{"hello":"world"},"c","d","e"]`,
			To:     20,
			HasErr: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := scanner.FindTo([]byte(tc.In), 0, tc.To)
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
